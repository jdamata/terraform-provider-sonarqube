package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetBinding for unmarshalling response body from getting project binding details
type GetBinding struct {
	Key                   string `json:"key"`
	Alm                   string `json:"alm"`
	Repository            string `json:"repository"`
	URL                   string `json:"url"`
	SummaryCommentEnabled bool   `json:"summaryCommentEnabled"`
	Monorepo              bool   `json:"monorepo"`
}

// Returns the resource represented by this file.
func resourceSonarqubeGithubBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeGithubBindingCreate,
		Read:   resourceSonarqubeGithubBindingRead,
		Delete: resourceSonarqubeGithubBindingDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"almsetting": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"monorepo": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "no",
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"summarycommentenabled": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "true",
				ForceNew: true,
			},
		},
	}
}

func resourceSonarqubeGithubBindingCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/set_github_binding"

	sonarQubeURL.RawQuery = url.Values{
		"almSetting":            []string{d.Get("almsetting").(string)},
		"monorepo":              []string{d.Get("monorepo").(string)},
		"project":               []string{d.Get("project").(string)},
		"repository":            []string{d.Get("repository").(string)},
		"summaryCommentEnabled": []string{d.Get("summarycommentenabled").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeGithubBindingCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	id := fmt.Sprintf("%v/%v", d.Get("project").(string), d.Get("repository").(string))
	d.SetId(id)

	return resourceSonarqubeGithubBindingRead(d, m)
}

func resourceSonarqubeGithubBindingRead(d *schema.ResourceData, m interface{}) error {
	idSlice := strings.Split(d.Id(), "/")
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
		"resourceSonarqubeGithubBindingRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	BindingReadResponse := GetBinding{}
	err = json.NewDecoder(resp.Body).Decode(&BindingReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeGithubBindingRead: Failed to decode json into struct: %+v", err)
	}
	// Loop over all branches to see if the main branch we need exists.
	if idSlice[1] == BindingReadResponse.Repository && BindingReadResponse.Alm == "github" {
		d.Set("project", idSlice[0])
		d.Set("repository", idSlice[1])
		d.Set("almsetting", BindingReadResponse.Key)
		return nil
	}
	return fmt.Errorf("resourceSonarqubeGithubBindingRead: Failed to find github binding: %+v", d.Id())

}

func resourceSonarqubeGithubBindingDelete(d *schema.ResourceData, m interface{}) error {
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
		"resourceSonarqubeGithubBindingDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
