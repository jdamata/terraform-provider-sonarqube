package sonarqube

import (
	"encoding/json"
	"errors"
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
	query := url.Values{
		"name":       []string{d.Get("name").(string)},
		"project":    []string{d.Get("project").(string)},
		"visibility": []string{d.Get("visibility").(string)},
	}
	sonarQubeURL.RawQuery = query.Encode()

	req, err := http.NewRequest("POST", sonarQubeURL.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeProjectCreate")
		return err
	}
	resp, err := m.(*ProviderConfiguration).httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeProjectCreate")
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}

	projectResponse := CreateProjectResponse{}
	err = json.NewDecoder(resp.Body).Decode(&projectResponse)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeProjectCreate")
	}

	d.SetId(projectResponse.Project.Key)
	return nil
}

func resourceSonarqubeProjectRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/projects/search"
	query := url.Values{
		"project": []string{d.Id()},
	}
	sonarQubeURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", sonarQubeURL.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeProjectRead")
		return err
	}
	resp, err := m.(*ProviderConfiguration).httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeProjectRead")
		return err
	}

	defer resp.Body.Close()
	log.WithField("status code", resp.StatusCode).Info("Response from server")
	if resp.StatusCode != http.StatusOK {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}

	projectReadResponse := GetProject{}
	err = json.NewDecoder(resp.Body).Decode(&projectReadResponse)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeProjectRead")
	}

	for _, value := range projectReadResponse.Components {
		if d.Id() == value.Key {
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
	query := url.Values{
		"project": []string{d.Id()},
	}
	sonarQubeURL.RawQuery = query.Encode()
	req, err := http.NewRequest("POST", sonarQubeURL.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeProjectDelete")
		return err
	}
	resp, err := m.(*ProviderConfiguration).httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeProjectDelete")
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}

	return nil
}

func resourceSonarqubeProjectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeProjectRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
