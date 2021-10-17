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
func resourceQualityProfile() *schema.Resource {
	return &schema.Resource{
		Create: CreateQualityProfile,
		Read:   ReadQualityProfile,
		Delete: DeleteQualityProfile,
		Importer: &schema.ResourceImporter{
			State: ImportQualityProfile,
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
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is the default profile",
				Default:     false,
				ForceNew:    true,
			},
		},
	}
}

func CreateQualityProfile(d *schema.ResourceData, m interface{}) error {
	qualityProfile := m.(*ProviderConfiguration).sonarQubeURL
	qualityProfile.Path = "api/qualityprofiles/create"
	qualityProfile.RawQuery = url.Values{
		"name":     []string{d.Get("name").(string)},
		"language": []string{d.Get("language").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		qualityProfile.String(),
		http.StatusOK,
		"CreateQualityProfile",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityProfileResponse := CreateQualityProfileResponse{}
	err = json.NewDecoder(resp.Body).Decode(&qualityProfileResponse)
	if err != nil {
		return fmt.Errorf("CreateQualityProfile: Failed to decode json into struct: %+v", err)
	}

	if d.HasChange("is_default") {
		err = setDefaultQualityProfile(d, m)
		if err != nil {
			return err
		}
	}

	d.SetId(qualityProfileResponse.Profile.Name)
	return ReadQualityProfile(d, m)
}

func setDefaultQualityProfile(d *schema.ResourceData, m interface{}) error {
	qualityProfile := m.(*ProviderConfiguration).sonarQubeURL
	qualityProfile.Path = "api/qualityprofiles/set_default"
	qualityProfile.RawQuery = url.Values{
		"name":     []string{d.Get("name").(string)},
		"language": []string{d.Get("language").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		qualityProfile.String(),
		http.StatusOK,
		"CreateQualityProfile",
	)
	defer resp.Body.Close()
	return err
}

func ReadQualityProfile(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualityprofiles/search"

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"ReadQualityProfile",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	getQualityProfileResponse := GetQualityProfileList{}
	err = json.NewDecoder(resp.Body).Decode(&getQualityProfileResponse)
	if err != nil {
		return fmt.Errorf("ReadQualityProfile: Failed to decode json into struct: %+v", err)
	}

	for _, value := range getQualityProfileResponse.Profiles {
		if d.Id() == value.Name {
			d.SetId(value.Name)
			d.Set("language", value.Language)
			return nil
		}
	}

	return fmt.Errorf("ReadQualityProfile: Failed to find project: %+v", d.Id())
}

func DeleteQualityProfile(d *schema.ResourceData, m interface{}) error {
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
		"DeleteQualityProfile",
	)
	if err != nil {
		return fmt.Errorf("DeleteQualityProfile: Failed to delete quality profile: %+v", err)
	}
	defer resp.Body.Close()

	return nil

}

func ImportQualityProfile(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := ReadQualityProfile(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
