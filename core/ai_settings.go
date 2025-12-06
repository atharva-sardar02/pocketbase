package core

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// AISettings defines the AI Query feature configuration.
type AISettings struct {
	// Enabled controls whether the AI Query feature is active.
	Enabled bool `form:"enabled" json:"enabled"`

	// Provider specifies the LLM provider to use.
	// Valid values: "openai", "ollama", "anthropic", "custom"
	Provider string `form:"provider" json:"provider"`

	// BaseURL is the API base URL for the LLM provider.
	// For OpenAI: https://api.openai.com/v1
	// For Ollama: http://localhost:11434/v1
	BaseURL string `form:"baseUrl" json:"baseUrl"`

	// APIKey is the API key for authenticating with the LLM provider.
	// This field is encrypted at rest when stored in the database.
	APIKey string `form:"apiKey" json:"apiKey,omitempty"`

	// Model specifies the LLM model to use (e.g., "gpt-4o-mini", "llama2").
	Model string `form:"model" json:"model"`

	// Temperature controls the randomness of the LLM output.
	// Range: 0.0 (deterministic) to 1.0 (creative).
	// Default: 0.1 (low randomness for consistent filter generation).
	Temperature float64 `form:"temperature" json:"temperature"`
}

// Validate makes AISettings validatable by implementing [validation.Validatable] interface.
func (c AISettings) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.Provider,
			validation.When(c.Enabled, validation.Required),
			validation.In("openai", "ollama", "anthropic", "custom"),
		),
		validation.Field(
			&c.BaseURL,
			validation.When(c.Enabled, validation.Required),
			is.URL,
		),
		validation.Field(
			&c.APIKey,
			validation.When(c.Enabled && c.Provider != "ollama", validation.Required),
		),
		validation.Field(
			&c.Model,
			validation.When(c.Enabled, validation.Required),
		),
		validation.Field(
			&c.Temperature,
			validation.Min(0.0),
			validation.Max(1.0),
		),
	)
}

// ValidateProvider checks if the provider is valid.
func (c AISettings) ValidateProvider() error {
	validProviders := []string{"openai", "ollama", "anthropic", "custom"}
	for _, valid := range validProviders {
		if c.Provider == valid {
			return nil
		}
	}
	return validation.NewError("validation_invalid_provider", "Invalid provider. Must be one of: openai, ollama, anthropic, custom")
}

// ValidateTemperature checks if the temperature is within valid range.
func (c AISettings) ValidateTemperature() error {
	if c.Temperature < 0.0 || c.Temperature > 1.0 {
		return validation.NewError("validation_invalid_temperature", "Temperature must be between 0.0 and 1.0")
	}
	return nil
}


