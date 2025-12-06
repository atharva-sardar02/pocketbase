package ai

import "fmt"

// AIClientError is the base error type for AI client errors.
type AIClientError struct {
	Message string
	Code    int
	Err     error
}

func (e *AIClientError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AIClientError) Unwrap() error {
	return e.Err
}

// AIRateLimitError represents a 429 rate limit response from the LLM API.
type AIRateLimitError struct {
	AIClientError
	RetryAfter int // seconds
}

func NewAIRateLimitError(retryAfter int) *AIRateLimitError {
	return &AIRateLimitError{
		AIClientError: AIClientError{
			Message: "Rate limit exceeded",
			Code:    429,
		},
		RetryAfter: retryAfter,
	}
}

// AIAuthError represents a 401 authentication error from the LLM API.
type AIAuthError struct {
	AIClientError
}

func NewAIAuthError(message string) *AIAuthError {
	return &AIAuthError{
		AIClientError: AIClientError{
			Message: message,
			Code:    401,
		},
	}
}

// AITimeoutError represents a timeout error when calling the LLM API.
type AITimeoutError struct {
	AIClientError
}

func NewAITimeoutError() *AITimeoutError {
	return &AITimeoutError{
		AIClientError: AIClientError{
			Message: "Request timeout: LLM API call exceeded 30 seconds",
			Code:    0,
		},
	}
}

