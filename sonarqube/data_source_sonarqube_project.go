package sonarqube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubeProject() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a Sonarqube project resource",
		Read:        dataSourceSonarqubeProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the project",
			},
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The project key of the project",
			},
			"visibility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project visibility",
			},
		},
	}
}

func dataSourceSonarqubeProjectRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("project").(string))
	return resourceSonarqubeProjectRead(d, m)
}
