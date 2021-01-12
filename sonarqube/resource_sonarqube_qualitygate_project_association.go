package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
			"gatename": {
				Type:     schema.TypeString,
				Required: true,
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
	sonarQubeURL.RawQuery = url.Values{
		"gateName":   []string{d.Get("gatename").(string)},
		"projectKey": []string{d.Get("projectkey").(string)},
	}.Encode()

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

	id := fmt.Sprintf("%v/%v", d.Get("gatename").(string), d.Get("projectkey").(string))
	d.SetId(id)
	return resourceSonarqubeQualityGateProjectAssociationRead(d, m)
}

func resourceSonarqubeQualityGateProjectAssociationRead(d *schema.ResourceData, m interface{}) error {
	idSlice := strings.Split(d.Id(), "/")
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/search"
	sonarQubeURL.RawQuery = url.Values{
		"gateName": []string{idSlice[0]},
	}.Encode()

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

	for _, value := range qualityGateAssociationReadResponse.Results {
		if idSlice[1] == value.Key {
			d.Set("gatename", idSlice[0])
			d.Set("projectkey", value.Key)
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
