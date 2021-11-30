package sonarqube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubeProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSonarqubeProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSonarqubeProjectRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("project").(string))
	return resourceSonarqubeProjectRead(d, m)
}
