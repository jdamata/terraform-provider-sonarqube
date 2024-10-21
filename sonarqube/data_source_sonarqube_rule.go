package sonarqube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubeRule() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a Sonarqube rule resource",
		Read:        dataSourceSonarqubeRuleRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The key of the sonarqube rule. Should be <repo>:<name>. https://next.sonarqube.com/sonarqube/web_api/api/rules?query=api%2Frules%2Fcreate",
			},
			"markdown_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule description",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule name",
			},
			"severity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule severity",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule status",
			},
			"template_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Key of the template rule",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule type",
			},
		},
	}
}

func dataSourceSonarqubeRuleRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("key").(string))
	return resourceSonarqubeRuleRead(d, m)
}
