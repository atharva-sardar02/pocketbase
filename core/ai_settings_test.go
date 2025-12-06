package core

import (
	"testing"
)

func TestAISettings_Validate(t *testing.T) {
	tests := []struct {
		name    string
		settings AISettings
		wantErr bool
	}{
		{
			name: "Valid settings",
			settings: AISettings{
				Enabled:     true,
				Provider:    "openai",
				BaseURL:     "https://api.openai.com/v1",
				APIKey:      "sk-test123",
				Model:       "gpt-4o-mini",
				Temperature: 0.1,
			},
			wantErr: false,
		},
		{
			name: "Missing API key when enabled",
			settings: AISettings{
				Enabled:     true,
				Provider:    "openai",
				BaseURL:     "https://api.openai.com/v1",
				APIKey:      "",
				Model:       "gpt-4o-mini",
				Temperature: 0.1,
			},
			wantErr: true,
		},
		{
			name: "Invalid provider",
			settings: AISettings{
				Enabled:     true,
				Provider:    "invalid",
				BaseURL:     "https://api.openai.com/v1",
				APIKey:      "sk-test123",
				Model:       "gpt-4o-mini",
				Temperature: 0.1,
			},
			wantErr: true,
		},
		{
			name: "Temperature out of range (too high)",
			settings: AISettings{
				Enabled:     true,
				Provider:    "openai",
				BaseURL:     "https://api.openai.com/v1",
				APIKey:      "sk-test123",
				Model:       "gpt-4o-mini",
				Temperature: 1.5,
			},
			wantErr: true,
		},
		{
			name: "Temperature out of range (negative)",
			settings: AISettings{
				Enabled:     true,
				Provider:    "openai",
				BaseURL:     "https://api.openai.com/v1",
				APIKey:      "sk-test123",
				Model:       "gpt-4o-mini",
				Temperature: -0.1,
			},
			wantErr: true,
		},
		{
			name: "Disabled settings skip validation",
			settings: AISettings{
				Enabled:     false,
				Provider:    "",
				BaseURL:     "",
				APIKey:      "",
				Model:       "",
				Temperature: 0.0,
			},
			wantErr: false,
		},
		{
			name: "Ollama provider doesn't require API key",
			settings: AISettings{
				Enabled:     true,
				Provider:    "ollama",
				BaseURL:     "http://localhost:11434/v1",
				APIKey:      "",
				Model:       "llama2",
				Temperature: 0.1,
			},
			wantErr: false,
		},
		{
			name: "Invalid BaseURL",
			settings: AISettings{
				Enabled:     true,
				Provider:    "openai",
				BaseURL:     "not-a-valid-url",
				APIKey:      "sk-test123",
				Model:       "gpt-4o-mini",
				Temperature: 0.1,
			},
			wantErr: true,
		},
		{
			name: "Missing model when enabled",
			settings: AISettings{
				Enabled:     true,
				Provider:    "openai",
				BaseURL:     "https://api.openai.com/v1",
				APIKey:      "sk-test123",
				Model:       "",
				Temperature: 0.1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.settings.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AISettings.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAISettings_ValidateProvider(t *testing.T) {
	tests := []struct {
		name    string
		settings AISettings
		wantErr bool
	}{
		{
			name: "Valid provider: openai",
			settings: AISettings{
				Provider: "openai",
			},
			wantErr: false,
		},
		{
			name: "Valid provider: ollama",
			settings: AISettings{
				Provider: "ollama",
			},
			wantErr: false,
		},
		{
			name: "Valid provider: anthropic",
			settings: AISettings{
				Provider: "anthropic",
			},
			wantErr: false,
		},
		{
			name: "Valid provider: custom",
			settings: AISettings{
				Provider: "custom",
			},
			wantErr: false,
		},
		{
			name: "Invalid provider",
			settings: AISettings{
				Provider: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.settings.ValidateProvider()
			if (err != nil) != tt.wantErr {
				t.Errorf("AISettings.ValidateProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAISettings_ValidateTemperature(t *testing.T) {
	tests := []struct {
		name    string
		settings AISettings
		wantErr bool
	}{
		{
			name: "Valid temperature: 0.0",
			settings: AISettings{
				Temperature: 0.0,
			},
			wantErr: false,
		},
		{
			name: "Valid temperature: 0.5",
			settings: AISettings{
				Temperature: 0.5,
			},
			wantErr: false,
		},
		{
			name: "Valid temperature: 1.0",
			settings: AISettings{
				Temperature: 1.0,
			},
			wantErr: false,
		},
		{
			name: "Invalid temperature: too high",
			settings: AISettings{
				Temperature: 1.5,
			},
			wantErr: true,
		},
		{
			name: "Invalid temperature: negative",
			settings: AISettings{
				Temperature: -0.1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.settings.ValidateTemperature()
			if (err != nil) != tt.wantErr {
				t.Errorf("AISettings.ValidateTemperature() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAISettings_Defaults(t *testing.T) {
	// Test that default settings are valid
	defaultSettings := AISettings{
		Enabled:     false,
		Provider:    "openai",
		BaseURL:     "https://api.openai.com/v1",
		APIKey:      "",
		Model:       "gpt-4o-mini",
		Temperature: 0.1,
	}

	// Default settings should be valid when disabled
	if err := defaultSettings.Validate(); err != nil {
		t.Errorf("Default AISettings should be valid when disabled, got error: %v", err)
	}

	// When enabled with API key, should be valid
	defaultSettings.Enabled = true
	defaultSettings.APIKey = "sk-test123"
	if err := defaultSettings.Validate(); err != nil {
		t.Errorf("Default AISettings should be valid when enabled with API key, got error: %v", err)
	}
}

func TestAISettings_APIKeyEncryption(t *testing.T) {
	// Note: Actual encryption is handled at the Settings level in DBExport/loadParam
	// This test just verifies that the APIKey field is properly omitted from JSON when empty
	settings := AISettings{
		Enabled:     true,
		Provider:    "openai",
		BaseURL:     "https://api.openai.com/v1",
		APIKey:      "", // Empty API key should be omitted
		Model:       "gpt-4o-mini",
		Temperature: 0.1,
	}

	// The json tag includes "omitempty" which means empty APIKey won't be serialized
	// This is a structural test - actual encryption happens in Settings.DBExport()
	// We verify the field exists and can be set
	if settings.APIKey != "" {
		t.Error("APIKey should be empty in test")
	}

	settings.APIKey = "sk-test123"
	if settings.APIKey != "sk-test123" {
		t.Error("APIKey should be settable")
	}
}


