package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// GetQualityProfileRules for unmarshalling response body from getting quality profile association
type GetQualityProfileRules struct {
	Paging  Paging                          `json:"paging"`
	Results []GetQualityProfileRulesResults `json:"results"`
}

// GetQualityProfileRulesResults used in GetQualityProfileRules
type GetQualityProfileRulesResults struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Key      string `json:"key"`
	Selected bool   `json:"selected"`
}

// Returns the resource represented by this file.
func resourceSonarqubeQualityProfileRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeQualityProfileRulesActivate,
		//Read:   resourceSonarqubeQualityProfileRulesRead,
		Delete: resourceSonarqubeQualityProfileRulesDeactivate,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeQualityProfileRulesImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"quality_profile": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Quality profile name",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 100),
				),
			},
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project name",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 100),
				),
			},
			"language": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Quality profile language",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						[]string{"cs", "css", "flex", "go", "java", "js", "jsp", "kotlin", "php", "py", "ruby", "scala", "ts", "vbnet", "web", "xml"},
						false,
					),
				),
			},
		},
	}
}

func resourceSonarqubeQualityProfileRulesActivate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualityprofiles/activate_rule"
	sonarQubeURL.RawQuery = url.Values{
		"key":      []string{d.Get("key").(string)},
		"params":   []string{d.Get("params").(string)},
		"reset":    []string{d.Get("reset").(string)},
		"rule":     []string{d.Get("rule").(string)},
		"severity": []string{d.Get("severity").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeQualityProfileRulesActivate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//id := fmt.Sprintf("%v/%v", d.Get("quality_profile").(string), d.Get("project").(string))
	//d.SetId(id)
	//return resourceSonarqubeQualityProfileRulesRead(d, m)
	return nil
}

func resourceSonarqubeQualityProfileRulesDeactivate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualityprofiles/deactivate_rule"
	sonarQubeURL.RawQuery = url.Values{
		"key":  []string{d.Get("key").(string)},
		"rule": []string{d.Get("rule").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeQualityProfileRulesDeactivate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func resourceSonarqubeQualityProfileRulesRead(d *schema.ResourceData, m interface{}) error {

	// Get all active rules for qualityprofile
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/rules/search"
	sonarQubeURL.RawQuery = url.Values{
		"f":          []string{"isTemplate,name,lang,langName,severity,status,sysTags,tags,templateKey,actives,params"},
		"facets":     []string{"languages,types"},
		"s":          []string{"name"},
		"activation": []string{"true"},
		"qprofile":   []string{d.Get("qprofile").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeQualityProfileRulesRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	getQualityProfileResponse := GetQualityProfileRules{}
	err = json.NewDecoder(resp.Body).Decode(&getQualityProfileResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileRulesRead: Failed to decode json into struct: %+v", err)
	}
	rulesMap := make(map[string]string)
	for _, value := range getQualityProfileResponse.Results {
		rulesMap[value.Key] = value.Name
	}

	d.SetId(d.Id())
	d.Set("rules", rulesMap)
	return nil
}

func resourceSonarqubeQualityProfileRulesImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeQualityProfileRulesRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
