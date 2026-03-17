package sonarqube

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceSonarqubeQualityProfileDeactivatedRules() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to list deactivated rules on a Sonarqube quality profile.",
		Read:        dataSourceSonarqubeQualityProfileDeactivatedRulesRead,
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
				Description: "List of deactivated rules for the quality profile.",
			},
		},
	}
}

func dataSourceSonarqubeQualityProfileDeactivatedRulesRead(d *schema.ResourceData, m interface{}) error {
	rules, err := readQualityProfileRules(d.Get("key").(string), false, m)
	if err != nil {
		return err
	}

	d.SetId(d.Get("key").(string) + "-deactivated")
	return d.Set("rules", flattenQualityProfileRules(rules))
}
