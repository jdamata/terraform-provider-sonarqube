package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubeQualityProfileActiveRules() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to list active rules on a Sonarqube quality profile.",
		Read:        dataSourceSonarqubeQualityProfileActiveRulesRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Quality Profile key. Can be obtained through api/qualityprofiles/search.",
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule key.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule name.",
						},
						"severity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule severity.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule type.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule status.",
						},
						"language": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule language.",
						},
						"repo": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule repository.",
						},
					},
				},
				Description: "List of active rules for the quality profile.",
			},
		},
	}
}

func dataSourceSonarqubeQualityProfileActiveRulesRead(d *schema.ResourceData, m interface{}) error {
	rules, err := readQualityProfileRules(d.Get("key").(string), true, m)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s-active", d.Get("key").(string)))
	return d.Set("rules", flattenQualityProfileRules(rules))
}

func readQualityProfileRules(profileKey string, active bool, m interface{}) ([]Rule, error) {
	page := 1
	pageSize := 500
	rules := []Rule{}

	for {
		sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/rules/search"
		sonarQubeURL.RawQuery = url.Values{
			"qprofile":   []string{profileKey},
			"activation": []string{strconv.FormatBool(active)},
			"p":          []string{strconv.Itoa(page)},
			"ps":         []string{strconv.Itoa(pageSize)},
		}.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"GET",
			sonarQubeURL.String(),
			http.StatusOK,
			"readQualityProfileRules",
		)
		if err != nil {
			return nil, err
		}

		ruleReadResponse := GetRule{}
		if err := json.NewDecoder(resp.Body).Decode(&ruleReadResponse); err != nil {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("readQualityProfileRules: Failed to decode json into struct: %+v", err)
		}
		_ = resp.Body.Close()

		rules = append(rules, ruleReadResponse.Rule...)

		if len(rules) >= ruleReadResponse.Total || len(ruleReadResponse.Rule) == 0 {
			break
		}
		page++
	}

	return rules, nil
}

func flattenQualityProfileRules(rules []Rule) []interface{} {
	ruleList := []interface{}{}

	for _, rule := range rules {
		values := map[string]interface{}{
			"key":      rule.RuleKey,
			"name":     rule.Name,
			"severity": rule.Severity,
			"type":     rule.Type,
			"status":   rule.Status,
			"language": rule.Lang,
			"repo":     rule.Repo,
		}

		ruleList = append(ruleList, values)
	}

	return ruleList
}
