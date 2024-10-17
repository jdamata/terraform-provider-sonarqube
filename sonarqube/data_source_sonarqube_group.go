package sonarqube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubeGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a Sonarqube Group resource",
		Read:        dataSourceSonarqubeGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the group.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The group description.",
			},
		},
	}
}

func dataSourceSonarqubeGroupRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))
	return resourceSonarqubeGroupRead(d, m)
}
