package main

import (
	"bytes"
	"encoding/json"
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

func encodeObject(obj interface{}) *bytes.Buffer {
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(obj)
	return buffer
}
