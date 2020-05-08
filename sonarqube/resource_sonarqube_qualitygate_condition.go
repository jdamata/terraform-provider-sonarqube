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
	url := m.(*ProviderConfiguration).url
	url.Path = "api/qualitygates/create_condition"
	url.ForceQuery = true
	url.RawQuery = fmt.Sprintf("gateId=%v&error=%v&metric=%s&op=%s",
		d.Get("gateid").(int),
		d.Get("error").(int),
		d.Get("metric").(string),
		d.Get("op").(string),
	)

	req, err := http.NewRequest("POST", url.String(), http.NoBody)
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
	url := m.(*ProviderConfiguration).url
	url.Path = "api/qualitygates/show"
	url.ForceQuery = true
	url.RawQuery = fmt.Sprintf("id=%v",
		d.Get("gateid").(int),
	)

	req, err := http.NewRequest("GET", url.String(), http.NoBody)
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
	conditionID, err := strconv.Atoi(d.Id())
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionUpdate")
		return err
	}

	url := m.(*ProviderConfiguration).url
	url.Path = "api/qualitygates/update_condition"
	url.ForceQuery = true
	url.RawQuery = fmt.Sprintf("gateid=%v&id=%v&error=%v&metric=%s&op=%v",
		d.Get("gateid").(int),
		conditionID,
		d.Get("error").(int),
		d.Get("metric").(string),
		d.Get("op").(string),
	)

	req, err := http.NewRequest("POST", url.String(), http.NoBody)
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
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionDelete")
		return err
	}

	url := m.(*ProviderConfiguration).url
	url.Path = "api/qualitygates/delete_condition"
	url.ForceQuery = true
	url.RawQuery = fmt.Sprintf("id=%v",
		id,
	)
	req, err := http.NewRequest("POST", url.String(), http.NoBody)
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
