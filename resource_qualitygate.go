package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

const (
	sonarURL    = "127.0.0.1:9000"
	contentType = "application/json"
)

// Returns the resource represented by this file.
func qualityGate() *schema.Resource {
	return &schema.Resource{
		// Operations required by every Terraform resource.
		Create: qualityGateCreate,
		// Read:   qualityGateRead,
		// Update: qualityGateUpdate,
		// Delete: qualityGateDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"organization": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
		},
	}
}

func qualityGateCreate(d *schema.ResourceData, m interface{}) error {
	log.Info("Creating quality gate")
	qualityGateName := d.Get("name").(string)
	qualityGate := QualityGate{Name: qualityGateName}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(qualityGate)
	resp, err := http.Post(sonarURL+"/api/qualitygates/create", contentType, buffer)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateCreate")
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}
	/*
	 * Why return nil?
	 * Please take a look at the rules for update the state in Terraform defined here:
	 * https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	 */
	return nil
}
