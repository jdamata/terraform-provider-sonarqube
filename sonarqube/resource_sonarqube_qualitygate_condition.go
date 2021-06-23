package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CreateQualityGateConditionResponse for unmarshalling response body of condition creation
type CreateQualityGateConditionResponse struct {
	ID      string `json:"id"`
	Metric  string `json:"metric"`
	OP      string `json:"op"`
	Error   string `json:"error"`
	Warning string `json:"warning"`
}

// CreateQualityGateConditionResponse for unmarshalling response body of condition creation
type CreateQualityGateConditionResponse_v7 struct {
	ID      int    `json:"id"`
	Metric  string `json:"metric"`
	OP      string `json:"op"`
	Error   string `json:"error"`
	Warning string `json:"warning"`
}

// Returns the resource represented by this file.
func resourceSonarqubeQualityGateCondition() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeQualityGateConditionCreate,
		Read:   resourceSonarqubeQualityGateConditionRead,
		Update: resourceSonarqubeQualityGateConditionUpdate,
		Delete: resourceSonarqubeQualityGateConditionDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"gatename": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"gateid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"threshold": {
				Type:     schema.TypeString,
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
	sonarQubeVersion := m.(*ProviderConfiguration).sonarQubeVersion

	// Sonarqube versions less than 8.0 require gateid instead of gatename
	if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
		sonarQubeURL.RawQuery = url.Values{
			"gateId": []string{d.Get("gateid").(string)},
			"error":  []string{d.Get("threshold").(string)},
			"metric": []string{d.Get("metric").(string)},
			"op":     []string{d.Get("op").(string)},
		}.Encode()
	} else {
		sonarQubeURL.RawQuery = url.Values{
			"gateName": []string{d.Get("gatename").(string)},
			"error":    []string{d.Get("threshold").(string)},
			"metric":   []string{d.Get("metric").(string)},
			"op":       []string{d.Get("op").(string)},
		}.Encode()
	}

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

	// Sonarqube versions less than 8.0 require gateid instead of gatename
	if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
		// Decode response into struct
		qualityGateConditionResponse := CreateQualityGateConditionResponse_v7{}
		err = json.NewDecoder(resp.Body).Decode(&qualityGateConditionResponse)
		if err != nil {
			return fmt.Errorf("createQualityGateConditionResponse: Failed to decode json into struct: %+v", err)
		}
		d.SetId(strconv.Itoa(qualityGateConditionResponse.ID))
	} else {
		// Decode response into struct
		qualityGateConditionResponse := CreateQualityGateConditionResponse{}
		err = json.NewDecoder(resp.Body).Decode(&qualityGateConditionResponse)
		if err != nil {
			return fmt.Errorf("createQualityGateConditionResponse: Failed to decode json into struct: %+v", err)
		}
		d.SetId(qualityGateConditionResponse.ID)
	}

	return resourceSonarqubeQualityGateConditionRead(d, m)
}

func resourceSonarqubeQualityGateConditionRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/show"
	sonarQubeVersion := m.(*ProviderConfiguration).sonarQubeVersion

	// Sonarqube versions less than 8.0 require gateid instead of gatename
	if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
		sonarQubeURL.RawQuery = url.Values{
			"id": []string{d.Get("gateid").(string)},
		}.Encode()
	} else {
		sonarQubeURL.RawQuery = url.Values{
			"name": []string{d.Get("gatename").(string)},
		}.Encode()
	}

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

	// Sonarqube versions less than 8.0 require gateid instead of gatename
	if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
		// Decode response into struct
		getQualityGateConditionResponse := GetQualityGate_v7{}
		err = json.NewDecoder(resp.Body).Decode(&getQualityGateConditionResponse)
		if err != nil {
			return fmt.Errorf("getQualityGateConditionResponse: Failed to decode json into struct: %+v", err)
		}
		for _, value := range getQualityGateConditionResponse.Conditions {
			if d.Id() == strconv.Itoa(value.ID) {
				d.SetId(strconv.Itoa(value.ID))
				d.Set("gateid", strconv.Itoa(getQualityGateConditionResponse.ID))
				d.Set("threshold", value.Error)
				d.Set("metric", value.Metric)
				d.Set("op", value.OP)
				return nil
			}
		}
	} else {
		// Decode response into struct
		getQualityGateConditionResponse := GetQualityGate{}
		err = json.NewDecoder(resp.Body).Decode(&getQualityGateConditionResponse)
		if err != nil {
			return fmt.Errorf("getQualityGateConditionResponse: Failed to decode json into struct: %+v", err)
		}
		for _, value := range getQualityGateConditionResponse.Conditions {
			if d.Id() == value.ID {
				d.SetId(value.ID)
				d.Set("gatename", getQualityGateConditionResponse.Name)
				d.Set("threshold", value.Error)
				d.Set("metric", value.Metric)
				d.Set("op", value.OP)
				return nil
			}
		}
	}

	return fmt.Errorf("getQualityGateConditionResponse: Failed to find quality gate condition: %+v", d.Id())
}

func resourceSonarqubeQualityGateConditionUpdate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/update_condition"
	sonarQubeURL.RawQuery = url.Values{
		"id":     []string{d.Id()},
		"error":  []string{d.Get("threshold").(string)},
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
