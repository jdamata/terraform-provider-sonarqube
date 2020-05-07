package sonarqube

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Returns the resource represented by this file.
func qualityGate() *schema.Resource {
	return &schema.Resource{
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

func buildQualityGateStruct(d *schema.ResourceData) QualityGate {
	qualityGateName := d.Get("name").(string)
	qualityGate := QualityGate{Name: qualityGateName}
	return qualityGate
}

func qualityGateCreate(d *schema.ResourceData, m interface{}) error {
	qualityGate := buildQualityGateStruct(d)
	buffer := encodeObject(qualityGate)

	req, err := http.NewRequest("POST", m.(*ProviderConfiguration).sonarURL+"/api/qualitygates/create", buffer)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateCreate")
		return err
	}

	req.SetBasicAuth(m.(*ProviderConfiguration).sonarUser, m.(*ProviderConfiguration).sonarPass)
	resp, err := m.(*ProviderConfiguration).httpClient.Do(req)
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
	d.SetId(qualityGateResponse.ID)
	/*
	 * Why return nil?
	 * Please take a look at the rules for update the state in Terraform defined here:
	 * https://www.terraform.io/docs/extend/writing-custom-providers.html#error-handling-amp-partial-state
	 */
	return nil
}

func qualityGateRead(d *schema.ResourceData, m interface{}) error {
	qualityGate := buildQualityGateStruct(d)
	buffer := encodeObject(qualityGate)

	req, err := http.NewRequest("GET", m.(*ProviderConfiguration).sonarURL+"/api/qualitygates/show", buffer)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateRead")
		return err
	}

	req.SetBasicAuth(m.(*ProviderConfiguration).sonarUser, m.(*ProviderConfiguration).sonarPass)
	resp, err := m.(*ProviderConfiguration).httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateRead")
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
		log.WithError(err).Error("resourceQualityGateRead")
		return err
	}

	d.Set("name", qualityGateResponse.Name)
	return nil
}

func qualityGateDelete(d *schema.ResourceData, m interface{}) error {
	qualityGate := buildQualityGateStruct(d)
	buffer := encodeObject(qualityGate)

	req, err := http.NewRequest("POST", m.(*ProviderConfiguration).sonarURL+"/api/qualitygates/destroy", buffer)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateDelete")
		return err
	}

	req.SetBasicAuth(m.(*ProviderConfiguration).sonarUser, m.(*ProviderConfiguration).sonarPass)
	resp, err := m.(*ProviderConfiguration).httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateDelete")
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
