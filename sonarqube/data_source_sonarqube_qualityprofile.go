package sonarqube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubeQualityProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a Sonarqube qualityprofile resource",
		Read:        dataSourceSonarqubeQualityProfileRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Quality Profile",
			},
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key of the Quality Profile",
			},
			"language": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Quality Profile language",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Quality Profile default",
			},
		},
	}
}

func dataSourceSonarqubeQualityProfileRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))
	return resourceSonarqubeQualityProfileRead(d, m)
}
