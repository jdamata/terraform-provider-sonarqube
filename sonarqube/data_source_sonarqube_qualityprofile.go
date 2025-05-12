package sonarqube

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"net/url"
	"strings"
)

func dataSourceSonarqubeQualityProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a Sonarqube qualityprofile resource",
		Read:        dataSourceSonarqubeQualityProfileRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the Quality Profile",
				ExactlyOneOf: []string{"name", "key"},
				RequiredWith: []string{"language"},
			},
			"key": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The key of the Quality Profile",
				ExactlyOneOf: []string{"name", "key"},
			},
			"language": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Quality Profile language",
				RequiredWith: []string{"name"},
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Quality Profile default",
			},
		},
	}
}

func dataSourceSonarqubeQualityProfileRead(d *schema.ResourceData, m interface{}) error {
	key, hasKey := d.GetOk("key")
	name, hasName := d.GetOk("name")
	language, hasLanguage := d.GetOk("language")

	if hasKey && (hasName || hasLanguage) {
		return fmt.Errorf("dataSourceSonarqubeQualityProfileRead: when 'key' is set, 'name' and 'language' should not be set")
	}

	if (hasName && !hasLanguage) || (!hasName && hasLanguage) {
		return fmt.Errorf("dataSourceSonarqubeQualityProfileRead: 'name' and 'language' must be set together")
	}

	if !hasKey && !hasName {
		return fmt.Errorf("dataSourceSonarqubeQualityProfileRead: either 'key' or both 'name' and 'language' must be set")
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualityprofiles/search"

	RawQuery := url.Values{}
	if !hasKey {
		RawQuery.Set("qualityProfile", name.(string))
		RawQuery.Set("language", language.(string))
	}
	sonarQubeURL.RawQuery = RawQuery.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"dataSourceSonarqubeQualityProfileRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	getQualityProfileResponse := GetQualityProfileList{}
	err = json.NewDecoder(resp.Body).Decode(&getQualityProfileResponse)
	if err != nil {
		return fmt.Errorf("dataSourceSonarqubeQualityProfileRead: Failed to decode json into struct: %+v", err)
	}

	// bundle identifiable information for informative error
	qualityProfileIdent := fmt.Sprintf("key: %s", key.(string))
	if !hasKey {
		qualityProfileIdent = fmt.Sprintf("name: %s & language: %s", name.(string), language.(string))
	}

	if len(getQualityProfileResponse.Profiles) <= 0 {
		return fmt.Errorf("dataSourceSonarqubeQualityProfileRead: failed to find quality profile with name %s", qualityProfileIdent)
	}

	var foundQualityProfile GetQualityProfile
	if hasKey {
		for _, qualityProfile := range getQualityProfileResponse.Profiles {
			if key.(string) == qualityProfile.Key {
				foundQualityProfile = qualityProfile
				break
			}
		}
	} else {
		if len(getQualityProfileResponse.Profiles) > 1 {
			return fmt.Errorf("dataSourceSonarqubeQualityProfileRead: found more than one quality profile with %s", qualityProfileIdent)
		}
		foundQualityProfile = getQualityProfileResponse.Profiles[0]
	}

	if foundQualityProfile.Key != "" {
		d.SetId(foundQualityProfile.Key)
		var errs []error
		errs = append(errs, d.Set("name", foundQualityProfile.Name))
		errs = append(errs, d.Set("language", foundQualityProfile.Language))
		errs = append(errs, d.Set("key", foundQualityProfile.Key))
		errs = append(errs, d.Set("is_default", foundQualityProfile.IsDefault))
		return errors.Join(errs...)
	}

	return fmt.Errorf("dataSourceSonarqubeQualityProfileRead: Failed to find quality profile %s", qualityProfileIdent)
}
