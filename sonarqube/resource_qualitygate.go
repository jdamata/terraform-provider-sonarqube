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
func qualityGate() *schema.Resource {
	return &schema.Resource{
		Create: qualityGateCreate,
		Read:   qualityGateRead,
		Delete: qualityGateDelete,
		Importer: &schema.ResourceImporter{
			State: qualityGateImport,
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

func qualityGateCreate(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("%s/api/qualitygates/create?name=%s",
		m.(*ProviderConfiguration).sonarURL,
		d.Get("name").(string),
	)
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
	if resp.StatusCode != http.StatusOK {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}

	qualityGateResponse, err := getQualityGateResponse(resp)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateCreate")
		return err
	}

	d.SetId(strconv.FormatInt(qualityGateResponse.ID, 10))
	return nil
}

func qualityGateRead(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("%s/api/qualitygates/show?id=%s",
		m.(*ProviderConfiguration).sonarURL,
		d.Id(),
	)
	req, err := http.NewRequest("GET", url, http.NoBody)
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

	d.SetId(strconv.FormatInt(qualityGateResponse.ID, 10))
	d.Set("name", qualityGateResponse.Name)
	return nil
}

func qualityGateDelete(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		log.WithError(err).Error("resourceQualityGateDelete")
		return err
	}

	url := fmt.Sprintf("%s/api/qualitygates/destroy?id=%v",
		m.(*ProviderConfiguration).sonarURL,
		id,
	)
	req, err := http.NewRequest("POST", url, http.NoBody)
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

func qualityGateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := qualityGateRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
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
