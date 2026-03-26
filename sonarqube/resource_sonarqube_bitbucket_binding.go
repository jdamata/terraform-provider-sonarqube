package sonarqube

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Returns the resource represented by this file.
func resourceSonarqubeBitbucketBinding() *schema.Resource {
	return &schema.Resource{
		Description: `Provides a Sonarqube Bitbucket Data Center binding resource. This can be used to create and manage the binding between a
Bitbucket Data Center repository and a SonarQube project`,
		Create: resourceSonarqubeBitbucketBindingCreate,
		// You can update any project binding with the same API call as the CREATE
		Update: resourceSonarqubeBitbucketBindingCreate,
		Read:   resourceSonarqubeBitbucketBindingRead,
		Delete: resourceSonarqubeBitbucketBindingDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeBitbucketBindingImport,
		},
		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"alm_setting": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bitbucket Server ALM setting key",
			},
			"monorepo": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "false",
				Description: "Is this project part of a monorepo. Default value: false",
			},
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SonarQube project key. Changing this will force a new resource to be created",
			},
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bitbucket Server repository key (Bitbucket Project Key)",
			},
			"slug": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bitbucket repository slug",
			},
		},
	}
}

func checkBitbucketBindingSupport(conf *ProviderConfiguration) error {
	if strings.ToLower(conf.sonarQubeEdition) == "community" {
		return fmt.Errorf("Bitbucket Bindings are not supported in the Community edition of SonarQube. You are using: SonarQube %s version %s", conf.sonarQubeEdition, conf.sonarQubeVersion)
	}
	return nil
}

func resourceSonarqubeBitbucketBindingCreate(d *schema.ResourceData, m interface{}) error {
	if err := checkBitbucketBindingSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/set_bitbucket_binding"

	sonarQubeURL.RawQuery = url.Values{
		"almSetting": []string{d.Get("alm_setting").(string)},
		"monorepo":   []string{d.Get("monorepo").(string)},
		"project":    []string{d.Get("project").(string)},
		"repository": []string{d.Get("repository").(string)},
		"slug":       []string{d.Get("slug").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeBitbucketBindingCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	id := fmt.Sprintf("%v/%v", d.Get("project").(string), d.Get("repository").(string))
	d.SetId(id)

	return resourceSonarqubeBitbucketBindingRead(d, m)
}

func resourceSonarqubeBitbucketBindingRead(d *schema.ResourceData, m interface{}) error {
	if err := checkBitbucketBindingSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

	idSlice := strings.SplitN(d.Id(), "/", 2)
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/get_binding"
	sonarQubeURL.RawQuery = url.Values{
		"project": []string{idSlice[0]},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeBitbucketBindingRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	BindingReadResponse := GetBinding{}
	err = json.NewDecoder(resp.Body).Decode(&BindingReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeBitbucketBindingRead: Failed to decode json into struct: %+v", err)
	}
	// Check if the binding matches the expected repository and ALM type
	if idSlice[1] == BindingReadResponse.Repository && BindingReadResponse.Alm == "bitbucket" {
		errs := []error{}
		errs = append(errs, d.Set("project", idSlice[0]))
		errs = append(errs, d.Set("repository", idSlice[1]))
		errs = append(errs, d.Set("slug", BindingReadResponse.Slug))
		errs = append(errs, d.Set("alm_setting", BindingReadResponse.Key))
		errs = append(errs, d.Set("monorepo", strconv.FormatBool(BindingReadResponse.Monorepo)))

		return errors.Join(errs...)
	}
	return fmt.Errorf("resourceSonarqubeBitbucketBindingRead: Failed to find bitbucket binding: %+v", d.Id())
}

func resourceSonarqubeBitbucketBindingDelete(d *schema.ResourceData, m interface{}) error {
	if err := checkBitbucketBindingSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/delete_binding"
	sonarQubeURL.RawQuery = url.Values{
		"project": []string{d.Get("project").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeBitbucketBindingDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeBitbucketBindingImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeBitbucketBindingRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
