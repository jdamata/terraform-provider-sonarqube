package sonarqube

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Rule struct
type Rule struct {
	Key         string       `json:"key"`
	Repo        string       `json:"repo"`
	Name        string       `json:"name"`
	CreatedAt   string       `json:"createdAt"`
	HTMLDesc    string       `json:"htmlDesc"`
	MdDesc      string       `json:"mdDesc"`
	Severity    string       `json:"severity"`
	Status      string       `json:"status"`
	IsTemplate  bool         `json:"isTemplate"`
	TemplateKey string       `json:"templateKey"`
	SysTags     []string     `json:"sysTags"`
	Lang        string       `json:"lang"`
	LangName    string       `json:"langName"`
	Params      []RuleParams `json:"params"`
	Scope       string       `json:"scope"`
	IsExternal  bool         `json:"isExternal"`
	Type        string       `json:"type"`
}

// RuleParams struct
type RuleParams struct {
	Key          string `json:"key"`
	HTMLDesc     string `json:"htmlDesc"`
	DefaultValue string `json:"defaultValue"`
	Type         string `json:"type"`
}

// Returns the resource represented by this file.
func resourceSonarqubeCustomRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeCustomRuleCreate,
		Read:   resourceSonarqubeCustomRuleRead,
		Update: resourceSonarqubeCustomRuleUpdate,
		Delete: resourceSonarqubeCustomRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeCustomRuleImport,
		},

		// Define the fields of this schema.
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
				Description: "Rule name",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 200),
				),
			},
			"template_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key of the template rule in order to create a custom rule",
			},
			"severity": {
				Type:        schema.TypeString,
				Description: "Rule severity",
				Optional:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						[]string{"INFO", "MINOR", "MAJOR", "CRITICAL", "BLOCKER"},
						false,
					),
				),
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Rule status",
				Default:     "READY",
				Optional:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						[]string{"BETA", "DEPRECATED", "READY", "REMOVED"},
						false,
					),
				),
			},
			"rule_type": {
				Type:        schema.TypeString,
				Description: "Rule type",
				ForceNew:    true,
				Optional:    true,
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

func resourceSonarqubeCustomRuleCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/rules/create"

	rawQuery := url.Values{
		"custom_key":           []string{d.Get("key").(string)},
		"markdown_description": []string{d.Get("markdown_description").(string)},
		"name":                 []string{d.Get("name").(string)},
		"status":               []string{d.Get("status").(string)},
		"template_key":         []string{d.Get("template_key").(string)},
	}

	if severity, ok := d.GetOk("severity"); ok {
		rawQuery.Add("severity", severity.(string))
	}

	if ruleType, ok := d.GetOk("rule_type"); ok {
		rawQuery.Add("rule_type", ruleType.(string))
	}

	sonarQubeURL.RawQuery = rawQuery.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeCustomRuleCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	rule := Rule{}
	err = json.NewDecoder(resp.Body).Decode(&rule)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeCustomRuleCreate: Failed to decode json into struct: %+v", err)
	}

	log.Printf("[DEBUG][resourceSonarqubeCustomRuleCreate] Name is %s", rule.Name)
	d.SetId(rule.Key)
	return resourceSonarqubeCustomRuleRead(d, m)
}

func resourceSonarqubeCustomRuleRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/rules/show"
	sonarQubeURL.RawQuery = url.Values{
		"key": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeCustomRuleRead",
	)
	if err != nil {
		return fmt.Errorf("Error reading Sonarqube rule: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	rule := Rule{}
	err = json.NewDecoder(resp.Body).Decode(&rule)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeCustomRuleRead: Failed to decode json into struct: %+v", err)
	}

	log.Printf("[DEBUG][resourceSonarqubeCustomRuleRead] Checking '%s' with '%s'", d.Id(), rule.Key)
	log.Printf("[DEBUG][resourceSonarqubeCustomRuleRead] MdDesc: '%s'", rule.MdDesc)
	log.Printf("[DEBUG][resourceSonarqubeCustomRuleRead] Name: '%s'", rule.Name)
	log.Printf("[DEBUG][resourceSonarqubeCustomRuleRead] Status: '%s'", rule.Status)
	log.Printf("[DEBUG][resourceSonarqubeCustomRuleRead] Severity: '%s'", rule.Severity)
	log.Printf("[DEBUG][resourceSonarqubeCustomRuleRead] Type: '%s'", rule.Type)
	log.Printf("[DEBUG][resourceSonarqubeCustomRuleRead] TemplateKey: '%s'", rule.TemplateKey)

	d.SetId(rule.Key)
	d.Set("markdown_description", rule.MdDesc)
	d.Set("name", rule.Name)
	d.Set("status", rule.Status)
	d.Set("severity", rule.Severity)
	d.Set("rule_type", rule.Type)
	d.Set("template_key", rule.TemplateKey)

	return nil
}

func resourceSonarqubeCustomRuleUpdate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/rules/update"

	rawQuery := url.Values{
		"key":                  []string{d.Id()},
		"markdown_description": []string{d.Get("markdown_description").(string)},
		"name":                 []string{d.Get("name").(string)},
		"status":               []string{d.Get("status").(string)},
		"template_key":         []string{d.Get("template_key").(string)},
	}

	if severity, ok := d.GetOk("severity"); ok {
		rawQuery.Add("severity", severity.(string))
	}

	sonarQubeURL.RawQuery = rawQuery.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeCustomRuleUpdate",
	)
	if err != nil {
		return fmt.Errorf("Error updating Sonarqube rule: %+v", err)
	}
	defer resp.Body.Close()

	return resourceSonarqubeCustomRuleRead(d, m)
}

func resourceSonarqubeCustomRuleDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/rules/delete"
	sonarQubeURL.RawQuery = url.Values{
		"key": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeCustomRuleDelete",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeCustomRuleDelete: error deleting Sonarqube rule: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeCustomRuleImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeCustomRuleRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
