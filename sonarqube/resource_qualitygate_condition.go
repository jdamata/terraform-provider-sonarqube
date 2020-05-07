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
func qualityGateCondition() *schema.Resource {
	return &schema.Resource{
		Create: qualityGateConditionCreate,
		Read:   qualityGateConditionRead,
		Update: qualityGateConditionUpdate,
		Delete: qualityGateConditionDelete,
		Importer: &schema.ResourceImporter{
			State: qualityGateConditionImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"gateid": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"error": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func qualityGateConditionCreate(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("%s/api/qualitygates/create_condition?gateId=%v&error=%v&metric=%s",
		m.(*ProviderConfiguration).sonarURL,
		d.Get("gateid").(int),
		d.Get("error").(int),
		d.Get("metric").(string),
	)
	req, err := http.NewRequest("POST", url, http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionCreate")
		return err
	}

	req.SetBasicAuth(m.(*ProviderConfiguration).sonarUser, m.(*ProviderConfiguration).sonarPass)
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

func qualityGateConditionRead(d *schema.ResourceData, m interface{}) error {
	url := fmt.Sprintf("%s/api/qualitygates/show?id=%s",
		m.(*ProviderConfiguration).sonarURL,
		d.Id(),
	)
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionRead")
		return err
	}

	req.SetBasicAuth(m.(*ProviderConfiguration).sonarUser, m.(*ProviderConfiguration).sonarPass)
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

	qualityGateConditionResponse := GetQualityGate{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateConditionResponse)
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionRead")
	}

	for _, value := range qualityGateConditionResponse.Conditions {
		if d.Id() == strconv.FormatInt(value.ID, 10) {
			d.SetId(strconv.FormatInt(value.ID, 10))
			d.Set("gateid", value.ID)
			d.Set("error", value.Error)
			d.Set("metric", value.Metric)
		}
	}

	return nil
}

func qualityGateConditionUpdate(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionUpdate")
		return err
	}

	url := fmt.Sprintf("%s/api/qualitygates/update_condition?id=%v&error=%v&metric=%s",
		m.(*ProviderConfiguration).sonarURL,
		id,
		d.Get("error").(int64),
		d.Get("metric").(string),
	)
	req, err := http.NewRequest("POST", url, http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionUpdate")
		return err
	}

	req.SetBasicAuth(m.(*ProviderConfiguration).sonarUser, m.(*ProviderConfiguration).sonarPass)
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

	return qualityGateConditionRead(d, m)
}

func qualityGateConditionDelete(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionDelete")
		return err
	}

	url := fmt.Sprintf("%s/api/qualitygates/delete_condition?id=%v",
		m.(*ProviderConfiguration).sonarURL,
		id,
	)
	req, err := http.NewRequest("POST", url, http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourcequalityGateConditionDelete")
		return err
	}

	req.SetBasicAuth(m.(*ProviderConfiguration).sonarUser, m.(*ProviderConfiguration).sonarPass)
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

func qualityGateConditionImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := qualityGateConditionRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
