package sonarqube

import (
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Deprecate this
func getResponseBodyAsString(resp *http.Response) (string, error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithField("status code", resp.StatusCode).WithError(err).Error("getResponseBodyAsString")
		return "", err
	}
	return string(bodyBytes), nil
}

func httpRequestHelper(client http.Client, method string, sonarqubeURL string, expectedResponseCode int, errormsg string) http.Response {
	// Prepare request
	req, err := http.NewRequest(method, sonarqubeURL, http.NoBody)
	if err != nil {
		log.WithError(err).Error(errormsg)
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Error(errormsg)
	}

	// Check response code
	if resp.StatusCode != expectedResponseCode {
		log.Error(errormsg + "StatusCode code does not match expectedResponseCode")
	}

	return *resp
}
