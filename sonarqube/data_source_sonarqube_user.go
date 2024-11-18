package sonarqube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubeUser() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a Sonarqube User resource",
		Read:        dataSourceSonarqubeUserRead,
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The login name of the user",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the user",
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email of the user",
			},
			"is_local": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the user is local",
			},
		},
	}
}

func dataSourceSonarqubeUserRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("login_name").(string))
	return resourceSonarqubeUserRead(d, m)
}
