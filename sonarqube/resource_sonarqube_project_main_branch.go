package sonarqube

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-version"
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
		Update:      resourceSonarqubeProjectMainBranchUpdate,
		Delete:      resourceSonarqubeProjectMainBranchDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeProjectMainBranchImport,
		},
		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
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
		http.StatusNoContent,
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
	idSlice := strings.SplitN(d.Id(), "/", 2)
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
	branchReadResponse := GetBranches{}
	err = json.NewDecoder(resp.Body).Decode(&branchReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeProjectMainBranchRead: Failed to decode json into struct: %+v", err)
	}
	// Loop over all branches to see if the main branch we need exists.
	for _, value := range branchReadResponse.Branches {
		if idSlice[1] == value.Name && value.IsMain {
			errProject := d.Set("project", idSlice[0])
			errName := d.Set("name", value.Name)
			return errors.Join(errProject, errName)
		}
	}
	return fmt.Errorf("resourceSonarqubeProjectMainBranchRead: Failed to find project main branch: %+v", d.Id())
}

// getBranches retrieves all branches for a given project
func getBranches(conf *ProviderConfiguration, projectKey string) (*GetBranches, error) {
	sonarQubeURL := conf.sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/project_branches/list"
	sonarQubeURL.RawQuery = url.Values{
		"project": []string{projectKey},
	}.Encode()

	resp, err := httpRequestHelper(
		conf.httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"getBranches",
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	branchResponse := &GetBranches{}
	err = json.NewDecoder(resp.Body).Decode(branchResponse)
	if err != nil {
		return nil, fmt.Errorf("getBranches: Failed to decode json into struct: %+v", err)
	}

	return branchResponse, nil
}

// deleteBranch removes a non-main branch from a project
func deleteBranch(conf *ProviderConfiguration, projectKey, branchName string) error {
	sonarQubeURL := conf.sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/project_branches/delete"
	sonarQubeURL.RawQuery = url.Values{
		"branch":  []string{branchName},
		"project": []string{projectKey},
	}.Encode()

	resp, err := httpRequestHelper(
		conf.httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"deleteBranch",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeProjectMainBranchUpdate(d *schema.ResourceData, m interface{}) error {
	conf := m.(*ProviderConfiguration)
	sonarQubeURL := conf.sonarQubeURL

	projectKey := d.Get("project").(string)
	targetBranchName := d.Get("name").(string)

	// set_main endpoint was introduced in SonarQube 10.2
	minimumVersionForSetMain, _ := version.NewVersion("10.2")
	useSetMainEndpoint := conf.sonarQubeVersion.GreaterThanOrEqual(minimumVersionForSetMain)

	var endpoint, paramName string
	if useSetMainEndpoint {
		endpoint = "/api/project_branches/set_main"
		paramName = "branch"
	} else {
		// For older versions using rename, we need to check if the target branch already exists
		// If it does and it's not the main branch, delete it first to avoid conflicts
		branches, err := getBranches(conf, projectKey)
		if err != nil {
			return fmt.Errorf("resourceSonarqubeProjectMainBranchUpdate: Failed to get branches: %+v", err)
		}

		for _, branch := range branches.Branches {
			if branch.Name == targetBranchName && !branch.IsMain {
				// Delete the existing non-main branch before renaming
				if err := deleteBranch(conf, projectKey, targetBranchName); err != nil {
					return fmt.Errorf("resourceSonarqubeProjectMainBranchUpdate: Failed to delete existing branch: %+v", err)
				}
				break
			}
		}

		endpoint = "/api/project_branches/rename"
		paramName = "name"
	}

	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + endpoint

	sonarQubeURL.RawQuery = url.Values{
		paramName: []string{targetBranchName},
		"project": []string{projectKey},
	}.Encode()

	resp, err := httpRequestHelper(
		conf.httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeProjectMainBranchUpdate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	id := fmt.Sprintf("%v/%v", projectKey, targetBranchName)
	d.SetId(id)

	return resourceSonarqubeProjectMainBranchRead(d, m)
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
