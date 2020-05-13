package sonarqube

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func httpRequestHelper(client http.Client, method string, sonarqubeURL string, expectedResponseCode int, errormsg string) (http.Response, error) {
	// Prepare request
	req, err := http.NewRequest(method, sonarqubeURL, http.NoBody)
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
		return *resp, fmt.Errorf("StatusCode: %v does not match expectedResponseCode: %v", resp.StatusCode, expectedResponseCode)
	}

	return *resp, nil
}
