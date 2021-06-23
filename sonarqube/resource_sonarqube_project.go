package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Project used in CreateProjectResponse
type Project struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	Qualifier string `json:"qualifier"`
}

// GetProject for unmarshalling response body from getting project details
type GetProject struct {
	Paging     Paging              `json:"paging"`
	Components []ProjectComponents `json:"components"`
}

// ProjectComponents used in GetProject
type ProjectComponents struct {
	Organization     string `json:"organization"`
	ProjectKey       string `json:"key"`
	Name             string `json:"name"`
	Qualifier        string `json:"qualifier"`
	Visibility       string `json:"visibility"`
	LastAnalysisDate string `json:"lastAnalysisDate"`
	Revision         string `json:"revision"`
}

// CreateProjectResponse for unmarshalling response body of project creation
type CreateProjectResponse struct {
	Project Project `json:"project"`
}

// Returns the resource represented by this file.
func resourceSonarqubeProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeProjectCreate,
		Read:   resourceSonarqubeProjectRead,
		Delete: resourceSonarqubeProjectDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeProjectImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "public",
				ForceNew: true,
			},
		},
	}
}

func resourceSonarqubeProjectCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/projects/create"
	sonarQubeURL.RawQuery = url.Values{
		"name":       []string{d.Get("name").(string)},
		"project":    []string{d.Get("project").(string)},
		"visibility": []string{d.Get("visibility").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeProjectCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	projectResponse := CreateProjectResponse{}
	err = json.NewDecoder(resp.Body).Decode(&projectResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeProjectCreate: Failed to decode json into struct: %+v", err)
	}

	d.SetId(projectResponse.Project.Key)
	return resourceSonarqubeProjectRead(d, m)
}

func resourceSonarqubeProjectRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/projects/search"
	sonarQubeURL.RawQuery = url.Values{
		"project": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeProjectRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	projectReadResponse := GetProject{}
	err = json.NewDecoder(resp.Body).Decode(&projectReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeProjectRead: Failed to decode json into struct: %+v", err)
	}

	// Loop over all projects to see if the project we need exists.
	for _, value := range projectReadResponse.Components {
		if d.Id() == value.ProjectKey {
			// If it does, set the values of that project
			d.SetId(value.ProjectKey)
			d.Set("name", value.Name)
			d.Set("project", value.ProjectKey)
			d.Set("visibility", value.Visibility)
			return nil
		}
	}

	return fmt.Errorf("resourceSonarqubeProjectRead: Failed to find project: %+v", d.Id())

}

func resourceSonarqubeProjectDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/projects/delete"
	sonarQubeURL.RawQuery = url.Values{
		"project": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeProjectDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeProjectImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeProjectRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
