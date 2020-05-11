package sonarqube

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Returns the resource represented by this file.
func resourceSonarqubeQualityGateCondition() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeQualityGateConditionCreate,
		Read:   resourceSonarqubeQualityGateConditionRead,
		Update: resourceSonarqubeQualityGateConditionUpdate,
		Delete: resourceSonarqubeQualityGateConditionDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"gateid": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"error": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
			},
			"op": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSonarqubeQualityGateConditionCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/create_condition"
	query := url.Values{
		"gateId": []string{strconv.Itoa(d.Get("gateid").(int))},
		"error":  []string{strconv.Itoa(d.Get("error").(int))},
		"metric": []string{d.Get("metric").(string)},
		"op":     []string{d.Get("op").(string)},
	}
	sonarQubeURL.RawQuery = query.Encode()

	req, err := http.NewRequest("POST", sonarQubeURL.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionCreate")
		return err
	}
	resp, err := m.(*ProviderConfiguration).httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionCreate")
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}

	qualityGateConditionResponse := CreateQualityGateConditionResponse{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateConditionResponse)
	if err != nil {
		log.WithError(err).Error("getQualityGateConditionResponse")
	}

	d.SetId(strconv.FormatInt(qualityGateConditionResponse.ID, 10))
	return nil
}

func resourceSonarqubeQualityGateConditionRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/show"
	query := url.Values{
		"id": []string{strconv.Itoa(d.Get("gateid").(int))},
	}
	sonarQubeURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", sonarQubeURL.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionRead")
		return err
	}
	resp, err := m.(*ProviderConfiguration).httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionRead")
		return err
	}

	defer resp.Body.Close()
	log.WithField("status code", resp.StatusCode).Info("Response from server")
	if resp.StatusCode != http.StatusOK {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}

	getQualityGateConditionResponse := GetQualityGate{}
	err = json.NewDecoder(resp.Body).Decode(&getQualityGateConditionResponse)
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionRead")
	}

	for _, value := range getQualityGateConditionResponse.Conditions {
		if d.Id() == strconv.FormatInt(value.ID, 10) {
			d.SetId(strconv.FormatInt(value.ID, 10))
			d.Set("gateid", getQualityGateConditionResponse.ID)
			d.Set("error", value.Error)
			d.Set("metric", value.Metric)
			d.Set("op", value.OP)
		}
	}

	return nil
}

func resourceSonarqubeQualityGateConditionUpdate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/update_condition"
	query := url.Values{
		"gateid": []string{strconv.Itoa(d.Get("gateid").(int))},
		"id":     []string{d.Id()},
		"error":  []string{strconv.Itoa(d.Get("error").(int))},
		"metric": []string{d.Get("metric").(string)},
		"op":     []string{d.Get("op").(string)},
	}
	sonarQubeURL.RawQuery = query.Encode()

	req, err := http.NewRequest("POST", sonarQubeURL.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionUpdate")
		return err
	}
	resp, err := m.(*ProviderConfiguration).httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionUpdate")
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}

	return resourceSonarqubeQualityGateConditionRead(d, m)
}

func resourceSonarqubeQualityGateConditionDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/delete_condition"
	query := url.Values{
		"id": []string{d.Id()},
	}
	sonarQubeURL.RawQuery = query.Encode()

	req, err := http.NewRequest("POST", sonarQubeURL.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionDelete")
		return err
	}
	resp, err := m.(*ProviderConfiguration).httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionDelete")
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}

	return nil
}
