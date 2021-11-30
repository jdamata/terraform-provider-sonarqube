package sonarqube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubeRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSonarqubeRuleRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"markdown_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSonarqubeRuleRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("key").(string))
	return resourceSonarqubeRuleRead(d, m)
}
