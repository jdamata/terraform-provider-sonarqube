package sonarqube

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var rulesSearchFilters = map[string]string{
	"languages":    "languages",
	"repositories": "repositories",
	"severities":   "severities",
	"statuses":     "statuses",
	"tags":         "tags",
	"types":        "types",
}

func dataSourceSonarqubeRules() *schema.Resource {
	filterSchema := func(description string) *schema.Schema {
		return &schema.Schema{
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: description,
		}
	}
	return &schema.Resource{
		Description: "Use this data source to list Sonarqube rules matching a set of filters, for example to activate them on a quality profile with `for_each`. All filters are optional and combine with AND; values within one filter combine with OR. SonarQube returns at most 10000 rules for one search, so use filters to narrow large rule sets.",
		Read:        dataSourceSonarqubeRulesRead,
		Schema: map[string]*schema.Schema{
			"search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Free-text search, the same as the search box on the SonarQube rules page.",
			},
			"languages":    filterSchema("Language keys, e.g. `java`, `cs`, `js`."),
			"repositories": filterSchema("Rule repository keys, e.g. `java`, `csharpsquid`."),
			"severities":   filterSchema("Severities: `INFO`, `MINOR`, `MAJOR`, `CRITICAL`, `BLOCKER`."),
			"statuses":     filterSchema("Statuses: `READY`, `BETA`, `DEPRECATED`, `REMOVED`."),
			"tags":         filterSchema("Rule tags."),
			"types":        filterSchema("Types: `CODE_SMELL`, `BUG`, `VULNERABILITY`, `SECURITY_HOTSPOT`."),
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key":      {Type: schema.TypeString, Computed: true, Description: "Rule key."},
						"name":     {Type: schema.TypeString, Computed: true, Description: "Rule name."},
						"severity": {Type: schema.TypeString, Computed: true, Description: "Rule severity."},
						"type":     {Type: schema.TypeString, Computed: true, Description: "Rule type."},
						"status":   {Type: schema.TypeString, Computed: true, Description: "Rule status."},
						"language": {Type: schema.TypeString, Computed: true, Description: "Rule language."},
						"repo":     {Type: schema.TypeString, Computed: true, Description: "Rule repository."},
					},
				},
				Description: "The rules matching the filters.",
			},
		},
	}
}

func dataSourceSonarqubeRulesRead(d *schema.ResourceData, m interface{}) error {
	filters := url.Values{}
	if search, ok := d.GetOk("search"); ok {
		filters.Set("q", search.(string))
	}
	for attribute, param := range rulesSearchFilters {
		if v, ok := d.GetOk(attribute); ok {
			values := []string{}
			for _, value := range v.(*schema.Set).List() {
				values = append(values, value.(string))
			}
			sort.Strings(values)
			filters.Set(param, strings.Join(values, ","))
		}
	}

	rules, err := searchRules(filters, m)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", schema.HashString(filters.Encode())))
	return d.Set("rules", flattenQualityProfileRules(rules))
}
