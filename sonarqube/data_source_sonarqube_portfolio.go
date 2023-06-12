package sonarqube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubePortfolio() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSonarqubePortfolioRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"qualifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"selection_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"branch": { // Only active for TAGS, REGEXP and REST
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": { // Only active for TAGS
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"regexp": { // Only active for REGEXP
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSonarqubePortfolioRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("key").(string))
	return resourceSonarqubePortfolioRead(d, m)
}
