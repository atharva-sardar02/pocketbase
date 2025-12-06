package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

const (
	// DefaultTimeout is the default timeout for LLM API calls (30 seconds).
	DefaultTimeout = 30 * time.Second

	// MaxRetries is the maximum number of retries for transient failures.
	MaxRetries = 2

	// RetryDelay is the delay between retries.
	RetryDelay = 1 * time.Second
)

// OpenAIClient handles communication with OpenAI-compatible LLM APIs.
type OpenAIClient struct {
	settings core.AISettings
	client   *http.Client
}

// NewOpenAIClient creates a new OpenAI client with the given settings.
func NewOpenAIClient(settings core.AISettings) *OpenAIClient {
	return &OpenAIClient{
		settings: settings,
		client: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// ChatCompletionRequest represents the request body for OpenAI chat completion API.
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

// Message represents a single message in the chat completion request.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse represents the response from OpenAI chat completion API.
type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

// Choice represents a single choice in the completion response.
type Choice struct {
	Message Message `json:"message"`
}

// APIError represents an error response from the OpenAI API.
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// SendCompletion sends a completion request to the LLM API and returns the generated text.
func (c *OpenAIClient) SendCompletion(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	// Build request
	reqBody := ChatCompletionRequest{
		Model: c.settings.Model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
		Temperature: c.settings.Temperature,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", &AIClientError{
			Message: "Failed to marshal request body",
			Err:     err,
		}
	}

	// Build URL
	url := fmt.Sprintf("%s/chat/completions", c.settings.BaseURL)

	// Retry logic
	var lastErr error
	for attempt := 0; attempt <= MaxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry
			select {
			case <-ctx.Done():
				return "", NewAITimeoutError()
			case <-time.After(RetryDelay):
			}
		}

		result, err := c.sendRequest(ctx, url, jsonBody)
		if err == nil {
			return result, nil
		}

		lastErr = err

		// Don't retry on certain errors
		if _, ok := err.(*AIAuthError); ok {
			return "", err
		}
		if _, ok := err.(*AITimeoutError); ok {
			return "", err
		}

		// Check if context is cancelled
		if ctx.Err() != nil {
			return "", NewAITimeoutError()
		}
	}

	return "", lastErr
}

// sendRequest performs a single HTTP request to the LLM API.
func (c *OpenAIClient) sendRequest(ctx context.Context, url string, body []byte) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", &AIClientError{
			Message: "Failed to create request",
			Err:     err,
		}
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.settings.APIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.settings.APIKey))
	}

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		// Check for timeout
		if ctx.Err() == context.DeadlineExceeded {
			return "", NewAITimeoutError()
		}
		return "", &AIClientError{
			Message: "Failed to send request",
			Err:     err,
		}
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", &AIClientError{
			Message: "Failed to read response body",
			Err:     err,
		}
	}

	// Handle error status codes
	switch resp.StatusCode {
	case 401:
		return "", NewAIAuthError("Invalid API key or authentication failed")
	case 429:
		retryAfter := 0
		if retryAfterStr := resp.Header.Get("Retry-After"); retryAfterStr != "" {
			fmt.Sscanf(retryAfterStr, "%d", &retryAfter)
		}
		return "", NewAIRateLimitError(retryAfter)
	case 200:
		// Parse successful response
		var completionResp ChatCompletionResponse
		if err := json.Unmarshal(respBody, &completionResp); err != nil {
			return "", &AIClientError{
				Message: "Failed to parse response",
				Err:     err,
			}
		}

		// Check for API error in response
		if completionResp.Error != nil {
			return "", &AIClientError{
				Message: completionResp.Error.Message,
				Code:    resp.StatusCode,
			}
		}

		// Extract content from first choice
		if len(completionResp.Choices) == 0 {
			return "", &AIClientError{
				Message: "No choices in response",
				Code:    resp.StatusCode,
			}
		}

		return completionResp.Choices[0].Message.Content, nil
	default:
		// Try to parse error response
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err == nil && apiErr.Message != "" {
			return "", &AIClientError{
				Message: apiErr.Message,
				Code:    resp.StatusCode,
			}
		}

		return "", &AIClientError{
			Message: fmt.Sprintf("Unexpected status code: %d", resp.StatusCode),
			Code:    resp.StatusCode,
		}
	}
}

