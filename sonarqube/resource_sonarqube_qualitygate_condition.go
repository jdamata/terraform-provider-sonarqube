package sonarqube

import (
	"encoding/json"
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
	sonarQubeURL.RawQuery = url.Values{
		"gateId": []string{strconv.Itoa(d.Get("gateid").(int))},
		"error":  []string{strconv.Itoa(d.Get("error").(int))},
		"metric": []string{d.Get("metric").(string)},
		"op":     []string{d.Get("op").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourcequalityGateConditionCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityGateConditionResponse := CreateQualityGateConditionResponse{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateConditionResponse)
	if err != nil {
		log.WithError(err).Error("getQualityGateConditionResponse: Failed to decode json into struct")
	}

	d.SetId(strconv.FormatInt(qualityGateConditionResponse.ID, 10))
	return nil
}

func resourceSonarqubeQualityGateConditionRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/show"
	sonarQubeURL.RawQuery = url.Values{
		"id": []string{strconv.Itoa(d.Get("gateid").(int))},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourcequalityGateConditionRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	getQualityGateConditionResponse := GetQualityGate{}
	err = json.NewDecoder(resp.Body).Decode(&getQualityGateConditionResponse)
	if err != nil {
		log.WithError(err).Error("getQualityGateConditionResponse: Failed to decode json into struct")
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
	sonarQubeURL.RawQuery = url.Values{
		"gateid": []string{strconv.Itoa(d.Get("gateid").(int))},
		"id":     []string{d.Id()},
		"error":  []string{strconv.Itoa(d.Get("error").(int))},
		"metric": []string{d.Get("metric").(string)},
		"op":     []string{d.Get("op").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourcequalityGateConditionUpdate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return resourceSonarqubeQualityGateConditionRead(d, m)
}

func resourceSonarqubeQualityGateConditionDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/delete_condition"
	sonarQubeURL.RawQuery = url.Values{
		"id": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourcequalityGateConditionDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
