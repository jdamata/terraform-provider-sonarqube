package main

import (
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func getResponseBodyAsString(resp *http.Response) string {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithField("status code", resp.StatusCode).WithError(err).Error("getResponseBodyAsString")
		return "Could not convert response body to string"
	}
	return string(bodyBytes)
}
