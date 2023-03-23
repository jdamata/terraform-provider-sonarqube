package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Portfolio used in CreatePortfolioResponse
type Portfolio struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Qualifier string `json:"qualifier"`
	Visibility string `json:"visibility"`
	SelectionMode string `json:"selectionMode"`
}

// CreatePortfolioResponse for unmarshalling response body of project creation
type CreatePortfolioResponse struct {
	Portfolio Portfolio `json:"portfolio"`
}

// GetPortfolio for unmarshalling response body of Portfolio get
type GetPortfolio struct {
	Portfolio Portfolio `json:"portfolio"`
}


// Returns the resource represented by this file.
func resourceSonarqubePortfolio() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubePortfolioCreate,
		Read:   resourceSonarqubePortfolioRead,
		Update: resourceSonarqubePortfolioUpdate,
		Delete: resourceSonarqubePortfoliotDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubePortfolioImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"qualifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "public", // TODO: Not sure if this should be public. The docs makes it sound like this is a global setting somewhere
				ForceNew: false, // TODO: Implement update
			},
			"selection_mode": {
				Type:     schema.TypeString,
				Optional: true, // TODO: Add extra call ic Create for when something other than "NONE" is specified
				Default:  "NONE", 
				ForceNew: false, // TODO: Implement update
				// TODO: Set "regexp" if mode="REGEXP", otherwise create/update fails
			},
		},
	}
}

func portfolioSetSelectionMode(d *schema.ResourceData, m interface{}, sonarQubeURL url.URL) error {
	var endpoint string
	switch selectionMode := d.Get("selection_mode"); selectionMode {
	case "NONE":
		endpoint = "/api/views/set_none_mode" 
		sonarQubeURL.RawQuery = url.Values{
			"portfolio":       []string{d.Get("key").(string)},
		}.Encode()

	case "MANUAL":
		endpoint = "/api/views/set_manual_mode"
		sonarQubeURL.RawQuery = url.Values{
			"portfolio":       []string{d.Get("key").(string)},
		}.Encode()
		
	case "TAGS":
		endpoint = "/api/views/set_tags_mode"
		sonarQubeURL.RawQuery = url.Values{
			"portfolio":       []string{d.Get("key").(string)},
			"tags":       []string{d.Get("tags").(string)}, // TODO: Support this, and validate csv and tags
		}.Encode()

	case "REGEXP":
		endpoint = "/api/views/set_regexp_mode"
		sonarQubeURL.RawQuery = url.Values{
			"portfolio":       []string{d.Get("key").(string)},
			"regexp":       []string{d.Get("regexp").(string)}, // TODO: Support this 
		}.Encode()

	case "REST":
		endpoint = "/api/views/set_remaining_projects_mode"
		sonarQubeURL.RawQuery = url.Values{
			"portfolio":       []string{d.Get("key").(string)},
		}.Encode()

	default:
		return fmt.Errorf("resourceSonarqubePortfolioCreate: selection_mode needs to be set to one of NONE, MANUAL, TAGS, REGEXP, REST")
	}

	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + endpoint

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubePortfolioCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubePortfolioCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/views/create"

	sonarQubeURL.RawQuery = url.Values{
		"description":       []string{d.Get("description").(string)},
		"key":    []string{d.Get("key").(string)},
		"name": []string{d.Get("name").(string)},
		"visibility": []string{d.Get("visibility").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubePortfolioCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = portfolioSetSelectionMode(d, m, sonarQubeURL)
	if err != nil {
		return err
	}

	// Decode response into struct
	portfolioResponse := CreatePortfolioResponse{}
	err = json.NewDecoder(resp.Body).Decode(&portfolioResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubePortfolioCreate: Failed to decode json into struct: %+v", err)
	}

	d.SetId(portfolioResponse.Portfolio.Key)
	return resourceSonarqubePortfolioRead(d, m)
}

func resourceSonarqubePortfolioRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/views/show"
	sonarQubeURL.RawQuery = url.Values{
		"key": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubePortfolioRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	portfolioReadResponse := GetPortfolio{}
	err = json.NewDecoder(resp.Body).Decode(&portfolioReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubePortfolioRead: Failed to decode json into struct: %+v", err)
	}

	d.SetId(portfolioReadResponse.Portfolio.Key)
	d.Set("name", portfolioReadResponse.Portfolio.Name)
	d.Set("qualifier", portfolioReadResponse.Portfolio.Qualifier)
	d.Set("visibility", portfolioReadResponse.Portfolio.Visibility)
	d.Set("selectionMode", portfolioReadResponse.Portfolio.SelectionMode)

	return nil

}

func resourceSonarqubePortfolioUpdate(d *schema.ResourceData, m interface{}) error {

	if d.HasChange("name") || d.HasChange("description") {
		sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/views/update"
		sonarQubeURL.RawQuery = url.Values{
			"key":    []string{d.Id()},
			"description": []string{d.Get("description").(string)},
			"name": []string{d.Get("name").(string)},
		}.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"POST",
			sonarQubeURL.String(),
			http.StatusNoContent,
			"resourceSonarqubePortfolioUpdate",
		)
		if err != nil {
			return fmt.Errorf("error updating Sonarqube Portfolio Name and Description: %+v", err)
		}
		defer resp.Body.Close()

	}


	return resourceSonarqubeProjectRead(d, m)
}

func resourceSonarqubePortfolioDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/views/delete"
	sonarQubeURL.RawQuery = url.Values{
		"key": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubePortfolioDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubePortfolioImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeProjectRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
