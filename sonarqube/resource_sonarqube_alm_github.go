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
		Description: `Provides a Sonarqube GitHub Alm/Devops Platform Integration resource. This can be used to create and manage a Alm/Devops
Platform Integration for GitHub.`,
		Create: resourceSonarqubeAlmGithubCreate,
		Read:   resourceSonarqubeAlmGithubRead,
		Update: resourceSonarqubeAlmGithubUpdate,
		Delete: resourceSonarqubeAlmGithubDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "GitHub App ID. Maximum length: 80",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "GitHub App Client ID. Maximum length: 80",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "GitHub App Client Secret. Maximum length: 160",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique key of the GitHUb instance setting. Maximum length: 200",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "GitHub App private key. Maximum length: 2500",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "GitHub API URL. Maximum length: 2000",
			},
			"webhook_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "GitHub App Webhook Secret. Maximum length: 160",
			},
		},
	}
}

func resourceSonarqubeAlmGithubCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/create_github"

	sonarQubeURL.RawQuery = url.Values{
		"appId":         []string{d.Get("app_id").(string)},
		"clientId":      []string{d.Get("client_id").(string)},
		"clientSecret":  []string{d.Get("client_secret").(string)},
		"key":           []string{d.Get("key").(string)},
		"privateKey":    []string{d.Get("private_key").(string)},
		"url":           []string{d.Get("url").(string)},
		"webhookSecret": []string{d.Get("webhook_secret").(string)},
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
			d.Set("app_id", value.AppID)
			d.Set("client_id", value.ClientID)
			return nil
		}
	}
	return fmt.Errorf("resourceSonarqubeGithubBindingRead: Failed to find github binding: %+v", d.Id())
}

func resourceSonarqubeAlmGithubUpdate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/update_github"
	sonarQubeURL.RawQuery = url.Values{
		"appId":         []string{d.Get("app_id").(string)},
		"clientId":      []string{d.Get("client_id").(string)},
		"clientSecret":  []string{d.Get("client_secret").(string)},
		"key":           []string{d.Id()},
		"newKey":        []string{d.Get("key").(string)},
		"privateKey":    []string{d.Get("private_key").(string)},
		"url":           []string{d.Get("url").(string)},
		"webhookSecret": []string{d.Get("webhook_secret").(string)},
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
