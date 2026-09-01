package sonarqube

import (
	"encoding/json"
	"errors"
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
		Description: "Provides a Sonarqube Project main branch resource. This can be used to create and manage a Sonarqube Projects main branch.",
		Create:      resourceSonarqubeProjectMainBranchCreate,
		Read:        resourceSonarqubeProjectMainBranchRead,
		Delete:      resourceSonarqubeProjectMainBranchDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeProjectMainBranchImport,
		},
		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name you want the main branch to have.",
			},
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key of the project. Maximum length 400. All letters, digits, dash, underscore, period or colon.",
			},
		},
	}
}

func resourceSonarqubeProjectMainBranchCreate(d *schema.ResourceData, m interface{}) error {
	project := d.Get("project").(string)
	name := d.Get("name").(string)

	// Check if branch already exists
	branches, err := getProjectBranches(project, m)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeProjectMainBranchCreate: Failed to get project branches: %w", err)
	}

	branchExists := false
	for _, branch := range branches {
		if branch.Name == name {
			branchExists = true
			break
		}
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL

	if branchExists {
		// If branch exists, set it as main
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/project_branches/set_main"
		sonarQubeURL.RawQuery = url.Values{
			"branch":  []string{name},
			"project": []string{project},
		}.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"POST",
			sonarQubeURL.String(),
			http.StatusNoContent,
			"resourceSonarqubeProjectMainBranchCreate/SetMain",
		)
		if err != nil {
			return err
		}
		resp.Body.Close()
	} else {
		// If branch does not exist, rename the current main branch
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/project_branches/rename"
		sonarQubeURL.RawQuery = url.Values{
			"name":    []string{name},
			"project": []string{project},
		}.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"POST",
			sonarQubeURL.String(),
			http.StatusNoContent,
			"resourceSonarqubeProjectMainBranchCreate/Rename",
		)
		if err != nil {
			return err
		}
		resp.Body.Close()
	}

	id := fmt.Sprintf("%v/%v", project, name)
	d.SetId(id)

	return resourceSonarqubeProjectMainBranchRead(d, m)
}

func resourceSonarqubeProjectMainBranchRead(d *schema.ResourceData, m interface{}) error {
	idSlice := strings.SplitN(d.Id(), "/", 2)
	project := idSlice[0]
	name := idSlice[1]

	branches, err := getProjectBranches(project, m)
	if err != nil {
		return err
	}

	// Loop over all branches to see if the main branch we need exists.
	for _, value := range branches {
		if name == value.Name && value.IsMain {
			errProject := d.Set("project", project)
			errName := d.Set("name", value.Name)
			return errors.Join(errProject, errName)
		}
	}
	return fmt.Errorf("resourceSonarqubeProjectMainBranchRead: Failed to find project main branch: %+v", d.Id())
}

func getProjectBranches(project string, m interface{}) ([]Branches, error) {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/project_branches/list"
	sonarQubeURL.RawQuery = url.Values{
		"project": []string{project},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"getProjectBranches",
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response into struct
	branchReadResponse := GetBranches{}
	err = json.NewDecoder(resp.Body).Decode(&branchReadResponse)
	if err != nil {
		return nil, fmt.Errorf("getProjectBranches: Failed to decode json into struct: %+v", err)
	}

	return branchReadResponse.Branches, nil
}

// TODO make the delete function read the default branch name of the sonarQube instance instead of assuming
func resourceSonarqubeProjectMainBranchDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/project_branches/rename"
	sonarQubeURL.RawQuery = url.Values{
		"name":    []string{"main"},
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

func resourceSonarqubeProjectMainBranchImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeProjectMainBranchRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
