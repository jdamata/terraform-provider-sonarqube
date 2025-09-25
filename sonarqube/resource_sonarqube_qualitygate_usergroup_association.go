package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetQualityGateUsergroupAssociation for unmarshalling response body from getting quality gate association
type GetQualityGateUsergroupAssociation struct {
	Paging Paging                                       `json:"paging"`
	Groups []GetQualityGateUsergroupAssociationProjects `json:"groups,omitempty"`
	Users  []GetQualityGateUsergroupAssociationProjects `json:"users,omitempty"`
}

// GetQualityGateUsergroupAssociationProjects used in GetQualityGateUsergroupAssociation
type GetQualityGateUsergroupAssociationProjects struct {
	Login       string `json:"login,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Selected    bool   `json:"selected"`
}

// Returns the resource represented by this file.
func resourceSonarqubeQualityGateUsergroupAssociation() *schema.Resource {
	return &schema.Resource{
		Description: `Provides a Sonarqube Quality Gate Usergroup association resource. This can be used to associate a Quality Gate to an User or to a Group.
The feature is available on SonarQube 9.2 or newer.`,
		Create: resourceSonarqubeQualityGateUsergroupAssociationCreate,
		Read:   resourceSonarqubeQualityGateUsergroupAssociationRead,
		Delete: resourceSonarqubeQualityGateUsergroupAssociationDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"login_name", "group_name"},
				Description:  "The name of the User to associate. Either `group_name` or `login_name` should be provided.",
			},
			"group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"login_name", "group_name"},
				Description:  "The name of the Group to associate. Either `group_name` or `login_name` should be provided.",
			},
			"gatename": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Quality Gate",
			},
		},
	}
}

func resourceSonarqubeQualityGateUsergroupAssociationCreate(d *schema.ResourceData, m interface{}) error {
	if err := checkGatePermissionFeatureSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	rawQuery := url.Values{
		"gateName": []string{d.Get("gatename").(string)},
	}

	if _, ok := d.GetOk("login_name"); ok {
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/add_user"
		rawQuery.Add("login", d.Get("login_name").(string))
	} else {
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/add_group"
		rawQuery.Add("groupName", d.Get("group_name").(string))
	}

	sonarQubeURL.RawQuery = rawQuery.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeQualityGateUsergroupAssociationCreate",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityGateUsergroupAssociationCreate: Failed creating Sonarqube quality gate usergroup association for quality gate '%s': %+v", d.Get("gatename").(string), err)
	}
	defer resp.Body.Close()

	if _, ok := d.GetOk("login_name"); ok {
		d.SetId(createGatePermissionId(d.Get("gatename").(string), "user", d.Get("login_name").(string)))
	} else {
		d.SetId(createGatePermissionId(d.Get("gatename").(string), "group", d.Get("group_name").(string)))
	}
	return resourceSonarqubeQualityGateUsergroupAssociationRead(d, m)
}

func resourceSonarqubeQualityGateUsergroupAssociationRead(d *schema.ResourceData, m interface{}) error {
	if err := checkGatePermissionFeatureSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.RawQuery = url.Values{
		"gateName": []string{d.Get("gatename").(string)},
		"selected": []string{"selected"},
	}.Encode()

	if _, ok := d.GetOk("login_name"); ok {
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/search_users"
	} else {
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/search_groups"
	}

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeQualityGateUsergroupAssociationRead",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityGateUsergroupAssociationRead: Failed to call quality gate usergroup association api: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityGateUsergroupAssociationReadResponse := GetQualityGateUsergroupAssociation{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateUsergroupAssociationReadResponse)
	if err != nil {
		return fmt.Errorf("resourceQualityGateRead: Failed to decode json into struct: %+v", err)
	}

	if _, ok := d.GetOk("login_name"); ok {
		// Loop over all groups to see if the group we need exists.
		login := d.Get("login_name").(string)
		for _, value := range qualityGateUsergroupAssociationReadResponse.Users {
			if strings.EqualFold(value.Login, login) {
				return d.Set("login_name", value.Login)
			}
		}
	} else {
		// Loop over all groups to see if the group we need exists.
		groupName := d.Get("group_name").(string)
		for _, value := range qualityGateUsergroupAssociationReadResponse.Groups {
			if strings.EqualFold(value.Name, groupName) {
				return nil
			}
		}
	}
	return fmt.Errorf("resourceSonarqubeQualityGateUsergroupAssociationRead: Failed to call quality gate usergroup association api: %+v", err)
}

func resourceSonarqubeQualityGateUsergroupAssociationDelete(d *schema.ResourceData, m interface{}) error {
	if err := checkGatePermissionFeatureSupport(m.(*ProviderConfiguration)); err != nil {
		return err
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	rawQuery := url.Values{
		"gateName": []string{d.Get("gatename").(string)},
	}

	if _, ok := d.GetOk("login_name"); ok {
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/remove_user"
		rawQuery.Add("login", d.Get("login_name").(string))
	} else {
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/qualitygates/remove_group"
		rawQuery.Add("groupName", d.Get("group_name").(string))
	}

	sonarQubeURL.RawQuery = rawQuery.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeQualityGateProjectAssociationDelete",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeQualityGateUsergroupAssociationDelete: Failed to call quality gate usergroup association api: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}

func createGatePermissionId(gateName string, targetType string, target string) string {
	return gateName + "[" + targetType + "/" + target + "]"
}

func checkGatePermissionFeatureSupport(conf *ProviderConfiguration) error {
	minimumVersion, _ := version.NewVersion("9.2")
	if conf.sonarQubeVersion.LessThan(minimumVersion) {
		return fmt.Errorf("minimum required SonarQube version for quality gate permissions is %s", minimumVersion)
	}
	return nil
}
