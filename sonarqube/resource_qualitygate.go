package sonarqube

import (
	"encoding/json"
	"errors"
	"fmt"
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
				ForceNew: true,
			},
		},
	}
}

func qualityGateCreate(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("%s/api/qualitygates/create?name=%s&organization=%s",
		m.(*ProviderConfiguration).sonarURL,
		d.Get("name").(string),
		d.Get("organization").(string),
	)
	log.Info(url)
	req, err := http.NewRequest("POST", url, http.NoBody)
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
	req, err := http.NewRequest("GET", m.(*ProviderConfiguration).sonarURL+"/api/qualitygates/show", http.NoBody)
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
	req, err := http.NewRequest("POST", m.(*ProviderConfiguration).sonarURL+"/api/qualitygates/destroy", http.NoBody)
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
