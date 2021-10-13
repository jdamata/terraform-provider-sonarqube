package sonarqube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type GetRule struct {
	Rule  []Rule `json:"rules"`
	Total int    `json:"total"`
	P     int    `json:"p"`
	PS    int    `json:"ps"`
}

type Rule struct {
	keyRule   string `json:"key"`
	Repo      string `json:"repo"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	// UpdatedAt   string   `json:"updatedAt"`
	HtmlDesc string `json:"htmlDesc"`
	Severity string `json:"severity"`
	Status   string `json:"status"`
	// InternalKey string   `json:"internalKey"`
	IsTemplate string `json:"isTemplate"`
	// Tags        []string `json:"tags"`
	SysTags    []string `json:"sysTags"`
	Lang       string   `json:"lang"`
	LangName   string   `json:"langName"`
	Scope      string   `json:"scope"`
	IsExternal string   `json:"isExternal"`
	Type       string   `json:"type"`
	Params     Params   `json:"params"`
}

type Params struct {
	KeyParm      string `json:"key"`
	Desc         string `json:"desc"`
	DefaultValue int    `json:"defaultValue"`
	Type         string `json:"type"`
}

type CreateRuleResponse struct {
	Rule Rule `json:"rule"`
}

func resourceSonarqubeRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeRuleCreate,
		Read:   resourceSonarqubeRuleRead,
		Update: resourceSonarqubeRuleUpdate,
		Delete: resourceSonarqubeRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeRuleImporter,
		},

		Schema: map[string]*schema.Schema{
			"custom_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key of the custom rule",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 200),
				),
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
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 200),
				),
			},
			"params": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parameters as semi-colon list of =, for example 'params=key1=v1;key2=v2' (Only for custom rule)",
			},
			"prevent_reactivation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If set to true and if the rule has been deactivated (status 'REMOVED'), a status 409 will be returned",
				Default:     "false",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						[]string{"true", "false", "yes", "no"},
						false,
					),
				),
			},
			"severity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule severity",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						[]string{"INFO", "MINOR", "MAJOR", "CRITICAL", "BLOCKER"},
						false,
					),
				),
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule status",
				Default:     "READY",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						[]string{"BETA", "DEPRECATED", "READY", "REMOVED"},
						false,
					),
				),
			},
			"template_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Key of the template rule in order to create a custom rule (mandatory for custom rule)",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule type",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						[]string{"CODE_SMELL", "BUG", "VULNERABILITY", "SECURITY_HOTSPOT"},
						false,
					),
				),
			},
		},
	}
}
