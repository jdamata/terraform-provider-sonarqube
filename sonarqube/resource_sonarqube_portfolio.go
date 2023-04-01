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

// Portfolio used in Portfolio
type Portfolio struct {
	Key           string   `json:"key"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Qualifier     string   `json:"qualifier"`
	Visibility    string   `json:"visibility"`
	SelectionMode string   `json:"selectionMode"`
	Branch        string   `json:"branch"`
	Tags          []string `json:"tags"`
	Regexp        string   `json:"regexp"`
}

// Returns the resource represented by this file.
func resourceSonarqubePortfolio() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubePortfolioCreate,
		Read:   resourceSonarqubePortfolioRead,
		Update: resourceSonarqubePortfolioUpdate,
		Delete: resourceSonarqubePortfolioDelete,
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
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "public",
				ForceNew:     true, // TODO: There currently isn't an API to update this in-place, even though it's possible in the UI 
				ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
			},
			"selection_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "NONE",
				ForceNew:     false,
				ValidateFunc: validation.StringInSlice([]string{"NONE", "MANUAL", "TAGS", "REGEXP", "REST"}, false),
			},
			"branch": { // Only active for TAGS, REGEXP and REST
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Description: "Which branch to analyze. If nothing, or \"\" is specified, the main branch is used.",
			},
			"tags": { // Only active for TAGS
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      false,
				ConflictsWith: []string{"regexp"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"regexp": { // Only active for REGEXP
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      false,
				ConflictsWith: []string{"tags"},
				ValidateFunc:  validation.StringIsValidRegExp,
			},
			// TODO: MANUAL
			// "selectedProjects": [],
			// "projects": [],
		},
	}
}

func checkPortfolioSupport(conf *ProviderConfiguration) error {
	if strings.ToLower(conf.sonarQubeEdition) != "enterprise" {
		return fmt.Errorf("Portfolios are only supported in the Enterprise edition of SonarQube. You are using: SonarQube %s version %s", conf.sonarQubeEdition, conf.sonarQubeVersion)
	}
	return nil
}

// Validate the regexp and tag fields if the corresponding selection_mode is chosen
func validatePortfolioResource(d *schema.ResourceData) error {
	switch selectionMode := d.Get("selection_mode"); selectionMode {
	case "NONE", "MANUAL", "REST":
		return nil

	case "TAGS":
		tags := d.Get("tags").([]interface{})
		if len(tags) == 0 {
			return fmt.Errorf("validatePortfolioResource: When selection_mode is set to TAGS, you need atleast 1 tag, got: %+v", d.Get("tags"))
		}

		for _, tag := range d.Get("tags").([]interface{}) {
			tagString := fmt.Sprint(tag)
			if len(tagString) == 0 {
				return fmt.Errorf("validatePortfolioResource: When selection_mode is set to TAGS, each tag must be non 0, got: %s", tagString)
			}
		}
		return nil

	case "REGEXP":
		regexp := d.Get("regexp").(string)
		if len(regexp) == 0 {
			return fmt.Errorf("validatePortfolioResource: When selection_mode is set to REGEXP, regexp must be set, got: \"%s\"", regexp)
		}
		return nil

	default:
		return fmt.Errorf("resourceSonarqubePortfolioCreate: selection_mode needs to be set to one of NONE, MANUAL, TAGS, REGEXP, REST")
	}
}

func portfolioSetSelectionMode(d *schema.ResourceData, m interface{}, sonarQubeURL url.URL) error {
	var endpoint string
	switch selectionMode := d.Get("selection_mode"); selectionMode {
	case "NONE":
		endpoint = "/api/views/set_none_mode"
		sonarQubeURL.RawQuery = url.Values{
			"portfolio": []string{d.Get("key").(string)},
		}.Encode()

	case "MANUAL":
		endpoint = "/api/views/set_manual_mode"
		sonarQubeURL.RawQuery = url.Values{
			"portfolio": []string{d.Get("key").(string)},
		}.Encode()

	case "TAGS":
		endpoint = "/api/views/set_tags_mode"

		var tags []string
		for _, v := range d.Get("tags").([]interface{}) {
			tags = append(tags, fmt.Sprint(v))
		}
		tagsCSV := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(tags)), ","), "[]")

		urlParameters := url.Values{
			"portfolio": []string{d.Get("key").(string)},
			"tags":      []string{tagsCSV},
		}

		// SonarQube handles "" like it actually is a name of a branch, see PR for reference: TODO: Add Link to PR
		branch := d.Get("branch").(string)
		if len(branch) > 0 {
			urlParameters.Add("branch", branch)
		}

		sonarQubeURL.RawQuery = urlParameters.Encode()

	case "REGEXP":
		endpoint = "/api/views/set_regexp_mode"

		urlParameters := url.Values{
			"portfolio": []string{d.Get("key").(string)},
			"regexp":    []string{d.Get("regexp").(string)},
		}

		// SonarQube handles "" like it actually is a name of a branch, see PR for reference: TODO: Add Link to PR
		branch := d.Get("branch").(string)
		if len(branch) > 0 {
			urlParameters.Add("branch", branch)
		}

		sonarQubeURL.RawQuery = urlParameters.Encode()

	case "REST":
		endpoint = "/api/views/set_remaining_projects_mode"

		urlParameters := url.Values{
			"portfolio": []string{d.Get("key").(string)},
		}

		// SonarQube handles "" like it actually is a name of a branch, see PR for reference: TODO: Add Link to PR
		branch := d.Get("branch").(string)
		if len(branch) > 0 {
			urlParameters.Add("branch", branch)
		}

		sonarQubeURL.RawQuery = urlParameters.Encode()

	default:
		return fmt.Errorf("resourceSonarqubePortfolioCreate: selection_mode needs to be set to one of NONE, MANUAL, TAGS, REGEXP, REST")
	}

	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + endpoint

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubePortfolioCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubePortfolioCreate(d *schema.ResourceData, m interface{}) error {
	if err := checkPortfolioSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

	err := validatePortfolioResource(d)
	if err != nil {
		return err
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/views/create"

	sonarQubeURL.RawQuery = url.Values{
		"description": []string{d.Get("description").(string)},
		"key":         []string{d.Get("key").(string)},
		"name":        []string{d.Get("name").(string)},
		"visibility":  []string{d.Get("visibility").(string)},
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

	err = portfolioSetSelectionMode(d, m, m.(*ProviderConfiguration).sonarQubeURL)
	if err != nil {
		return err
	}

	// Decode response into struct
	portfolioResponse := Portfolio{}
	err = json.NewDecoder(resp.Body).Decode(&portfolioResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubePortfolioCreate: Failed to decode json into struct: %+v", err)
	}

	d.SetId(portfolioResponse.Key)
	return resourceSonarqubePortfolioRead(d, m)
}

func resourceSonarqubePortfolioRead(d *schema.ResourceData, m interface{}) error {
	if err := checkPortfolioSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

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
	portfolioReadResponse := Portfolio{}
	err = json.NewDecoder(resp.Body).Decode(&portfolioReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubePortfolioRead: Failed to decode json into struct: %+v", err)
	}

	d.SetId(portfolioReadResponse.Key)
	d.Set("key", portfolioReadResponse.Key)
	d.Set("name", portfolioReadResponse.Name)
	d.Set("description", portfolioReadResponse.Desc)
	d.Set("qualifier", portfolioReadResponse.Qualifier)
	d.Set("visibility", portfolioReadResponse.Visibility)
	d.Set("selection_mode", portfolioReadResponse.SelectionMode)

	// These fields may or may not be set in the reposnse from SonarQube depending on the selection_mode
	if len(portfolioReadResponse.Tags) > 0 {
		d.Set("tags", portfolioReadResponse.Tags)
	}
	if len(portfolioReadResponse.Branch) > 0 {
		d.Set("branch", portfolioReadResponse.Branch)
	}
	if len(portfolioReadResponse.Regexp) > 0 {
		d.Set("regexp", portfolioReadResponse.Regexp)
	}

	return nil
}

func resourceSonarqubePortfolioUpdate(d *schema.ResourceData, m interface{}) error {
	if err := checkPortfolioSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

	err := validatePortfolioResource(d)
	if err != nil {
		return err
	}


	if d.HasChanges("name", "description") {
		sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/views/update"
		sonarQubeURL.RawQuery = url.Values{
			"key":         []string{d.Id()},
			"description": []string{d.Get("description").(string)},
			"name":        []string{d.Get("name").(string)},
		}.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"POST",
			sonarQubeURL.String(),
			http.StatusOK,
			"resourceSonarqubePortfolioUpdate",
		)
		if err != nil {
			return fmt.Errorf("error updating Sonarqube Portfolio Name and Description: %+v", err)
		}
		defer resp.Body.Close()
	}

	if d.HasChanges("selection_mode", "branch", "tags", "regexp") {
		err := portfolioSetSelectionMode(d, m, m.(*ProviderConfiguration).sonarQubeURL)
		if err != nil {
			return fmt.Errorf("error updating Sonarqube selection mode: %+v", err)
		}
	}

	return resourceSonarqubePortfolioRead(d, m)
}

func resourceSonarqubePortfolioDelete(d *schema.ResourceData, m interface{}) error {
	if err := checkPortfolioSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

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
	if err := resourceSonarqubePortfolioRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
