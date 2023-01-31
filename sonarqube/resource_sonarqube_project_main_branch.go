package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetBranches for unmarshalling response body from getting project branch details
type GetBranches struct {
	Branches []Branches `json:"branches"`
}
type Status struct {
	QualityGateStatus string `json:"qualityGateStatus"`
}
type Branches struct {
	Name              string `json:"name"`
	IsMain            bool   `json:"isMain"`
	Type              string `json:"type"`
	Status            Status `json:"status"`
	AnalysisDate      string `json:"analysisDate"`
	ExcludedFromPurge bool   `json:"excludedFromPurge"`
}

// Returns the resource represented by this file.
func resourceSonarqubeProjectMainBranch() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeProjectMainBranchCreate,
		Read:   resourceSonarqubeProjectMainBranchRead,
		Delete: resourceSonarqubeProjectMainBranchDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // Investigate if we should do this or just update :)
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSonarqubeProjectMainBranchCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/project_branches/rename"

	sonarQubeURL.RawQuery = url.Values{
		"name":    []string{d.Get("name").(string)},
		"project": []string{d.Get("project").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent, // We know its 204 right ?
		"resourceSonarqubeProjectMainBranchCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	id := fmt.Sprintf("%v/%v", d.Get("project").(string), d.Get("name").(string))
	d.SetId(id)

	return resourceSonarqubeProjectMainBranchRead(d, m)
}

func resourceSonarqubeProjectMainBranchRead(d *schema.ResourceData, m interface{}) error {
	idSlice := strings.Split(d.Id(), "/") // This is kinda borrowed from project assocation as it also return a 204
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/project_branches/list"
	sonarQubeURL.RawQuery = url.Values{
		"project": []string{idSlice[0]},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeProjectMainBranchRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	projectReadResponse := GetBranches{}
	err = json.NewDecoder(resp.Body).Decode(&projectReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeProjectMainBranchRead: Failed to decode json into struct: %+v", err)
	}
	// Loop over all branches to see if the main branch we need exists.
	for _, value := range projectReadResponse.Branches {
		if idSlice[1] == value.Name && value.IsMain {
			d.Set("project", idSlice[0])
			d.Set("name", value.Name)
			return nil
		}
	}
	return fmt.Errorf("resourceSonarqubeProjectMainBranchRead: Failed to find project main branch: %+v", d.Id())

}

// I believe that the default main branch name is main. so delete could just be setting it back to main. alternatively we can read the RULE that sets up default main branch name and that way get the proper one
func resourceSonarqubeProjectMainBranchDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/project_branches/rename"
	sonarQubeURL.RawQuery = url.Values{
		"name":    []string{"main"}, // See comment we might wanna make it that
		"project": []string{d.Get("project").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeProjectMainBranchDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
