package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Inactives struct {
	QProfile string   `json:"qProfile"`
	Inherit  string   `json:"inherit"`
	Severity string   `json:"severity"`
	Params   []Params `json:"params"`
}

type GetInactiveRules struct {
	Rule      Rule        `json:"rule"`
	Inactives []Inactives `json:"inactives"`
}

func resourceSonarqubeQualityProfileRuleDeactivate() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeQualityProfileRuleDeactivateCreate,
		Delete: resourceSonarqubeQualityProfileRuleDeactivateDelete,
		Read:   resourceSonarqubeQualityProfileRuleDeactivateRead,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeQualityProfileRuleDeactivateImporter,
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

func resourceSonarqubeQualityProfileRuleDeactivateCreate(d *schema.ResourceData, m interface{}) error {
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
		"resourceSonarqubeQualityProfileRuleDeactivateCreate",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileRuleDeactivateCreate: Failed to deactivate rule: %+v", err)
	}
	defer resp.Body.Close()

	d.SetId(d.Get("rule").(string))
	return resourceSonarqubeQualityProfileRuleDeactivateRead(d, m)
}

func resourceSonarqubeQualityProfileRuleDeactivateDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualityprofiles/activate_rule"

	sonarQubeURL.RawQuery = url.Values{
		"key":  []string{d.Get("key").(string)},
		"rule": []string{d.Get("rule").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeQualityProfileRuleDeactivateDelete",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileRuleDeactivateDelete: Failed to activate rule: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeQualityProfileRuleDeactivateRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/rules/show"
	sonarQubeURL.RawQuery = url.Values{
		"key":     []string{d.Id()},
		"actives": []string{"false"},
	}.Encode()
	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeQualityProfileRuleDeactivateRead",
	)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	inactiveRuleReadResponse := GetInactiveRules{}
	err = json.NewDecoder(resp.Body).Decode(&inactiveRuleReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileRuleDeactivateRead: Failed to decode json into struct: %+v", err)
	}

	if d.Id() == inactiveRuleReadResponse.Rule.RuleKey {
		d.SetId(inactiveRuleReadResponse.Rule.RuleKey)
		return nil
	}

	return fmt.Errorf("resourceSonarqubeQualityProfileRuleDeactivateRead: Failed to find project: %+v", d.Id())
}

func resourceSonarqubeQualityProfileRuleDeactivateImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeQualityProfileRuleDeactivateRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
