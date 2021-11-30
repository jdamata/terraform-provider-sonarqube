package sonarqube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubeUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSonarqubeUserRead,
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_local": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSonarqubeUserRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("login_name").(string))
	return resourceSonarqubeUserRead(d, m)
}
