package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type Rule struct {
	RuleKey     string   `json:"key"`
	Repo        string   `json:"repo"`
	Name        string   `json:"name"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
	HtmlDesc    string   `json:"htmlDesc"`
	MdDesc      string   `json:"mdDesc"`
	Severity    string   `json:"severity"`
	Status      string   `json:"status"`
	InternalKey string   `json:"internalKey"`
	IsTemplate  string   `json:"isTemplate"`
	Tags        []string `json:"tags"`
	SysTags     []string `json:"sysTags"`
	Lang        string   `json:"lang"`
	LangName    string   `json:"langName"`
	Scope       string   `json:"scope"`
	IsExternal  string   `json:"isExternal"`
	Type        string   `json:"type"`
	Params      Params   `json:"params"`
}

type Params struct {
	ParmKey      string `json:"key"`
	Desc         string `json:"desc"`
	DefaultValue int    `json:"defaultValue"`
	Type         string `json:"type"`
}

type GetRule struct {
	Rule  []Rule `json:"rules"`
	Total int    `json:"total"`
	P     int    `json:"p"`
	PS    int    `json:"ps"`
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
			"key": {
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
				ForceNew:    true,
				Description: "Rule name",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 200),
				),
			},
			"params": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parameters as semi-colon list of =, for example 'params=key1=v1;key2=v2' (Only for custom rule)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Required:    true,
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

func resourceSonarqubeRuleCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/rules/create"
	sonarQubeURL.RawQuery = url.Values{
		"key":                  []string{d.Get("key").(string)},
		"markdown_description": []string{d.Get("markdown_description").(string)},
		"name":                 []string{d.Get("name").(string)},
		"params":               []string{d.Get("params").(string)},
		"prevent_reactivation": []string{d.Get("prevent_reactivation").(string)},
		"severity":             []string{d.Get("severity").(string)},
		"status":               []string{d.Get("status").(string)},
		"template_key":         []string{d.Get("template_key").(string)},
		"type":                 []string{d.Get("type").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeRuleCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ruleCreateResponse := CreateRuleResponse{}
	err = json.NewDecoder(resp.Body).Decode(&ruleCreateResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeRuleCreate: Failed to decode json into struct: %+v", err)
	}

	d.SetId(ruleCreateResponse.Rule.RuleKey)
	return resourceSonarqubeRuleRead(d, m)
}

func resourceSonarqubeRuleRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/rules/search"
	sonarQubeURL.RawQuery = url.Values{
		"key": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeRuleRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ruleReadResponse := GetRule{}
	err = json.NewDecoder(resp.Body).Decode(&ruleReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeRuleRead: Failed to decode json into struct: %+v", err)
	}

	for _, value := range ruleReadResponse.Rule {
		if d.Id() == value.RuleKey {
			d.SetId(value.RuleKey)
			d.Set("name", value.Name)
			d.Set("markdown_description", value.MdDesc)
			return nil
		}
	}
	return fmt.Errorf("resourceSonarqubeRuleRead: Failed to find project: %+v", d.Id())
}

func resourceSonarqubeRuleDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/rules/delete"
	sonarQubeURL.RawQuery = url.Values{
		"key": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeRuleDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeRuleImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeRuleRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func resourceSonarqubeRuleUpdate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/rules/update"
	sonarQubeURL.RawQuery = url.Values{
		"key":                  []string{d.Id()},
		"markdown_description": []string{d.Get("markdown_description").(string)},
		"name":                 []string{d.Get("name").(string)},
		"params":               []string{d.Get("params").(string)},
		"severity":             []string{d.Get("severity").(string)},
		"status":               []string{d.Get("status").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeRuleUpdate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return resourceSonarqubeQualityGateConditionRead(d, m)
}