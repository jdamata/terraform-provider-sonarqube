package sonarqube

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

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

	resp := httpRequestHelper(
		*m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeProjectCreate",
	)
	defer resp.Body.Close()

	// Decode response into struct
	projectResponse := CreateProjectResponse{}
	err := json.NewDecoder(resp.Body).Decode(&projectResponse)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeProjectCreate: Failed to decode json into struct")
	}

	d.SetId(projectResponse.Project.Key)
	return nil
}

func resourceSonarqubeProjectRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/projects/search"
	sonarQubeURL.RawQuery = url.Values{
		"project": []string{d.Id()},
	}.Encode()

	resp := httpRequestHelper(
		*m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeProjectRead",
	)

	defer resp.Body.Close()

	// Decode response into struct
	projectReadResponse := GetProject{}
	err := json.NewDecoder(resp.Body).Decode(&projectReadResponse)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeProjectRead: Failed to decode json into struct")
	}

	// Loop over all projects to see if the project we need exists.
	for _, value := range projectReadResponse.Components {
		if d.Id() == value.Key {
			// If it does, set the values of that project
			d.SetId(value.Key)
			d.Set("name", value.Name)
			d.Set("key", value.Key)
			d.Set("visibility", value.Visibility)
		}
	}

	return nil
}

func resourceSonarqubeProjectDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/projects/delete"
	sonarQubeURL.RawQuery = url.Values{
		"project": []string{d.Id()},
	}.Encode()

	resp := httpRequestHelper(
		*m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeProjectDelete",
	)

	defer resp.Body.Close()
	return nil
}

func resourceSonarqubeProjectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeProjectRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
