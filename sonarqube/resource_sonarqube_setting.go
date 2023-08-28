package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Setting struct {
	Key         string              `json:"key"`
	Value       string              `json:"value"`
	Values      []string            `json:"values"`
	Inherited   bool                `json:"inherited"`
	FieldValues []map[string]string `json:"fieldValues"`
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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Setting value. To reset a value, please use the reset web service.",
				ExactlyOneOf: []string{"value", "values", "field_values"},
			},
			"values": {
				Type:         schema.TypeList,
				Optional:     true,
				Description:  "Setting multi values for the supplied key",
				ExactlyOneOf: []string{"value", "values", "field_values"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"field_values": {
				Type:         schema.TypeList,
				Optional:     true,
				Description:  "Setting field values for the supplied key",
				ExactlyOneOf: []string{"value", "values", "field_values"},
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: schema.TypeString,
				},
			},
			"component": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Component key. Only keys for projects, applications, portfolios or subportfolios are accepted.",
			},
		},
	}
}

func resourceSonarqubeSettingsCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/settings/set"
	sonarQubeURL.RawQuery = getCreateOrUpdateQueryRawQuery([]string{d.Get("key").(string)}, d)

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

	d.SetId(buildId(d))
	return resourceSonarqubeSettingsRead(d, m)
}

func resourceSonarqubeSettingsRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/settings/values"

	key := d.Get("key").(string)
	params := url.Values{
		"keys": []string{key},
	}
	component, componentOk := d.GetOk("component")
	if componentOk {
		params.Set("component", component.(string))
	}

	sonarQubeURL.RawQuery = params.Encode()

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
		if key == value.Key {
			d.Set("key", value.Key)
			d.Set("value", value.Value)
			d.Set("values", value.Values)
			d.Set("field_values", value.FieldValues)
			// Field 'component' is not included in the response object, so it is imported from the parameter.
			if componentOk {
				d.Set("component", component.(string))
			}
			d.SetId(buildId(d))
			return nil
		}
	}
	return fmt.Errorf("resourceSonarqubeSettingsRead: Failed to find setting: %+v", d.Id())
}

func resourceSonarqubeSettingsDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/settings/reset"
	params := url.Values{
		"keys": []string{d.Get("key").(string)},
	}
	component, componentOk := d.GetOk("component")
	if componentOk {
		params.Set("component", component.(string))
	}

	sonarQubeURL.RawQuery = params.Encode()

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
	// the ID consists of the key and the component (foo/bar)
	idSplit := strings.Split(d.Id(), "/")
	d.Set("key", idSplit[0])
	if len(idSplit) > 1 {
		d.Set("component", idSplit[1])
	}
	if err := resourceSonarqubeSettingsRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func resourceSonarqubeSettingsUpdate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/settings/set"
	sonarQubeURL.RawQuery = getCreateOrUpdateQueryRawQuery([]string{d.Get("key").(string)}, d)

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

	return resourceSonarqubeSettingsRead(d, m)
}
func getCreateOrUpdateQueryRawQuery(key []string, d *schema.ResourceData) string {
	// build the base query
	RawQuery := url.Values{
		"key": key,
	}
	// check optional
	if component, ok := d.GetOk("component"); ok {
		RawQuery.Add("component", component.(string))
	}
	// Add in value/values/fieldValues as appropriate
	// single value
	if value, ok := d.GetOk("value"); ok {
		RawQuery.Add("value", value.(string))
	} else {
		// array of strings
		if values, ok := d.GetOk("values"); ok {
			for _, value := range values.([]interface{}) {
				RawQuery.Add("values", value.(string))
			}
		} else {
			// array of objects of key/value pairs
			fieldValues := d.Get("field_values").([]interface{})
			for _, value := range fieldValues {
				b, _ := json.Marshal(value)
				fv := string(b)
				RawQuery.Add("fieldValues", fv)
			}
		}
	}
	return RawQuery.Encode()
}

func buildId(d *schema.ResourceData) string {
	component, ok := d.GetOk("component")
	// the ID consists of the key and the component (foo/bar)
	if ok {
		return fmt.Sprintf("%s/%s", d.Get("key").(string), component.(string))
	} else {
		return fmt.Sprintf("%s", d.Get("key").(string))
	}
}
