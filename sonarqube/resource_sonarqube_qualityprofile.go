package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// QualityProfile struct
type QualityProfile struct {
	IsDefault    bool   `json:"isDefault"`
	IsInherited  bool   `json:"isInherited"`
	Language     string `json:"language"`
	LanguageName string `json:"languageName"`
	Name         string `json:"name"`
	Key          string `json:"key"`
}

// CreateQualityProfileResponse for unmarshalling response body from creating quality profiles
type CreateQualityProfileResponse struct {
	Profile  QualityProfile `json:"profile"`
	Warnings []string       `json:"warnings"`
}

// GetQualityProfileList for unmarshalling response body of quality gate get
type GetQualityProfileList struct {
	Profiles []GetQualityProfile `json:"profiles"`
}

// GetQualityProfile for unmarshalling response body of quality gate get
type GetQualityProfile struct {
	Key                       string                   `json:"key"`
	Name                      string                   `json:"name"`
	Language                  string                   `json:"language"`
	LanguageName              string                   `json:"languageName"`
	IsInherited               bool                     `json:"isInherited"`
	IsBuiltIn                 bool                     `json:"isBuiltIn"`
	ActiveRuleCount           int                      `json:"activeRuleCount"`
	ActiveDeprecatedRuleCount int                      `json:"activeDeprecatedRuleCount"`
	IsDefault                 bool                     `json:"isDefault"`
	RuleUpdatedAt             string                   `json:"ruleUpdatedAt"`
	LastUsed                  string                   `json:"lastUsed"`
	Actions                   GetQualityProfileActions `json:"actions"`
}

// GetQualityProfileActions for unmarshalling response body of quality gate get
type GetQualityProfileActions struct {
	Edit              bool `json:"edit"`
	SetAsDefault      bool `json:"setAsDefault"`
	Copy              bool `json:"copy"`
	Delete            bool `json:"delete"`
	AssociateProjects bool `json:"associateProjects"`
}

// Returns the resource represented by this file.
func resourceSonarqubeQualityProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeQualityProfileCreate,
		Read:   resourceSonarqubeQualityProfileRead,
		Delete: resourceSonarqubeQualityProfileDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeQualityProfileImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Quality profile name",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 100),
				),
			},
			"key": {
				Type:        schema.TypeString,
				Description: "Quality profile key",
				Computed:    true,
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
			"isDefault": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is the default profile",
			},
		},
	}
}

func resourceSonarqubeQualityProfileCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualityprofiles/create"
	sonarQubeURL.RawQuery = url.Values{
		"name":     []string{d.Get("name").(string)},
		"language": []string{d.Get("language").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeQualityProfileCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityProfileResponse := CreateQualityProfileResponse{}
	err = json.NewDecoder(resp.Body).Decode(&qualityProfileResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileCreate: Failed to decode json into struct: %+v", err)
	}

	if d.Get("isDefault").(bool) {
		sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
		sonarQubeURL.Path = "api/qualityprofiles/set_default"
		sonarQubeURL.RawQuery = url.Values{
			"qualityProfile": []string{d.Get("name").(string)},
			"language":       []string{d.Get("language").(string)},
		}.Encode()

		respDefault, errDefault := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"POST",
			sonarQubeURL.String(),
			http.StatusOK,
			"resourceSonarqubeQualityProfileCreate",
		)
		if errDefault != nil {
			return err
		}
		defer respDefault.Body.Close()
	}

	d.SetId(qualityProfileResponse.Profile.Key)
	return resourceSonarqubeQualityProfileRead(d, m)
}

func resourceSonarqubeQualityProfileRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualityprofiles/search"

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeQualityProfileRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	getQualityProfileResponse := GetQualityProfileList{}
	err = json.NewDecoder(resp.Body).Decode(&getQualityProfileResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileRead: Failed to decode json into struct: %+v", err)
	}

	for _, value := range getQualityProfileResponse.Profiles {
		if d.Id() == value.Key {
			d.SetId(value.Key)
			d.Set("name", value.Name)
			d.Set("language", value.Language)
			d.Set("key", value.Key)
			return nil
		}
	}

	return fmt.Errorf("resourceSonarqubeQualityProfileRead: Failed to find project: %+v", d.Id())
}

func resourceSonarqubeQualityProfileDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualityprofiles/delete"
	sonarQubeURL.RawQuery = url.Values{
		"qualityProfile": []string{d.Get("name").(string)},
		"language":       []string{d.Get("language").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeQualityProfileDelete",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityProfileDelete: Failed to delete quality profile: %+v", err)
	}
	defer resp.Body.Close()

	return nil

}

func resourceSonarqubeQualityProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeQualityProfileRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
