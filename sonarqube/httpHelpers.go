package sonarqube

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
	"regexp"
)

// ErrorResponse struct
type ErrorResponse struct {
	Errors []ErrorMessage `json:"errors,omitempty"`
}

// ErrorMessage struct
type ErrorMessage struct {
	Message string `json:"msg,omitempty"`
}

// Paging used in /search API endpoints
type Paging struct {
	PageIndex int64 `json:"pageIndex"`
	PageSize  int64 `json:"pageSize"`
	Total     int64 `json:"total"`
}

// helper function to make api request to sonarqube
func httpRequestHelper(client *retryablehttp.Client, method string, sonarqubeURL string, expectedResponseCode int, resource string) (http.Response, error) {
	// Prepare request
	req, err := retryablehttp.NewRequest(method, sonarqubeURL, http.NoBody)
	if err != nil {
		return http.Response{}, fmt.Errorf("failed to create request for resource %s: %w", resource, censorHttpError(err))
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return http.Response{}, fmt.Errorf("failed to send request for resource %s: %w", resource, censorHttpError(err))
	}

	// Check response code
	if resp.StatusCode != expectedResponseCode {
		if resp.Body == http.NoBody {
			// No error message in the body
			return *resp, fmt.Errorf("statusCode: %v does not match expectedResponseCode: %v for resource %s", resp.StatusCode, expectedResponseCode, resource)
		}

		// The response body has content, try to decode the error message
		errorResponse := ErrorResponse{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return *resp, fmt.Errorf("failed to decode error response json into struct for resource %s: %+v", resource, err)
		}
		if len(errorResponse.Errors) == 0 {
			return *resp, fmt.Errorf("statusCode: %v does not match expectedResponseCode for resource %s: %v. No error message found in the response body", resp.StatusCode, resource, expectedResponseCode)
		}
		return *resp, fmt.Errorf("API returned an error for resource %s: %+v", resource, errorResponse.Errors[0].Message)
	}

	return *resp, nil
}

func censorHttpError(error error) error {
	sanitizedError := sanitizeSensitiveURLs(error.Error())
	return errors.New(sanitizedError)
}

var (
	regexBasicAuth = regexp.MustCompile(`(https?://)([^:]+:)([^@]+)`)
	regexToken     = regexp.MustCompile(`([&?]token=)([^&"' ]*)`)
	regexPassword  = regexp.MustCompile(`([&?]password=)([^&"' ]*)`)
	regexSecret    = regexp.MustCompile(`([&?]secret=)([^&"' ]*)`)
)

func sanitizeSensitiveURLs(input string) string {
	outputString := regexBasicAuth.ReplaceAllString(input, "${1}***:***")

	outputString = regexToken.ReplaceAllString(outputString, "${1}***")
	outputString = regexPassword.ReplaceAllString(outputString, "${1}***")
	outputString = regexSecret.ReplaceAllString(outputString, "${1}***")

	return outputString
}
