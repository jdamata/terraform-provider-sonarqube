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

// GetQualityGate for unmarshalling response body of quality gate get
type GetQualityGate struct {
	ID         string                               `json:"id"`
	Name       string                               `json:"name"`
	Conditions []CreateQualityGateConditionResponse `json:"conditions"`
	IsBuiltIn  bool                                 `json:"isBuiltIn"`
	Actions    QualityGateActions                   `json:"actions"`
}

// GetQualityGate for unmarshalling response body of quality gate get for sonarqube versions less than 8
type GetQualityGate_v7 struct {
	ID         int                                     `json:"id"`
	Name       string                                  `json:"name"`
	Conditions []CreateQualityGateConditionResponse_v7 `json:"conditions"`
	IsBuiltIn  bool                                    `json:"isBuiltIn"`
	Actions    QualityGateActions                      `json:"actions"`
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

// CreateQualityGateResponse for unmarshalling response body of quality gate creation for sonarqube versions less than 8
type CreateQualityGateResponse_v7 struct {
	ID   int    `json:"id"`
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
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Name of the quality gate to set as default",
				Default:     false,
				ForceNew:    true,
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

	if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
		// Decode response into struct
		qualityGateResponse := CreateQualityGateResponse_v7{}
		err = json.NewDecoder(resp.Body).Decode(&qualityGateResponse)
		if err != nil {
			return fmt.Errorf("resourceQualityGateCreate: Failed to decode json into struct: %+v", err)
		}
		d.SetId(strconv.Itoa(qualityGateResponse.ID))
	} else {
		// Decode response into struct
		qualityGateResponse := CreateQualityGateResponse{}
		err = json.NewDecoder(resp.Body).Decode(&qualityGateResponse)
		if err != nil {
			return fmt.Errorf("resourceQualityGateCreate: Failed to decode json into struct: %+v", err)
		}
		if d.Get("is_default").(bool) {
			err := setDefaultQualityGate(d, m, d.Get("is_default").(bool))
			if err != nil {
				return err
			}
		}
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

	if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
		// Decode response into struct
		qualityGateReadResponse := GetQualityGate_v7{}
		err = json.NewDecoder(resp.Body).Decode(&qualityGateReadResponse)
		if err != nil {
			return fmt.Errorf("resourceQualityGateRead: Failed to decode json into struct: %+v", err)
		}
		d.SetId(strconv.Itoa(qualityGateReadResponse.ID))
		d.Set("name", qualityGateReadResponse.Name)
		d.Set("is_default", qualityGateReadResponse.Actions.SetAsDefault)
	} else {
		// Decode response into struct
		qualityGateReadResponse := GetQualityGate{}
		err = json.NewDecoder(resp.Body).Decode(&qualityGateReadResponse)
		if err != nil {
			return fmt.Errorf("resourceQualityGateRead: Failed to decode json into struct: %+v", err)
		}
		d.SetId(qualityGateReadResponse.Name)
		d.Set("name", qualityGateReadResponse.Name)
		// Api returns if true if set as default is available. when is_default=true setAsDefault=false so is_default=tue
		d.Set("is_default", !qualityGateReadResponse.Actions.SetAsDefault)
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

		err := setDefaultQualityGate(d, m, false)
		if err != nil {
			return err
		}
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

func setDefaultQualityGate(d *schema.ResourceData, m interface{}, setDefault bool) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/set_as_default"
	if setDefault {
		sonarQubeURL.RawQuery = url.Values{
			"name": []string{d.Get("name").(string)},
		}.Encode()
	} else {
		sonarQubeURL.RawQuery = url.Values{
			"name": []string{"Sonar way"},
		}.Encode()
	}

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"setDefaultQualityGate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
