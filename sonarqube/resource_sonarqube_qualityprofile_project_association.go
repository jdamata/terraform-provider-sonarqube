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

// GetQualityProfileProjectAssociation for unmarshalling response body from getting quality profile association
type GetQualityProfileProjectAssociation struct {
	Paging  Paging                                       `json:"paging"`
	Results []GetQualityProfileProjectAssociationResults `json:"results"`
}

// GetQualityProfileProjectAssociationResults used in GetQualityProfileProjectAssociation
type GetQualityProfileProjectAssociationResults struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Key      string `json:"key"`
	Selected bool   `json:"selected"`
}

// Returns the resource represented by this file.
func resourceSonarqubeQualityProfileProjectAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeQualityProfileProjectAssociationCreate,
		Read:   resourceSonarqubeQualityProfileProjectAssociationRead,
		Delete: resourceSonarqubeQualityProfileProjectAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeQualityProfileProjectAssociationImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"quality_profile": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Quality profile name",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 100),
				),
			},
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project name",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 100),
				),
			},
			"language": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Quality profile language",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						[]string{"cs", "css", "flex", "go", "java", "js", "jsp", "kotlin", "php", "py", "ruby", "scala", "ts", "vbnet", "web", "xml"},
						false,
					),
				),
			},
		},
	}
}

func resourceSonarqubeQualityProfileProjectAssociationCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualityprofiles/add_project"
	sonarQubeURL.RawQuery = url.Values{
		"language":       []string{d.Get("language").(string)},
		"project":        []string{d.Get("project").(string)},
		"qualityProfile": []string{d.Get("quality_profile").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeQualityProfileProjectAssociationCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	id := fmt.Sprintf("%v/%v", d.Get("quality_profile").(string), d.Get("project").(string))
	d.SetId(id)
	return resourceSonarqubeQualityProfileProjectAssociationRead(d, m)
}

func resourceSonarqubeQualityProfileProjectAssociationRead(d *schema.ResourceData, m interface{}) error {
	var language string
	var qualityProfile string

	// Id is composed of qualityProfile name and project name
	idSlice := strings.Split(d.Id(), "/")

	// Call api/qualityprofiles/search to return the qualityProfileID
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualityprofiles/search"

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeQualityProfileProjectAssociationRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	getQualityProfileResponse := GetQualityProfileList{}
	err = json.NewDecoder(resp.Body).Decode(&getQualityProfileResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileProjectAssociationRead: Failed to decode json into struct: %+v", err)
	}

	var qualityProfileID string
	for _, value := range getQualityProfileResponse.Profiles {
		if idSlice[0] == value.Name {
			qualityProfileID = value.Key
			language = value.Language
			qualityProfile = value.Name
		}
	}

	// With the qualityProfileID we can check if the project name is associated
	sonarQubeURL.Path = "api/qualityprofiles/projects"
	sonarQubeURL.RawQuery = url.Values{
		"key": []string{qualityProfileID},
	}.Encode()

	resp, err = httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeQualityProfileProjectAssociationRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	getQualityProfileProjectResponse := GetQualityProfileProjectAssociation{}
	err = json.NewDecoder(resp.Body).Decode(&getQualityProfileProjectResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileProjectAssociationRead: Failed to decode json into struct: %+v", err)
	}

	for _, value := range getQualityProfileProjectResponse.Results {
		if idSlice[1] == value.Name {
			d.SetId(d.Id())
			d.Set("project", value.Name)
			d.Set("quality_profile", qualityProfile)
			d.Set("language", language)
			return nil
		}
	}

	return fmt.Errorf("resourceSonarqubeQualityProfileProjectAssociationRead: Failed to find project association: %+v", d.Id())

}

func resourceSonarqubeQualityProfileProjectAssociationDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualityprofiles/remove_project"
	sonarQubeURL.RawQuery = url.Values{
		"language":       []string{d.Get("language").(string)},
		"project":        []string{d.Get("project").(string)},
		"qualityProfile": []string{d.Get("quality_profile").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeQualityProfileProjectAssociationDelete",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileProjectAssociationDelete: Failed to delete quality profile: %+v", err)
	}
	defer resp.Body.Close()

	return nil

}

func resourceSonarqubeQualityProfileProjectAssociationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeQualityProfileProjectAssociationRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
