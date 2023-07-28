package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Webhook struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Url    string `json:"url"`
	Secret string `json:"secret"`
}

type CreateWebhookResponse struct {
	Webhook *Webhook `json:"webhook"`
}

type ListWebhooksResponse struct {
	Webhooks []*Webhook `json:"webhooks"`
}

// Returns the resource represented by this file.
func resourceSonarqubeWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeWebhookCreate,
		Read:   resourceSonarqubeWebhookRead,
		Update: resourceSonarqubeWebhookUpdate,
		Delete: resourceSonarqubeWebhookDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeWebhookImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secret": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSonarqubeWebhookCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/webhooks/create"

	params := url.Values{
		"name": []string{d.Get("name").(string)},
		"url":  []string{d.Get("url").(string)},
	}
	if secret, ok := d.GetOk("secret"); ok {
		params.Set("secret", secret.(string))
	}
	sonarQubeURL.RawQuery = params.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceWebhookCreate",
	)
	if err != nil {
		return fmt.Errorf("resourceWebhookCreate: Failed to call %s: %+v", sonarQubeURL.Path, err)
	}
	defer resp.Body.Close()

	webhookResponse := CreateWebhookResponse{}
	err = json.NewDecoder(resp.Body).Decode(&webhookResponse)
	if err != nil {
		return fmt.Errorf("resourceWebhookCreate: Failed to decode json into struct: %+v", err)
	}

	d.SetId(webhookResponse.Webhook.Key)

	return resourceSonarqubeWebhookRead(d, m)
}

// unfortunately, there doesn't seem to be a way to get a webhook by its ID. the best we can do is list all webhooks and
// loop through the result until we find the one we're looking for.
func resourceSonarqubeWebhookRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/webhooks/list"

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceWebhookRead",
	)
	if err != nil {
		return fmt.Errorf("resourceWebhookRead: Failed to call %s: %+v", sonarQubeURL.Path, err)
	}
	defer resp.Body.Close()

	webhookResponse := ListWebhooksResponse{}
	err = json.NewDecoder(resp.Body).Decode(&webhookResponse)
	if err != nil {
		return fmt.Errorf("resourceWebhookRead: Failed to decode json into struct: %+v", err)
	}

	for _, webhook := range webhookResponse.Webhooks {
		if webhook.Key == d.Id() {
			d.Set("name", webhook.Name)
			d.Set("secret", webhook.Secret)
			d.Set("url", webhook.Url)

			return nil
		}
	}

	return fmt.Errorf("resourceWebhookRead: Failed to find webhook with key %s", d.Id())
}

func resourceSonarqubeWebhookUpdate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/webhooks/update"

	params := url.Values{
		"webhook": []string{d.Id()},
		"name":    []string{d.Get("name").(string)},
		"url":     []string{d.Get("url").(string)},
	}
	if secret, ok := d.GetOk("secret"); ok {
		params.Set("secret", secret.(string))
	}
	sonarQubeURL.RawQuery = params.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceWebhookUpdate",
	)
	if err != nil {
		return fmt.Errorf("resourceWebhookUpdate: Failed to update webhook: %+v", err)
	}
	defer resp.Body.Close()

	return resourceSonarqubeWebhookRead(d, m)
}

func resourceSonarqubeWebhookDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/webhooks/delete"

	sonarQubeURL.RawQuery = url.Values{
		"webhook": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceWebhookDelete",
	)
	if err != nil {
		return fmt.Errorf("resourceWebhookDelete: Failed to delete webhook: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeWebhookImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeWebhookRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
