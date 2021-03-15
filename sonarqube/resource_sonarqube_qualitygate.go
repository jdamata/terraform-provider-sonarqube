package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetQualityGate for unmarshalling response body of quality gate get
type GetQualityGate struct {
	ID         string                               `json:"id"`
	Name       string                               `json:"name"`
	Conditions []CreateQualityGateConditionResponse `json:"conditions"`
	IsBuiltIn  bool                                 `json:"isBuiltIn"`
	Actions    QualityGateActions                   `json:"actions"`
}

// QualityGateActions used in GetQualityGate
type QualityGateActions struct {
	Rename            bool `json:"rename"`
	SetAsDefault      bool `json:"setAsDefault"`
	Copy              bool `json:"copy"`
	AssociateProjects bool `json:"associateProjects"`
	Delete            bool `json:"delete"`
	ManageConditions  bool `json:"manageConditions"`
}

// CreateQualityGateResponse for unmarshalling response body of quality gate creation
type CreateQualityGateResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

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
	sonarQubeVersion := m.(*ProviderConfiguration).sonarQubeVersion
	sonarQubeURL.RawQuery = url.Values{
		"name": []string{d.Get("name").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceQualityGateCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityGateResponse := CreateQualityGateResponse{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateResponse)
	if err != nil {
		return fmt.Errorf("resourceQualityGateCreate: Failed to decode json into struct: %+v", err)
	}

	if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
		d.SetId(qualityGateResponse.Name)
	} else {
		d.SetId(qualityGateResponse.Name)
	}

	return resourceSonarqubeQualityGateRead(d, m)
}

func resourceSonarqubeQualityGateRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/show"
	sonarQubeVersion := m.(*ProviderConfiguration).sonarQubeVersion

	// Sonarqube versions less than 8.0 require id instead of name
	if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
		sonarQubeURL.RawQuery = url.Values{
			"id": []string{d.Id()},
		}.Encode()
	} else {
		sonarQubeURL.RawQuery = url.Values{
			"name": []string{d.Id()},
		}.Encode()
	}

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceQualityGateRead",
	)
	if err != nil {
		return fmt.Errorf("resourceQualityGateRead: Failed to call api/qualitygates/show: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityGateReadResponse := GetQualityGate{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateReadResponse)
	if err != nil {
		return fmt.Errorf("resourceQualityGateRead: Failed to decode json into struct: %+v", err)
	}

	if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
		d.SetId(qualityGateReadResponse.Name)
		d.Set("name", qualityGateReadResponse.Name)
	} else {
		d.SetId(qualityGateReadResponse.Name)
		d.Set("name", qualityGateReadResponse.Name)
	}
	return nil
}

func resourceSonarqubeQualityGateDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/destroy"
	sonarQubeVersion := m.(*ProviderConfiguration).sonarQubeVersion

	// Sonarqube versions less than 8.0 require id instead of name
	if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
		sonarQubeURL.RawQuery = url.Values{
			"id": []string{d.Id()},
		}.Encode()
	} else {
		sonarQubeURL.RawQuery = url.Values{
			"name": []string{d.Id()},
		}.Encode()
	}

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceQualityGateDelete",
	)
	if err != nil {
		return fmt.Errorf("resourceQualityGateDelete: Failed to delete quality gate: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeQualityGateImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeQualityGateRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
