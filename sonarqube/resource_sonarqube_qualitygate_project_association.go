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
	QualityGate struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Default bool   `json:"default"`
	} `json:"qualityGate"`
}

// Returns the resource represented by this file.
func resourceSonarqubeQualityGateProjectAssociation() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Sonarqube Quality Gate Project association resource. This can be used to associate a Quality Gate to a Project",
		Create:      resourceSonarqubeQualityGateProjectAssociationCreate,
		Read:        resourceSonarqubeQualityGateProjectAssociationRead,
		Delete:      resourceSonarqubeQualityGateProjectAssociationDelete,
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
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the Quality Gate",
			},
			"projectkey": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key of the project. Maximum length 400. All letters, digits, dash, underscore, period or colon.",
			},
		},
	}
}

func resourceSonarqubeQualityGateProjectAssociationCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/select"

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
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/get_by_project"

	sonarQubeURL.RawQuery = url.Values{
		"project": []string{idSlice[1]},
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

	d.Set("projectkey", idSlice[1])
	d.Set("gatename", qualityGateAssociationReadResponse.QualityGate.Name)
	return nil
}

func resourceSonarqubeQualityGateProjectAssociationDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/deselect"

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
	if err := resourceSonarqubeQualityGateProjectAssociationRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
