package sonarqube

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Returns the resource represented by this file.
func resourceSonarqubeQualityGate() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeQualityGateCreate,
		Read:   resourceSonarqubeQualityGateRead,
		Delete: resourceSonarqubeQualityGateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeQualityGateImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSonarqubeQualityGateCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/create"
	sonarQubeURL.ForceQuery = true
	sonarQubeURL.RawQuery = fmt.Sprintf("name=%s",
		d.Get("name").(string),
	)

	req, err := http.NewRequest("POST", sonarQubeURL.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateCreate")
		return err
	}
	resp, err := m.(*ProviderConfiguration).httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateCreate")
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}

	qualityGateResponse := CreateQualityGateResponse{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateResponse)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateCreate")
	}

	d.SetId(strconv.FormatInt(qualityGateResponse.ID, 10))
	return nil
}

func resourceSonarqubeQualityGateRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/show"
	sonarQubeURL.ForceQuery = true
	sonarQubeURL.RawQuery = fmt.Sprintf("id=%s",
		d.Id(),
	)

	req, err := http.NewRequest("GET", sonarQubeURL.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateRead")
		return err
	}
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

	qualityGateReadResponse := GetQualityGate{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateReadResponse)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateRead")
	}

	d.SetId(strconv.FormatInt(qualityGateReadResponse.ID, 10))
	d.Set("name", qualityGateReadResponse.Name)
	return nil
}

func resourceSonarqubeQualityGateDelete(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		log.WithError(err).Error("resourceQualityGateDelete")
		return err
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/destroy"
	sonarQubeURL.ForceQuery = true
	sonarQubeURL.RawQuery = fmt.Sprintf("id=%v",
		id,
	)

	req, err := http.NewRequest("POST", sonarQubeURL.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateDelete")
		return err
	}
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

func resourceSonarqubeQualityGateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeQualityGateRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
