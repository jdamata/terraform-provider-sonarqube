package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type GatePermissionUser struct {
	LoginName string `json:"login,omitempty"`
	Name      string `json:"name,omitempty"`
}

type GetGatePermissionUsersResponse struct {
	Paging Paging               `json:"paging"`
	Users  []GatePermissionUser `json:"users"`
}

type GatePermissionGroup struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type GetGatePermissionGroupResponse struct {
	Paging Paging                `json:"paging"`
	Groups []GatePermissionGroup `json:"groups"`
}

// Returns the resource represented by this file.
func resourceSonarqubeQualityGatePermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeQualityGatePermissionCreate,
		Read:   resourceSonarqubeQualityGatePermissionRead,
		Delete: resourceSonarqubeQualityGatePermissionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeQualityGatePermissionImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"gatename": {
				Type:        schema.TypeString,
				Description: "Name of the quality gate",
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of permission (user or group)",
				Required:    true,
				ForceNew:    true,
			},
			"subject": {
				Type:        schema.TypeString,
				Description: "Login name or group name",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceSonarqubeQualityGatePermissionCreate(d *schema.ResourceData, m interface{}) error {
	if err := checkGatePermissionFeatureSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL

	switch t := d.Get("type").(string); t {
	case "user":
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/add_user"
		sonarQubeURL.RawQuery = url.Values{
			"gateName": []string{d.Get("gatename").(string)},
			"login":    []string{d.Get("subject").(string)},
		}.Encode()
	case "group":
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/add_group"
		sonarQubeURL.RawQuery = url.Values{
			"gateName":  []string{d.Get("gatename").(string)},
			"groupName": []string{d.Get("subject").(string)},
		}.Encode()
	default:
		return fmt.Errorf("Invalid value for 'type' parameter: '%s'", d.Get("type").(string))
	}

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceQualityGatePermissionCreate",
	)
	if err != nil {
		return fmt.Errorf("resourceQualityGatePermissionCreate: error creating Sonarqube %s '%s' permission on quality gate '%s': %w",
			d.Get("type").(string), d.Get("subject").(string), d.Get("gatename").(string), err)
	}
	defer resp.Body.Close()

	d.SetId(createGatePermissionId(d.Get("gatename").(string), d.Get("type").(string), d.Get("subject").(string)))

	return resourceSonarqubeQualityGatePermissionRead(d, m)
}

func resourceSonarqubeQualityGatePermissionRead(d *schema.ResourceData, m interface{}) error {
	if err := checkGatePermissionFeatureSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL

	switch t := d.Get("type").(string); t {
	case "user":
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/search_users"
	case "group":
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/search_groups"
	default:
		return fmt.Errorf("Invalid value for 'type' parameter: '%s'", d.Get("type").(string))
	}
	sonarQubeURL.RawQuery = url.Values{
		"gateName": []string{d.Get("gatename").(string)},
		"q":        []string{d.Get("subject").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeQualityGatePermissionRead",
	)
	if err != nil {
		return fmt.Errorf("error reading Sonarqube permissions on quality gate '%s': %w", d.Get("gatename").(string), err)
	}
	defer resp.Body.Close()

	readSuccess := false

	switch t := d.Get("type").(string); t {
	case "user":
		// Decode response into struct
		readResponse := GetGatePermissionUsersResponse{}
		err = json.NewDecoder(resp.Body).Decode(&readResponse)
		if err != nil {
			return fmt.Errorf("resourceQualityGatePermissionCreate: Failed to read user '%s' permission on quality gate '%s': %+v",
				d.Get("subject").(string), d.Get("gatename").(string), err)
		}
		// Loop over all returned members to see if the member we need exists.
		for _, value := range readResponse.Users {
			if d.Get("subject").(string) == value.LoginName {
				// If it does, set the values of that group membership
				d.SetId(createGatePermissionId(d.Get("gatename").(string), d.Get("type").(string), d.Get("subject").(string)))
				d.Set("subject", value.LoginName)
				readSuccess = true
				break
			}
		}
	case "group":
		// Decode response into struct
		readResponse := GetGatePermissionGroupResponse{}
		err = json.NewDecoder(resp.Body).Decode(&readResponse)
		if err != nil {
			return fmt.Errorf("resourceSonarqubeQualityGatePermissionRead: Failed to read group '%s' permission on quality gate '%s': %+v",
				d.Get("subject").(string), d.Get("gatename").(string), err)
		}
		// Loop over all returned members to see if the member we need exists.
		for _, value := range readResponse.Groups {
			if d.Get("subject").(string) == value.Name {
				// If it does, set the values of that group membership
				d.SetId(createGatePermissionId(d.Get("gatename").(string), d.Get("type").(string), d.Get("subject").(string)))
				d.Set("subject", value.Name)
				readSuccess = true
				break
			}
		}
	default:
		return fmt.Errorf("Invalid value for 'type' parameter: '%s'", d.Get("type").(string))
	}

	if !readSuccess {
		return fmt.Errorf("resourceSonarqubeQualityGatePermissionRead: Failed to read %s '%s' permission on quality gate '%s': not found",
			d.Get("type").(string), d.Get("subject").(string), d.Get("gatename").(string))
	}

	return nil
}

func resourceSonarqubeQualityGatePermissionDelete(d *schema.ResourceData, m interface{}) error {
	if err := checkGatePermissionFeatureSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL

	switch t := d.Get("type").(string); t {
	case "user":
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/remove_user"
		sonarQubeURL.RawQuery = url.Values{
			"gateName": []string{d.Get("gatename").(string)},
			"login":    []string{d.Get("subject").(string)},
		}.Encode()
	case "group":
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/remove_group"
		sonarQubeURL.RawQuery = url.Values{
			"gateName":  []string{d.Get("gatename").(string)},
			"groupName": []string{d.Get("subject").(string)},
		}.Encode()
	default:
		return fmt.Errorf("Invalid value for 'type' parameter: '%s'", d.Get("type").(string))
	}

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceQualityGatePermissionDelete",
	)
	if err != nil {
		return fmt.Errorf("resourceQualityGatePermissionDelete: error removing Sonarqube %s '%s' permission on quality gate '%s': %w",
			d.Get("type").(string), d.Get("subject").(string), d.Get("gatename").(string), err)
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeQualityGatePermissionImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := checkGatePermissionFeatureSupport(m.(*ProviderConfiguration)); err != nil {
		return nil, err
	}

	rgx := regexp.MustCompile(`(.*?)\[(.*?)/(.*?)\]`)
	rs := rgx.FindStringSubmatch(d.Id())
	gateName := rs[1]
	subjectType := rs[2]
	subject := rs[3]

	d.Set("gatename", gateName)
	d.Set("type", subjectType)
	d.Set("subject", subject)

	if err := resourceSonarqubeQualityGatePermissionRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func createGatePermissionId(gateName string, subjectType string, subject string) string {
	return gateName + "[" + subjectType + "/" + subject + "]"
}

func checkGatePermissionFeatureSupport(conf *ProviderConfiguration) error {
	minimumVersion, _ := version.NewVersion("9.2")
	if conf.sonarQubeVersion.LessThan(minimumVersion) {
		return fmt.Errorf("Minimum required SonarQube version for quality gate permissions is %s", minimumVersion)
	}
	return nil
}
