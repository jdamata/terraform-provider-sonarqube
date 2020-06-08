package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
)

func httpRequestHelper(client *retryablehttp.Client, method string, sonarqubeURL string, expectedResponseCode int, errormsg string) (http.Response, error) {
	// Prepare request
	req, err := retryablehttp.NewRequest(method, sonarqubeURL, http.NoBody)
	if err != nil {
		log.WithError(err).Error(errormsg)
		// Returning a blank http.Response object must be wrong. What am i suppose to do here??
		return http.Response{}, err
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error(errormsg)
		return *resp, err
	}

	// Check response code
	if resp.StatusCode != expectedResponseCode {
		if resp.Body == http.NoBody {
			// No error message in the body
			return *resp, fmt.Errorf("StatusCode: %v does not match expectedResponseCode: %v", resp.StatusCode, expectedResponseCode)
		}

		// The response body has content, try to decode the error
		// message
		errorResponse := ErrorResponse{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return *resp, fmt.Errorf("Failed to decode error response json into struct: %+v", err)
		}
		return *resp, fmt.Errorf("API returned an error: %+v", errorResponse.Errors[0].Message)
	}

	return *resp, nil
}
