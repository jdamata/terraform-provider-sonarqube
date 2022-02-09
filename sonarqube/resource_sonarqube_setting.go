package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Setting struct {
	Key         string        `json:"key"`
	Value       string        `json:"value"`
	Values      []string      `json:"values"`
	Inherited   bool          `json:"inherited"`
	FieldValues []FieldValues `json:"fieldValues"`
}

type FieldValues struct {
	Boolean string `json:"boolean"`
	Text    string `json:"text"`
}

type GetSettings struct {
	Setting            []Setting `json:"settings"`
	SetSecuredSettings []string  `json:"setSecuredSettings"`
}

func resourceSonarqubeSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeSettingsCreate,
		Read:   resourceSonarqubeSettingsRead,
		Update: resourceSonarqubeSettingsUpdate,
		Delete: resourceSonarqubeSettingsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeSettingsImporter,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Setting key",
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Setting value. To reset a value, please use the reset web service.",
				// ExactlyOneOf: []string{"value", "values"},
			},
			// TODO: add support for values
			//"values": {
			//	Type:        schema.TypeList,
			//	Optional:    true,
			//	Description: "Setting multi value. To set several values, the parameter must be called once for each value",
			//  ExactlyOneOf: []string{"value", "values"},
			//},
			//"fieldValues": {
			//	Type:        schema.TypeString,
			//	Optional:    true,
			//	Description: "Setting field values. To set several values, the parameter must be called once for each value.",
			//},
		},
	}
}

func resourceSonarqubeSettingsCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/settings/set"
	sonarQubeURL.RawQuery = url.Values{
		"key":   []string{d.Get("key").(string)},
		"value": []string{d.Get("value").(string)},
		// "fieldValues": []string{d.Get("fieldValues").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeSettingsCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d.SetId(d.Get("key").(string))
	return resourceSonarqubeRuleRead(d, m)
}

func resourceSonarqubeSettingsRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/settings/values"
	sonarQubeURL.RawQuery = url.Values{
		"keys": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeSettingsRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	settingReadResponse := GetSettings{}
	err = json.NewDecoder(resp.Body).Decode(&settingReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeSettingsRead: Failed to decode json into struct: %+v", err)
	}

	for _, value := range settingReadResponse.Setting {
		if d.Id() == value.Key {
			d.SetId(value.Key)
			d.Set("key", value.Key)
			d.Set("value", value.Value)
			// d.Set("fieldValues", value.FieldValues)
			return nil
		}
	}
	return fmt.Errorf("resourceSonarqubeSettingsRead: Failed to find setting: %+v", d.Id())
}

func resourceSonarqubeSettingsDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/settings/reset"
	sonarQubeURL.RawQuery = url.Values{
		"keys": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeSettingsDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeSettingsImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeSettingsRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func resourceSonarqubeSettingsUpdate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/settings/set"
	sonarQubeURL.RawQuery = url.Values{
		"key":   []string{d.Id()},
		"value": []string{d.Get("value").(string)},
		// "fieldValues": []string{d.Get("fieldValues").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeSettingsUpdate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return resourceSonarqubeRuleRead(d, m)
}
