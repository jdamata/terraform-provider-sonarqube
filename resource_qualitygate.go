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
	sonarURL    = "http://127.0.0.1:9000"
	contentType = "application/json"
)

// Returns the resource represented by this file.
func qualityGate() *schema.Resource {
	return &schema.Resource{
		// Operations required by every Terraform resource.
		Create: qualityGateCreate,
		Read:   qualityGateRead,
		// Update: qualityGateUpdate,
		Delete: qualityGateDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func qualityGateCreate(d *schema.ResourceData, m interface{}) error {
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

	qualityGateResponse, err := getQualityGateResponse(resp)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateCreate")
		return err
	}
	d.SetId(qualityGateResponse.Id)
	/*
	 * Why return nil?
	 * Please take a look at the rules for update the state in Terraform defined here:
	 * https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	 */
	return nil
}

func qualityGateRead(d *schema.ResourceData, m interface{}) error {
	resp, err := http.Get(sonarURL + "/api/qualitygates/show/" + d.Id())

	if err != nil {
		log.WithError(err).Error("qualityGateRead")
		return err
	}

	defer resp.Body.Close()
	log.WithField("status code", resp.StatusCode).Info("Response from server")
	if resp.StatusCode != http.StatusOK {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}

	qualityGateResponse, err := getQualityGateResponse(resp)
	if err != nil {
		log.WithError(err).Error("qualityGateRead")
		return err
	}

	d.Set("name", qualityGateResponse.Name)
	return nil
}

func qualityGateDelete(d *schema.ResourceData, m interface{}) error {
	qualityGateName := d.Get("name").(string)
	qualityGate := QualityGate{Name: qualityGateName}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(qualityGate)
	resp, err := http.Post(sonarURL+"/api/qualitygates/destroy", contentType, buffer)
	if err != nil {
		log.WithError(err).Error("qualityGateDelete")
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}

	return nil
}

func getQualityGateResponse(resp *http.Response) (QualityGateResponse, error) {
	qualityGateResponse := QualityGateResponse{}
	err := json.NewDecoder(resp.Body).Decode(&qualityGateResponse)
	if err != nil {
		log.WithError(err).Error("getQualityGateResponse")
		return qualityGateResponse, err
	}
	return qualityGateResponse, nil
}
