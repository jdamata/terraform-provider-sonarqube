package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceSonarqubeQualityProfileDeactivateRule manages deactivating a rule on a quality profile.
func resourceSonarqubeQualityProfileDeactivateRule() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Sonarqube Rules resource. This can be used to deactivate Sonarqube rules on a quality profile.",
		Create:      resourceSonarqubeQualityProfileDeactivateRuleCreate,
		Delete:      resourceSonarqubeQualityProfileDeactivateRuleDelete,
		Read:        resourceSonarqubeQualityProfileDeactivateRuleRead,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeQualityProfileDeactivateRuleImporter,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Quality Profile key. Can be obtained through api/qualityprofiles/search",
			},
			"rule": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule key",
			},
		},
	}
}

func resourceSonarqubeQualityProfileDeactivateRuleCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualityprofiles/deactivate_rule"
	sonarQubeURL.RawQuery = url.Values{
		"key":  []string{d.Get("key").(string)},
		"rule": []string{d.Get("rule").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeQualityProfileDeactivateRuleCreate",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileDeactivateRuleCreate: Failed to deactivate quality profile rule: %+v", err)
	}
	defer resp.Body.Close()

	d.SetId(d.Get("rule").(string))
	return resourceSonarqubeQualityProfileDeactivateRuleRead(d, m)
}

func resourceSonarqubeQualityProfileDeactivateRuleDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualityprofiles/activate_rule"
	sonarQubeURL.RawQuery = url.Values{
		"key":   []string{d.Get("key").(string)},
		"rule":  []string{d.Get("rule").(string)},
		"reset": []string{"true"},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeQualityProfileDeactivateRuleDelete",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileDeactivateRuleDelete: Failed to reactivate quality profile rule: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeQualityProfileDeactivateRuleRead(d *schema.ResourceData, m interface{}) error {
	ruleKey := d.Get("rule").(string)
	if ruleKey == "" {
		ruleKey = d.Id()
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/rules/show"
	sonarQubeURL.RawQuery = url.Values{
		"key":     []string{ruleKey},
		"actives": []string{"true"},
	}.Encode()
	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeQualityProfileDeactivateRuleRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	activeRuleReadResponse := GetActiveRules{}
	err = json.NewDecoder(resp.Body).Decode(&activeRuleReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileDeactivateRuleRead: Failed to decode json into struct: %+v", err)
	}

	profileKey := d.Get("key").(string)

	if activeRuleReadResponse.Rule.RuleKey == "" {
		d.SetId("")
		return nil
	}

	for _, active := range activeRuleReadResponse.Actives {
		if active.QProfile == profileKey {
			// Rule is still active on the profile, resource should be recreated.
			d.SetId("")
			return nil
		}
	}

	if err := d.Set("rule", activeRuleReadResponse.Rule.RuleKey); err != nil {
		return err
	}

	d.SetId(activeRuleReadResponse.Rule.RuleKey)
	return nil
}

func resourceSonarqubeQualityProfileDeactivateRuleImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeQualityProfileDeactivateRuleRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
