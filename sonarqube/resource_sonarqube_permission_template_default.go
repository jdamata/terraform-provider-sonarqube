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

type DefaultTemplate struct {
	TemplateId string `json:"templateId"`
	Qualifier  string `json:"qualifier"`
}

type SearchPermissionTemplatesResponse struct {
	PermissionTemplates []PermTemplate    `json:"permissionTemplates"`
	DefaultTemplates    []DefaultTemplate `json:"defaultTemplates"`
}

type PermTemplate struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func resourceSonarqubePermissionTemplateDefault() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Sonarqube Permission Template Default resource. This can be used to set a permission template as the default for a component qualifier.",
		Create:      resourceSonarqubePermissionTemplateDefaultCreate,
		Read:        resourceSonarqubePermissionTemplateDefaultRead,
		Update:      resourceSonarqubePermissionTemplateDefaultUpdate,
		Delete:      resourceSonarqubePermissionTemplateDefaultDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubePermissionTemplateDefaultImport,
		},

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the permission template to set as default.",
			},
			"qualifier": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Component qualifier: TRK (project), VW (portfolio, Enterprise+), APP (application, Enterprise+).",
				ValidateFunc: validation.StringInSlice([]string{"TRK", "VW", "APP"}, false),
			},
		},
	}
}

func resourceSonarqubePermissionTemplateDefaultCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/set_default_template"

	sonarQubeURL.RawQuery = url.Values{
		"templateId": []string{d.Get("template_id").(string)},
		"qualifier":  []string{d.Get("qualifier").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubePermissionTemplateDefaultCreate",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubePermissionTemplateDefaultCreate: Failed to call %s: %+v", sonarQubeURL.Path, err)
	}
	defer resp.Body.Close()

	d.SetId(d.Get("qualifier").(string))

	return resourceSonarqubePermissionTemplateDefaultRead(d, m)
}

func resourceSonarqubePermissionTemplateDefaultRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/search_templates"

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubePermissionTemplateDefaultRead",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubePermissionTemplateDefaultRead: Failed to call %s: %+v", sonarQubeURL.Path, err)
	}
	defer resp.Body.Close()

	templatesResponse := SearchPermissionTemplatesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&templatesResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubePermissionTemplateDefaultRead: Failed to decode json into struct: %+v", err)
	}

	for _, defaultTemplate := range templatesResponse.DefaultTemplates {
		if defaultTemplate.Qualifier == d.Id() {
			if err := d.Set("template_id", defaultTemplate.TemplateId); err != nil {
				return fmt.Errorf("resourceSonarqubePermissionTemplateDefaultRead: Failed to set template_id: %+v", err)
			}
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceSonarqubePermissionTemplateDefaultUpdate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/set_default_template"

	sonarQubeURL.RawQuery = url.Values{
		"templateId": []string{d.Get("template_id").(string)},
		"qualifier":  []string{d.Get("qualifier").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubePermissionTemplateDefaultUpdate",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubePermissionTemplateDefaultUpdate: Failed to call %s: %+v", sonarQubeURL.Path, err)
	}
	defer resp.Body.Close()

	return resourceSonarqubePermissionTemplateDefaultRead(d, m)
}

func resourceSonarqubePermissionTemplateDefaultDelete(d *schema.ResourceData, m interface{}) error {
	// There is no API endpoint to unset a default template in SonarQube.
	// Remove the resource from Terraform state only.
	return nil
}

func resourceSonarqubePermissionTemplateDefaultImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// import id is the qualifier: TRK, VW, or APP
	qualifier := d.Id()
	if qualifier != "TRK" && qualifier != "VW" && qualifier != "APP" {
		return nil, fmt.Errorf("resourceSonarqubePermissionTemplateDefaultImport: Import id '%s' must be one of TRK, VW, APP", qualifier)
	}

	if err := d.Set("qualifier", qualifier); err != nil {
		return nil, fmt.Errorf("resourceSonarqubePermissionTemplateDefaultImport: Failed to set qualifier: %+v", err)
	}

	if err := resourceSonarqubePermissionTemplateDefaultRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
