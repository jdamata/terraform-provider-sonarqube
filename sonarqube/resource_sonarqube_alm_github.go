package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetAlmGithub for unmarshalling response body from alm list definitions. With only github populated
type GetAlmGithub struct {
	Github []struct {
		Key      string `json:"key"`
		URL      string `json:"url"`
		AppID    string `json:"appId"`
		ClientID string `json:"clientId"`
	} `json:"github"`
}

// Returns the resource represented by this file.
func resourceSonarqubeAlmGithub() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeAlmGithubCreate,
		Read:   resourceSonarqubeAlmGithubRead,
		Update: resourceSonarqubeAlmGithubUpdate,
		Delete: resourceSonarqubeAlmGithubDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"appid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"clientid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"clientsecret": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"privatekey": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"webhooksecret": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceSonarqubeAlmGithubCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/create_github"

	sonarQubeURL.RawQuery = url.Values{
		"appId":         []string{d.Get("appid").(string)},
		"clientId":      []string{d.Get("clientid").(string)},
		"clientSecret":  []string{d.Get("clientsecret").(string)},
		"key":           []string{d.Get("key").(string)},
		"privateKey":    []string{d.Get("privatekey").(string)},
		"url":           []string{d.Get("url").(string)},
		"webhookSecret": []string{d.Get("webhooksecret").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeAlmGithubCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d.SetId(d.Get("key").(string))

	return resourceSonarqubeAlmGithubRead(d, m)
}

func resourceSonarqubeAlmGithubRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/list_definitions"
	sonarQubeURL.RawQuery = url.Values{}.Encode() // Dunno if you can keep it empty tbh?

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeAlmGithubRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	AlmGithubReadResponse := GetAlmGithub{}
	err = json.NewDecoder(resp.Body).Decode(&AlmGithubReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeAlmGithubRead: Failed to decode json into struct: %+v", err)
	}
	// Loop over all GitHub instances to see if the Alm instance exists.
	for _, value := range AlmGithubReadResponse.Github {
		if d.Id() == value.Key {
			d.Set("key", value.Key)
			d.Set("url", value.URL)
			d.Set("appid", value.AppID)
			d.Set("clientid", value.ClientID)
			return nil
		}
	}
	return fmt.Errorf("resourceSonarqubeGithubBindingRead: Failed to find github binding: %+v", d.Id())

}
func resourceSonarqubeAlmGithubUpdate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/update_github"
	sonarQubeURL.RawQuery = url.Values{
		"appId":         []string{d.Get("appid").(string)},
		"clientId":      []string{d.Get("clientid").(string)},
		"clientSecret":  []string{d.Get("clientsecret").(string)},
		"key":           []string{d.Id()},
		"newKey":        []string{d.Get("key").(string)},
		"privateKey":    []string{d.Get("privatekey").(string)},
		"url":           []string{d.Get("url").(string)},
		"webhookSecret": []string{d.Get("webhooksecret").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeAlmGithubUpdate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return resourceSonarqubeAlmGithubRead(d, m)
}

func resourceSonarqubeAlmGithubDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/delete"
	sonarQubeURL.RawQuery = url.Values{
		"key": []string{d.Get("key").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeAlmGithubDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}