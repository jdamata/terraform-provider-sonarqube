package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSonarqubeSecuredSetting() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Sonarqube Secured Setting resource. This can be used to manage Sonarqube secured (encrypted) settings. The value is write-only and cannot be read back from the API.",
		Create:      resourceSonarqubeSecuredSettingCreate,
		Read:        resourceSonarqubeSecuredSettingRead,
		Update:      resourceSonarqubeSecuredSettingUpdate,
		Delete:      resourceSonarqubeSecuredSettingDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeSecuredSettingImport,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Setting key.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Setting value. Stored encrypted server-side and never returned by the API (write-only).",
			},
			"component": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Component (project) key to scope the setting to. If omitted, the setting is global.",
			},
		},
	}
}

func resourceSonarqubeSecuredSettingCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/settings/set"

	params := url.Values{
		"key":   []string{d.Get("key").(string)},
		"value": []string{d.Get("value").(string)},
	}
	if component, ok := d.GetOk("component"); ok {
		params.Set("component", component.(string))
	}
	sonarQubeURL.RawQuery = params.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeSecuredSettingCreate",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeSecuredSettingCreate: Failed to call %s: %+v", sonarQubeURL.Path, err)
	}
	defer resp.Body.Close()

	key := d.Get("key").(string)
	if component, ok := d.GetOk("component"); ok {
		d.SetId(key + ":" + component.(string))
	} else {
		d.SetId(key)
	}

	return resourceSonarqubeSecuredSettingRead(d, m)
}

func resourceSonarqubeSecuredSettingRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/settings/values"

	key := d.Get("key").(string)
	params := url.Values{
		"keys": []string{key},
	}
	if component, ok := d.GetOk("component"); ok {
		params.Set("component", component.(string))
	}
	sonarQubeURL.RawQuery = params.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeSecuredSettingRead",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeSecuredSettingRead: Failed to call %s: %+v", sonarQubeURL.Path, err)
	}
	defer resp.Body.Close()

	settingReadResponse := GetSettings{}
	err = json.NewDecoder(resp.Body).Decode(&settingReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeSecuredSettingRead: Failed to decode json into struct: %+v", err)
	}

	for _, securedKey := range settingReadResponse.SetSecuredSettings {
		if securedKey == key {
			// Value intentionally not set — secured settings are write-only.
			return nil
		}
	}

	// Key not found in setSecuredSettings — was deleted out-of-band.
	d.SetId("")
	return nil
}

func resourceSonarqubeSecuredSettingUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceSonarqubeSecuredSettingCreate(d, m)
}

func resourceSonarqubeSecuredSettingDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/settings/reset"

	key := d.Get("key").(string)
	params := url.Values{
		"keys": []string{key},
	}
	if component, ok := d.GetOk("component"); ok {
		params.Set("component", component.(string))
	}
	sonarQubeURL.RawQuery = params.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeSecuredSettingDelete",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeSecuredSettingDelete: Failed to delete secured setting: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeSecuredSettingImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// import id in format {key} or {key}:{component}
	parts := strings.SplitN(d.Id(), ":", 2)

	if err := d.Set("key", parts[0]); err != nil {
		return nil, fmt.Errorf("resourceSonarqubeSecuredSettingImport: Failed to set key: %+v", err)
	}
	if len(parts) == 2 {
		if err := d.Set("component", parts[1]); err != nil {
			return nil, fmt.Errorf("resourceSonarqubeSecuredSettingImport: Failed to set component: %+v", err)
		}
	}

	if err := resourceSonarqubeSecuredSettingRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
