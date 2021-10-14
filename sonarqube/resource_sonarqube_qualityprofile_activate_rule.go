package sonarqube

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSonarqubeQualityProfileRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeQualityProfileRuleCreate,
		Delete: resourceSonarqubeQualityProfileRuleDelete,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Quality Profile key. Can be obtained through api/qualityprofiles/search",
			},
			"params": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parameters as semi-colon list of =, for example 'params=key1=v1;key2=v2' (Only for custom rule)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"reset": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Reset severity and parameters of activated rule. Set the values defined on parent profile or from rule default values.",
				Default:     false,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						[]string{"true", "false", "yes", "no"},
						false,
					),
				),
			},
			"rule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule key",
			},
			"severity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule severity",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						[]string{"INFO", "MINOR", "MAJOR", "CRITICAL", "BLOCKER"},
						false,
					),
				),
			},
		},
	}
}

func resourceSonarqubeQualityProfileRuleCreate(d *schema.ResourceData, m interface{}) error {
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
		http.StatusOK,
		"resourceSonarqubeQualityProfileRuleCreate",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileRuleCreate: Failed to delete quality profile: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeQualityProfileRuleDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualityprofiles/activate_rule"
	sonarQubeURL.RawQuery = url.Values{
		"key":  []string{d.Get("key").(string)},
		"rule": []string{d.Get("rule").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeQualityProfileRuleDelete",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileRuleDelete: Failed to delete quality profile: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}
