package sonarqube

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Returns the resource represented by this file.
func resourceSonarqubePlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubePluginCreate,
		Read:   resourceSonarqubePluginRead,
		Delete: resourceSonarqubePluginDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubePluginImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSonarqubePluginCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/plugins/install"
	sonarQubeURL.RawQuery = url.Values{
		"key": []string{d.Get("key").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubePluginCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d.SetId(d.Get("key").(string))
	return nil
}

func resourceSonarqubePluginRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/plugins/installed"

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubePluginRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	getInstalledPlugins := GetInstalledPlugins{}
	err = json.NewDecoder(resp.Body).Decode(&getInstalledPlugins)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubePluginRead: Failed to decode json into struct")
	}

	// Loop over all projects to see if the project we need exists.
	for _, value := range getInstalledPlugins.Plugins {
		if d.Id() == value.Key {
			// If it does, set the values of that project
			d.SetId(value.Key)
			d.Set("key", value.Key)
		}
	}

	return nil
}

func resourceSonarqubePluginDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/plugins/uninstall"
	sonarQubeURL.RawQuery = url.Values{
		"key": []string{d.Id()},
	}.Encode()

	log.Error(sonarQubeURL.String())
	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubePluginDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubePluginImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubePluginRead(d, m); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
