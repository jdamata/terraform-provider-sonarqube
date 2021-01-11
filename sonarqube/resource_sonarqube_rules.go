package sonarqube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Returns the resource represented by this file.
func resourceSonarqubeRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeRulesCreate,
		Read:   resourceSonarqubeRulesRead,
		Delete: resourceSonarqubeRulesDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeRulesImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"custom_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key of the custom rule",
			},
			"markdown_description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule description",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule name",
			},
			"prevent_reactivation": {
				Type:        schema.TypeBool,
				Description: "If set to true and if the rule has been deactivated (status 'REMOVED'), a status 409 will be returned",
				Default:     false,
			},
			"severity": {
				Type:        schema.TypeBool,
				Description: "Rule severity",
				Default:     false,
			},
		},
	}
}
