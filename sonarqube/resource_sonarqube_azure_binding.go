package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetAzureBinding for unmarshalling response body from getting project binding details
type GetAzureBinding struct {
	Key        string `json:"key"`
	Alm        string `json:"alm"`
	Repository string `json:"repository"`
	Slug       string `json:"slug"`
	URL        string `json:"url"`
	Monorepo   bool   `json:"monorepo"`
}

// Returns the resource represented by this file.
func resourceSonarqubeAzureBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeAzureBindingCreate,
		Read:   resourceSonarqubeAzureBindingRead,
		Delete: resourceSonarqubeAzureBindingDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeAzureBindingImport,
		},
		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"alm_setting": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"monorepo": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func checkAzureBindingSupport(conf *ProviderConfiguration) error {
	if strings.ToLower(conf.sonarQubeEdition) == "community" {
		return fmt.Errorf("Azure Devops Bindings are not supported in the Community edition of SonarQube. You are using: SonarQube %s version %s", conf.sonarQubeEdition, conf.sonarQubeVersion)
	}
	return nil
}

func resourceSonarqubeAzureBindingCreate(d *schema.ResourceData, m interface{}) error {
	if err := checkAzureBindingSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/set_azure_binding"

	sonarQubeURL.RawQuery = url.Values{
		"almSetting":  []string{d.Get("alm_setting").(string)},
		"monorepo":    []string{strconv.FormatBool(d.Get("monorepo").(bool))},
		"project":     []string{d.Get("project").(string)},
		"projectName": []string{d.Get("project_name").(string)},
		"repository":  []string{d.Get("repository").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeAzureBindingCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	id := fmt.Sprintf("%v/%v/%v", d.Get("project").(string), d.Get("project_name").(string), d.Get("repository").(string))
	d.SetId(id)

	return resourceSonarqubeAzureBindingRead(d, m)
}

func resourceSonarqubeAzureBindingRead(d *schema.ResourceData, m interface{}) error {
	if err := checkAzureBindingSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

	idSlice := strings.SplitN(d.Id(), "/", 2)
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/get_binding"
	sonarQubeURL.RawQuery = url.Values{
		"project": []string{idSlice[0]},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeAzureBindingRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	BindingReadResponse := GetAzureBinding{}
	err = json.NewDecoder(resp.Body).Decode(&BindingReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeAzureBindingRead: Failed to decode json into struct: %+v", err)
	}

	// if
	if idSlice[1] == BindingReadResponse.Repository && BindingReadResponse.Alm == "azure" {
		d.Set("project", BindingReadResponse.SummaryCommentEnabled)
		d.Set("repository", BindingReadResponse.Repository)
		d.Set("alm_setting", BindingReadResponse.Key)
		d.Set("monorepo", BindingReadResponse.Monorepo)

		return nil
	}
	return fmt.Errorf("resourceSonarqubeAzureBindingRead: Failed to find azure binding: %+v", d.Id())
}

func resourceSonarqubeAzureBindingDelete(d *schema.ResourceData, m interface{}) error {
	if err := checkAzureBindingSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/delete_binding"
	sonarQubeURL.RawQuery = url.Values{
		"project": []string{d.Get("project").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeAzureBindingDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeAzureBindingImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeAzureBindingRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
