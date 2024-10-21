package sonarqube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubeQualityGate() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a Sonarqube qualitygate resource",
		Read:        dataSourceSonarqubeQualityGateRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Quality Gate.",
			},
			"copy_from": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Origin of the Quality Gate",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Quality Gate default.",
			},
			"condition": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric": {
							Type:     schema.TypeString,
							Required: true,
						},
						"op": {
							Type:     schema.TypeString,
							Required: true,
						},
						"threshold": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Description: "List of Quality Gate conditions.",
			},
		},
	}
}

func dataSourceSonarqubeQualityGateRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("name").(string))
	return resourceSonarqubeQualityGateRead(d, m)
}
