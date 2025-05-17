package sonarqube

import (
	"testing"
)

func TestSanitizeSensitiveURLs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic auth credentials",
			input:    "https://username:password@example.com",
			expected: "https://***:***@example.com",
		},
		{
			name:     "token parameter",
			input:    "https://example.com/api?token=secret123&other=value",
			expected: "https://example.com/api?token=***&other=value",
		},
		{
			name:     "password parameter",
			input:    "https://example.com/api?password=secret123&other=value",
			expected: "https://example.com/api?password=***&other=value",
		},
		{
			name:     "secret parameter",
			input:    "https://example.com/api?secret=secret123&other=value",
			expected: "https://example.com/api?secret=***&other=value",
		},
		{
			name:     "multiple parameters",
			input:    "https://example.com/api?token=abc123&password=pass456&secret=sec789",
			expected: "https://example.com/api?token=***&password=***&secret=***",
		},
		{
			name:     "basic auth and parameter",
			input:    "https://user:pass@example.com/api?token=secret123",
			expected: "https://***:***@example.com/api?token=***",
		},
		{
			name:     "error message with URL",
			input:    "failed to connect to https://username:password@example.com/api?token=secret123",
			expected: "failed to connect to https://***:***@example.com/api?token=***",
		},
		{
			name:     "error message with URL",
			input:    "failed to connect to https://username:password@example.com/api?token=secret123 because of xyz",
			expected: "failed to connect to https://***:***@example.com/api?token=*** because of xyz",
		},
		{
			name:     "no sensitive data",
			input:    "https://example.com/api?param=value",
			expected: "https://example.com/api?param=value",
		},
		{
			name:     "quoted URL in JSON or string",
			input:    `error occurred with "https://user:pass@example.com?secret=abc"`,
			expected: `error occurred with "https://***:***@example.com?secret=***"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeSensitiveURLs(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeSensitiveURLs() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCensorHttpError(t *testing.T) {
	tests := []struct {
		name        string
		inputError  error
		expectedMsg string
	}{
		{
			name:        "error with basic auth",
			inputError:  errorWithMessage("failed to connect to https://user:pass@example.com"),
			expectedMsg: "failed to connect to https://***:***@example.com",
		},
		{
			name:        "error with parameter",
			inputError:  errorWithMessage("error in request to https://api.example.com?token=abc123"),
			expectedMsg: "error in request to https://api.example.com?token=***",
		},
		{
			name:        "error with no data",
			inputError:  errorWithMessage("connection timeout"),
			expectedMsg: "connection timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := censorHttpError(tt.inputError)
			if result.Error() != tt.expectedMsg {
				t.Errorf("censorHttpError() = %v, want %v", result.Error(), tt.expectedMsg)
			}
		})
	}
}

func errorWithMessage(msg string) error {
	return &testError{message: msg}
}

type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
