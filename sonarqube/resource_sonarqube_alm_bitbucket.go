package sonarqube

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// GetAlmBitbucket for unmarshalling response body from alm list definitions. With only bitbucket populated
type GetAlmBitbucket struct {
	Bitbucket []struct {
		Key                 string `json:"key"`
		URL                 string `json:"url"`
		PersonalAccessToken string `json:"personalAccessToken,omitempty"`
	} `json:"bitbucket"`
}

// Returns the resource represented by this file.
func resourceSonarqubeAlmBitbucket() *schema.Resource {
	return &schema.Resource{
		Description: `Provides a Sonarqube Bitbucket Data Center Alm/Devops Platform Integration resource. This can be used to create and manage a Alm/Devops
Platform Integration for Bitbucket Data Center.`,
		Create: resourceSonarqubeAlmBitbucketCreate,
		Read:   resourceSonarqubeAlmBitbucketRead,
		Update: resourceSonarqubeAlmBitbucketUpdate,
		Delete: resourceSonarqubeAlmBitbucketDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"key": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 200)),
				Description:      "Unique key of the Bitbucket instance setting. Maximum length: 200",
			},
			"personal_access_token": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 2000)),
				Description:      "Bitbucket personal access token. Maximum length: 2000",
			},
			"url": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 2000)),
				Description:      "Bitbucket server API URL. Maximum length: 2000",
			},
		},
	}
}

func resourceSonarqubeAlmBitbucketCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/create_bitbucket"

	sonarQubeURL.RawQuery = url.Values{
		"key":                 []string{d.Get("key").(string)},
		"personalAccessToken": []string{d.Get("personal_access_token").(string)},
		"url":                 []string{d.Get("url").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeAlmBitbucketCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d.SetId(d.Get("key").(string))

	return resourceSonarqubeAlmBitbucketRead(d, m)
}

func resourceSonarqubeAlmBitbucketRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/list_definitions"

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeAlmBitbucketRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	AlmBitbucketReadResponse := GetAlmBitbucket{}
	err = json.NewDecoder(resp.Body).Decode(&AlmBitbucketReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeAlmBitbucketRead: Failed to decode json into struct: %+v", err)
	}
	// Loop over all Bitbucket instances to see if the Alm instance exists.
	for _, value := range AlmBitbucketReadResponse.Bitbucket {
		if d.Id() == value.Key {
			errKey := d.Set("key", value.Key)
			errUrl := d.Set("url", value.URL)
			// The personal_access_token is a secured property that is not returned
			// d.Set("personal_access_token", value.PersonalAccessToken)
			return errors.Join(errKey, errUrl)
		}
	}
	return fmt.Errorf("resourceSonarqubeAlmBitbucketRead: Failed to find bitbucket alm setting: %+v", d.Id())
}

func resourceSonarqubeAlmBitbucketUpdate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/update_bitbucket"
	sonarQubeURL.RawQuery = url.Values{
		"key":                 []string{d.Id()},
		"newKey":              []string{d.Get("key").(string)},
		"personalAccessToken": []string{d.Get("personal_access_token").(string)},
		"url":                 []string{d.Get("url").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeAlmBitbucketUpdate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return resourceSonarqubeAlmBitbucketRead(d, m)
}

func resourceSonarqubeAlmBitbucketDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/delete"
	sonarQubeURL.RawQuery = url.Values{
		"key": []string{d.Get("key").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeAlmBitbucketDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
