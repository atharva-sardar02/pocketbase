package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAIClient_SendCompletion_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request format
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))

		// Verify request body
		var req ChatCompletionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, "gpt-4o-mini", req.Model)
		assert.Equal(t, 0.1, req.Temperature)
		assert.Len(t, req.Messages, 2)
		assert.Equal(t, "system", req.Messages[0].Role)
		assert.Equal(t, "user", req.Messages[1].Role)

		// Return mock response
		response := ChatCompletionResponse{
			Choices: []Choice{
				{
					Message: Message{
						Role:    "assistant",
						Content: `status = "active"`,
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	settings := core.AISettings{
		Enabled:     true,
		Provider:    "openai",
		BaseURL:     server.URL,
		APIKey:      "test-key",
		Model:       "gpt-4o-mini",
		Temperature: 0.1,
	}

	client := NewOpenAIClient(settings)
	result, err := client.SendCompletion(context.Background(), "system prompt", "user query")

	assert.NoError(t, err)
	assert.Equal(t, `status = "active"`, result)
}

func TestOpenAIClient_SendCompletion_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(35 * time.Second)
		w.WriteHeader(200)
	}))
	defer server.Close()

	settings := core.AISettings{
		Enabled:     true,
		Provider:    "openai",
		BaseURL:     server.URL,
		APIKey:      "test-key",
		Model:       "gpt-4o-mini",
		Temperature: 0.1,
	}

	client := NewOpenAIClient(settings)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	result, err := client.SendCompletion(ctx, "system", "user")

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.IsType(t, &AITimeoutError{}, err)
}

func TestOpenAIClient_SendCompletion_RateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(429)
		response := map[string]string{
			"error": "Rate limit exceeded",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	settings := core.AISettings{
		Enabled:     true,
		Provider:    "openai",
		BaseURL:     server.URL,
		APIKey:      "test-key",
		Model:       "gpt-4o-mini",
		Temperature: 0.1,
	}

	client := NewOpenAIClient(settings)
	result, err := client.SendCompletion(context.Background(), "system", "user")

	assert.Error(t, err)
	assert.Empty(t, result)
	rateLimitErr, ok := err.(*AIRateLimitError)
	assert.True(t, ok)
	assert.Equal(t, 60, rateLimitErr.RetryAfter)
}

func TestOpenAIClient_SendCompletion_AuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		response := map[string]string{
			"error": "Invalid API key",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	settings := core.AISettings{
		Enabled:     true,
		Provider:    "openai",
		BaseURL:     server.URL,
		APIKey:      "invalid-key",
		Model:       "gpt-4o-mini",
		Temperature: 0.1,
	}

	client := NewOpenAIClient(settings)
	result, err := client.SendCompletion(context.Background(), "system", "user")

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.IsType(t, &AIAuthError{}, err)
}

func TestOpenAIClient_SendCompletion_InvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		// Invalid JSON response
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	settings := core.AISettings{
		Enabled:     true,
		Provider:    "openai",
		BaseURL:     server.URL,
		APIKey:      "test-key",
		Model:       "gpt-4o-mini",
		Temperature: 0.1,
	}

	client := NewOpenAIClient(settings)
	result, err := client.SendCompletion(context.Background(), "system", "user")

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "Failed to parse response")
}

func TestOpenAIClient_SendCompletion_EmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ChatCompletionResponse{
			Choices: []Choice{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	settings := core.AISettings{
		Enabled:     true,
		Provider:    "openai",
		BaseURL:     server.URL,
		APIKey:      "test-key",
		Model:       "gpt-4o-mini",
		Temperature: 0.1,
	}

	client := NewOpenAIClient(settings)
	result, err := client.SendCompletion(context.Background(), "system", "user")

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "No choices in response")
}

func TestOpenAIClient_Retry(t *testing.T) {
	attempt := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++
		if attempt < 3 {
			// Return 500 error for first two attempts
			w.WriteHeader(500)
			response := map[string]string{
				"error": "Internal server error",
			}
			json.NewEncoder(w).Encode(response)
		} else {
			// Success on third attempt
			response := ChatCompletionResponse{
				Choices: []Choice{
					{
						Message: Message{
							Role:    "assistant",
							Content: "success",
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	settings := core.AISettings{
		Enabled:     true,
		Provider:    "openai",
		BaseURL:     server.URL,
		APIKey:      "test-key",
		Model:       "gpt-4o-mini",
		Temperature: 0.1,
	}

	client := NewOpenAIClient(settings)
	result, err := client.SendCompletion(context.Background(), "system", "user")

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	assert.Equal(t, 3, attempt) // Should have retried twice
}

func TestOpenAIClient_Retry_NoRetryOnAuthError(t *testing.T) {
	attempt := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++
		w.WriteHeader(401)
		response := map[string]string{
			"error": "Invalid API key",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	settings := core.AISettings{
		Enabled:     true,
		Provider:    "openai",
		BaseURL:     server.URL,
		APIKey:      "invalid-key",
		Model:       "gpt-4o-mini",
		Temperature: 0.1,
	}

	client := NewOpenAIClient(settings)
	result, err := client.SendCompletion(context.Background(), "system", "user")

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.IsType(t, &AIAuthError{}, err)
	assert.Equal(t, 1, attempt) // Should not retry on auth error
}

func TestOpenAIClient_NoAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify no Authorization header for Ollama
		assert.Empty(t, r.Header.Get("Authorization"))

		response := ChatCompletionResponse{
			Choices: []Choice{
				{
					Message: Message{
						Role:    "assistant",
						Content: "ollama response",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	settings := core.AISettings{
		Enabled:     true,
		Provider:    "ollama",
		BaseURL:     server.URL,
		APIKey:      "", // No API key for Ollama
		Model:       "llama2",
		Temperature: 0.1,
	}

	client := NewOpenAIClient(settings)
	result, err := client.SendCompletion(context.Background(), "system", "user")

	assert.NoError(t, err)
	assert.Equal(t, "ollama response", result)
}

