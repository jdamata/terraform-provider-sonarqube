package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type Rule struct {
	RuleKey     string   `json:"key"`
	Repo        string   `json:"repo"`
	Name        string   `json:"name"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
	HtmlDesc    string   `json:"htmlDesc,omitempty"`
	MdDesc      string   `json:"mdDesc,omitempty"`
	Severity    string   `json:"severity"`
	Status      string   `json:"status"`
	InternalKey string   `json:"internalKey"`
	IsTemplate  bool     `json:"isTemplate"`
	Tags        []string `json:"tags"`
	TemplateKey string   `json:"templateKey,omitempty"`
	SysTags     []string `json:"sysTags"`
	Lang        string   `json:"lang"`
	LangName    string   `json:"langName"`
	Scope       string   `json:"scope"`
	IsExternal  bool     `json:"isExternal"`
	Type        string   `json:"type"`
	Params      []Params `json:"params,omitempty"`
}

type Params struct {
	ParmKey      string `json:"key"`
	HtmlDesc     string `json:"htmlDesc"`
	DefaultValue string `json:"defaultValue"`
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
		Description: "Provides a Sonarqube Rules resource. This can be used to manage Sonarqube rules.",
		Create:      resourceSonarqubeRuleCreate,
		Read:        resourceSonarqubeRuleRead,
		Update:      resourceSonarqubeRuleUpdate,
		Delete:      resourceSonarqubeRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeRuleImporter,
		},

		Schema: map[string]*schema.Schema{
			"custom_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "key of the custom rule should only contain : a-z, 0-9, \\_",
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
				Type:     schema.TypeString,
				Optional: true,
				Description: `Parameters as semi-colon list of =, for example 'params=key1=v1;key2=v2' (Only for custom rule)
  - parameter order: expression=value;filePattern=value;message=value`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"prevent_reactivation": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `If set to true and if the rule has been deactivated (status 'REMOVED'), a status 409 will be returned
  - Possible values - true, false, yes, no`,
				Default: "false",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						[]string{"true", "false", "yes", "no"},
						false,
					),
				),
			},
			"severity": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Rule severity
  - Possible values - INFO, MINOR, MAJOR, CRITICAL, BLOCKER`,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						[]string{"INFO", "MINOR", "MAJOR", "CRITICAL", "BLOCKER"},
						false,
					),
				),
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Rule status
  - Possible values - BETA, DEPRECATED, READY, REMOVED
  - Default value - READY`,
				Default: "READY",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						[]string{"BETA", "DEPRECATED", "READY", "REMOVED"},
						false,
					),
				),
			},
			"template_key": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Key of the template rule in order to create a custom rule (mandatory for custom rule)
  - [Example values](https://docs.sonarqube.org/latest/user-guide/rules/#header-4)`,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Rule type
  - Possible values - CODE_SMELL, BUG, VULNERABILITY, SECURITY_HOTSPOT`,
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
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/rules/create"

	sonarQubeURL.RawQuery = url.Values{
		"customKey":           []string{d.Get("custom_key").(string)},
		"markdownDescription": []string{d.Get("markdown_description").(string)},
		"name":                []string{d.Get("name").(string)},
		"params":              []string{d.Get("params").(string)},
		"preventReactivation": []string{d.Get("prevent_reactivation").(string)},
		"severity":            []string{d.Get("severity").(string)},
		"status":              []string{d.Get("status").(string)},
		"templateKey":         []string{d.Get("template_key").(string)},
		"type":                []string{d.Get("type").(string)},
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
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/rules/search"
	sonarQubeURL.RawQuery = url.Values{
		"rule_key": []string{d.Id()},
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
			d.Set("markdown_description", value.MdDesc)
			d.Set("name", value.Name)
			d.Set("severity", value.Severity)
			d.Set("template_key", value.TemplateKey)
			d.Set("status", value.Status)
			d.Set("type", value.Type)
			return nil
		}
	}
	return fmt.Errorf("resourceSonarqubeRuleRead: Failed to find project: %+v", d.Id())
}

func resourceSonarqubeRuleDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/rules/delete"
	sonarQubeURL.RawQuery = url.Values{
		"key": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
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
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/rules/update"
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

	return resourceSonarqubeRuleRead(d, m)
}
