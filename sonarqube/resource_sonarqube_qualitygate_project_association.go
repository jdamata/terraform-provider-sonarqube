package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetQualityGateAssociation for unmarshalling response body from getting quality gate association
type GetQualityGateAssociation struct {
	Paging  Paging                              `json:"paging"`
	Results []GetQualityGateAssociationProjects `json:"results"`
}

// GetQualityGateAssociationProjects used in GetQualityGateAssociation
type GetQualityGateAssociationProjects struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	Selected bool   `json:"selected"`
}

// Returns the resource represented by this file.
func resourceSonarqubeQualityGateProjectAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeQualityGateProjectAssociationCreate,
		Read:   resourceSonarqubeQualityGateProjectAssociationRead,
		Delete: resourceSonarqubeQualityGateProjectAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeQualityGateProjectAssociationImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"gateid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"gatename": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"projectkey": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSonarqubeQualityGateProjectAssociationCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/select"
	sonarQubeVersion := m.(*ProviderConfiguration).sonarQubeVersion

	// Sonarqube versions less than 8.0 require gateid instead of gatename
	if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
		sonarQubeURL.RawQuery = url.Values{
			"gateId":     []string{d.Get("gateid").(string)},
			"projectKey": []string{d.Get("projectkey").(string)},
		}.Encode()
	} else {
		sonarQubeURL.RawQuery = url.Values{
			"gateName":   []string{d.Get("gatename").(string)},
			"projectKey": []string{d.Get("projectkey").(string)},
		}.Encode()
	}

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeQualityGateProjectAssociationCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
		id := fmt.Sprintf("%v/%v", d.Get("gateid").(string), d.Get("projectkey").(string))
		d.SetId(id)
	} else {
		id := fmt.Sprintf("%v/%v", d.Get("gatename").(string), d.Get("projectkey").(string))
		d.SetId(id)
	}

	return resourceSonarqubeQualityGateProjectAssociationRead(d, m)
}

func resourceSonarqubeQualityGateProjectAssociationRead(d *schema.ResourceData, m interface{}) error {
	idSlice := strings.Split(d.Id(), "/")
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/search"
	sonarQubeVersion := m.(*ProviderConfiguration).sonarQubeVersion

	// Sonarqube versions less than 8.0 require gateid instead of gatename
	if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
		sonarQubeURL.RawQuery = url.Values{
			"gateId": []string{idSlice[0]},
		}.Encode()
	} else {
		sonarQubeURL.RawQuery = url.Values{
			"gateName": []string{idSlice[0]},
		}.Encode()
	}

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeQualityGateProjectAssociationRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityGateAssociationReadResponse := GetQualityGateAssociation{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateAssociationReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityGateProjectAssociationRead: Failed to decode json into struct: %+v", err)
	}

	// Sonarqube versions less than 8.0 require gateid instead of gatename
	if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
		for _, value := range qualityGateAssociationReadResponse.Results {
			if idSlice[1] == value.Key {
				d.Set("gateid", idSlice[0])
				d.Set("projectkey", value.Key)
			}
		}
	} else {
		for _, value := range qualityGateAssociationReadResponse.Results {
			if idSlice[1] == value.Key {
				d.Set("gatename", idSlice[0])
				d.Set("projectkey", value.Key)
			}
		}
	}

	return nil
}

func resourceSonarqubeQualityGateProjectAssociationDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/deselect"
	sonarQubeURL.RawQuery = url.Values{
		"gateName":   []string{d.Get("gatename").(string)},
		"projectKey": []string{d.Get("projectkey").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeQualityGateProjectAssociationDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeQualityGateProjectAssociationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeQualityProfileRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
