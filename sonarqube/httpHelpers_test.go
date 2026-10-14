package sonarqube

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
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

func TestHttpRequestHelperErrorMessage(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{"errors list", `{"errors":[{"msg":"Project not found"}]}`, "API returned an error for resource test: Project not found"},
		{"result field", `{"result":"Password must be at least 12 characters long"}`, "API returned an error for resource test: Password must be at least 12 characters long"},
		{"v2 message field", `{"message":"Group not found"}`, "API returned an error for resource test: Group not found"},
		{"empty errors list", `{"errors":[]}`, "statusCode: 400 does not match expectedResponseCode: 200 for resource test. No error message found in the response body"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()
			client := retryablehttp.NewClient()
			client.RetryMax = 0
			client.Logger = nil

			resp, err := httpRequestHelper(client, "POST", server.URL, http.StatusOK, "test")
			if err == nil {
				t.Fatal("expected an error")
			}
			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", resp.StatusCode)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error = %q, want it to contain %q", err.Error(), tt.want)
			}
		})
	}
}

func TestReadPermissionPagesFollowsEveryPage(t *testing.T) {
	pages := map[string]string{
		"1": `{"paging":{"total":5},"groups":[{"name":"a"},{"name":"b"}]}`,
		"2": `{"paging":{"total":5},"groups":[{"name":"c"},{"name":"d"}]}`,
		"3": `{"paging":{"total":5},"groups":[{"name":"e"}]}`,
	}
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write([]byte(pages[r.URL.Query().Get("p")]))
	}))
	defer server.Close()
	endpoint, _ := url.Parse(server.URL)
	client := retryablehttp.NewClient()
	client.Logger = nil
	m := &ProviderConfiguration{httpClient: client}

	names := []string{}
	found, err := readPermissionPages(m, *endpoint, url.Values{"ps": {"2"}}, func(body io.Reader) (int, int64, error) {
		page := GetGroupPermissions{}
		if err := json.NewDecoder(body).Decode(&page); err != nil {
			return 0, 0, err
		}
		for _, g := range page.Groups {
			names = append(names, g.Name)
		}
		return len(page.Groups), page.Paging.Total, nil
	})
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if strings.Join(names, ",") != "a,b,c,d,e" || requests != 3 {
		t.Errorf("got %v after %d requests, want a..e after 3", names, requests)
	}
}
