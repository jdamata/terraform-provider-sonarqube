package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetGroupPermissions struct
type GetGroupPermissions struct {
	Paging Paging            `json:"paging"`
	Groups []GroupPermission `json:"groups"`
}

// GroupPermission struct
type GroupPermission struct {
	Id          string   `json:"id"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions,omitempty"`
}

// Returns the resource represented by this file.
func resourceSonarqubePermissions() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Sonarqube Permissions resource. This can be used to manage global and project permissions.",
		Create:      resourceSonarqubePermissionsCreate,
		Read:        resourceSonarqubePermissionsRead,
		Update:      resourceSonarqubePermissionsUpdate,
		Delete:      resourceSonarqubePermissionsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubePermissionsImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"login_name", "group_name"},
				Description:  "The name of the user that should get the specified permissions. Changing this forces a new resource to be created. Cannot be used with `group_name`.",
			},
			"group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"login_name", "group_name"},
				Description:  "The name of the Group that should get the specified permissions. Changing this forces a new resource to be created. Cannot be used with `login_name`",
			},
			"project_key": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"template_id", "template_name"},
				Description:   "Specify if you want to apply project level permissions. Changing this forces a new resource to be created. Cannot be used with `template_id & template_name`",
			},
			"template_id": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"project_key", "template_name"},
				Description:   "Specify if you want to apply the permissions to a permission template. Changing this forces a new resource to be created. Cannot be used with `project_key & template_name`",
			},
			"template_name": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"project_key", "template_id"},
				Description:   "Specify if you want to apply the permissions to a permission template. Changing this forces a new resource to be created. Cannot be used with `project_key & template_id`",
			},
			"permissions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of permissions that should be applied. Possible values are: `admin`, `codeviewer`, `issueadmin`, `securityhotspotadmin`, `scan`, `user`.",
			},
		},
	}
}

func resourceSonarqubePermissionsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) > 2 {
		return nil, fmt.Errorf("invalid import format, expected 'principal(:scope)' where principal is login_name or group_name and scope is project_key, template_id or template_name with optional prefix")
	}

	principal := parts[0]
	scope := ""
	if len(parts) == 2 {
		scope = parts[1]
	}

	// Determine if principal is a user or group by checking if it exists as a user
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/users/search"
	RawQuery := url.Values{
		"ps": []string{"100"},
		"q":  []string{principal},
	}
	sonarQubeURL.RawQuery = RawQuery.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubePermissionsImport",
	)

	if err != nil {
		return nil, fmt.Errorf("error searching for user during import: %+v", err)
	}
	defer resp.Body.Close()

	users := GetUser{}
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user response: %+v", err)
	}

	isUser := false
	for _, user := range users.Users {
		if strings.EqualFold(user.Login, principal) {
			isUser = true
			d.Set("login_name", user.Login)
			break
		}
	}

	if !isUser {
		// Assume it's a group
		d.Set("group_name", principal)
	}

	if scope != "" {
		// Determine the scope type (project_key, template_id, template_name)
		if strings.HasPrefix(scope, "p_") {
			d.Set("project_key", scope[2:])
		} else if strings.HasPrefix(scope, "t_") {
			d.Set("template_id", scope[2:])
		} else if strings.HasPrefix(scope, "tn_") {
			d.Set("template_name", scope[3:])
		}
	} else {
		scope = "global"
	}

	// Generate a deterministic ID
	if isUser {
		d.SetId(fmt.Sprintf("user-%s-%s-permissions", principal, scope))
	} else {
		d.SetId(fmt.Sprintf("group-%s-%s-permissions", principal, scope))
	}

	// Read the current state
	if err := resourceSonarqubePermissionsRead(d, m); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func resourceSonarqubePermissionsCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	permissions := expandPermissions(d.Get("permissions"))

	var principalName, scopeValue string

	if loginName, ok := d.GetOk("login_name"); ok {
		principalName = loginName.(string)
	} else if groupName, ok := d.GetOk("group_name"); ok {
		principalName = groupName.(string)
	}

	if projectKey, ok := d.GetOk("project_key"); ok {
		scopeValue = "p_" + projectKey.(string)
	} else if templateID, ok := d.GetOk("template_id"); ok {
		scopeValue = "t_" + templateID.(string)
	} else if templateName, ok := d.GetOk("template_name"); ok {
		scopeValue = "tn_" + templateName.(string)
	} else {
		scopeValue = "global"
	}

	// build the base query
	RawQuery := url.Values{}

	// if the permissions should be applied to a project
	// we append the project_key to the request
	if projectKey, ok := d.GetOk("project_key"); ok {
		RawQuery.Add("projectKey", projectKey.(string))
	}

	// we use different API endpoints and request params
	// based on the target principal type (group or user)
	// and if its a direct or template permission
	if _, ok := d.GetOk("login_name"); ok {
		// user permission
		RawQuery.Add("login", d.Get("login_name").(string))
		if templateID, ok := d.GetOk("template_id"); ok {
			// template user permission
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/add_user_to_template"
			RawQuery.Add("templateId", templateID.(string))
			// name provide instead of id
		} else if templateName, ok := d.GetOk("template_name"); ok {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/add_user_to_template"
			RawQuery.Add("templateName", templateName.(string))
		} else {
			// direct user permission
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/add_user"
		}

		d.SetId(fmt.Sprintf("user-%s-%s-permissions", principalName, scopeValue))
	} else {
		// group permission
		RawQuery.Add("groupName", d.Get("group_name").(string))
		if templateID, ok := d.GetOk("template_id"); ok {
			// template user permission
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/add_group_to_template"
			RawQuery.Add("templateId", templateID.(string))
			// name provide instead of id
		} else if templateName, ok := d.GetOk("template_name"); ok {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/add_group_to_template"
			RawQuery.Add("templateName", templateName.(string))
		} else {
			// direct user permission
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/add_group"
		}

		d.SetId(fmt.Sprintf("group-%s-%s-permissions", principalName, scopeValue))
	}

	// loop through all permissions that should be applied
	for _, permission := range permissions {
		CurrentRawQuery := RawQuery
		CurrentRawQuery.Del("permission")
		CurrentRawQuery.Add("permission", permission)
		sonarQubeURL.RawQuery = CurrentRawQuery.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"POST",
			sonarQubeURL.String(),
			http.StatusNoContent,
			"resourceSonarqubePermissionsCreate",
		)
		if err != nil {
			return fmt.Errorf("error creating Sonarqube permission: %+v", err)
		}
		defer resp.Body.Close()
	}

	return resourceSonarqubePermissionsRead(d, m)
}

func resourceSonarqubePermissionsRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL

	// build the base query
	RawQuery := url.Values{
		// set the page size to 100
		"ps": []string{"100"},
	}

	// if the permissions should be applied to a project
	// we append the project_key to the request
	if projectKey, ok := d.GetOk("project_key"); ok {
		RawQuery.Add("projectKey", projectKey.(string))
	}

	// we use different API endpoints and request params
	// based on the target principal type (group or user)
	// and if its a direct or template permission
	if loginName, ok := d.GetOk("login_name"); ok {
		// permission target is USER
		if templateID, ok := d.GetOk("template_id"); ok {
			// template user permission
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/template_users"
			RawQuery.Add("templateId", templateID.(string))
			// name provide instead of id
		} else if templateName, ok := d.GetOk("template_name"); ok {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/template_users"
			RawQuery.Add("templateName", templateName.(string))
		} else {
			// direct user permission
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/users"
			RawQuery.Add("q", loginName.(string))
		}
		sonarQubeURL.RawQuery = RawQuery.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"GET",
			sonarQubeURL.String(),
			http.StatusOK,
			"resourceSonarqubePermissionsRead",
		)
		if err != nil {
			return fmt.Errorf("error reading Sonarqube permissions: %+v", err)
		}
		defer resp.Body.Close()

		// Decode response into struct
		users := GetUser{}
		err = json.NewDecoder(resp.Body).Decode(&users)
		if err != nil {
			return fmt.Errorf("resourceSonarqubePermissionsRead: Failed to decode json into struct: %+v", err)
		}

		// Loop over all groups to see if the group we need exists.
		for _, value := range users.Users {
			if strings.EqualFold(value.Login, loginName.(string)) {
				d.Set("login_name", value.Login)
				d.Set("permissions", flattenPermissions(&value.Permissions))
				return nil
			}
		}

	} else {
		// permission target is GROUP
		groupName := d.Get("group_name").(string)

		if templateID, ok := d.GetOk("template_id"); ok {
			// template group permission
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/template_groups"
			RawQuery.Add("templateId", templateID.(string))
		} else if templateName, ok := d.GetOk("template_name"); ok {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/template_groups"
			RawQuery.Add("templateName", templateName.(string))
		} else {
			// direct group permission
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/groups"
			RawQuery.Add("q", groupName)
		}
		sonarQubeURL.RawQuery = RawQuery.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"GET",
			sonarQubeURL.String(),
			http.StatusOK,
			"resourceSonarqubePermissionsRead",
		)
		if err != nil {
			return fmt.Errorf("error reading Sonarqube permissions: %+v", err)
		}
		defer resp.Body.Close()

		// Decode response into struct
		groups := GetGroupPermissions{}
		err = json.NewDecoder(resp.Body).Decode(&groups)
		if err != nil {
			return fmt.Errorf("resourceSonarqubePermissionsRead: Failed to decode json into struct: %+v", err)
		}

		// Loop over all groups to see if the group we need exists.
		for _, value := range groups.Groups {
			if strings.EqualFold(value.Name, groupName) {
				d.Set("group_name", value.Name)
				d.Set("permissions", flattenPermissions(&value.Permissions))
				return nil
			}
		}
	}

	return fmt.Errorf("resourceSonarqubePermissionsRead: Unable to find group permissions for group: %+v", d.Id())
}

func resourceSonarqubePermissionsUpdate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeBasePath := sonarQubeURL.Path

	currentFlatPermissions, targetFlatPermissions := d.GetChange("permissions")
	currentPermissions := expandPermissions(currentFlatPermissions)
	targetPermissions := expandPermissions(targetFlatPermissions)

	RawQuery := url.Values{
		"ps": []string{"100"},
	}

	if projectKey, ok := d.GetOk("project_key"); ok {
		RawQuery.Add("projectKey", projectKey.(string))
	}

	if loginName, ok := d.GetOk("login_name"); ok {

		toAddPermissions, toRemovePermissions := calculatePermissionChanges(currentPermissions, targetPermissions)

		RawQuery = url.Values{}

		if projectKey, ok := d.GetOk("project_key"); ok {
			RawQuery.Add("projectKey", projectKey.(string))
		}

		RawQuery.Add("login", loginName.(string))
		if templateID, ok := d.GetOk("template_id"); ok {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeBasePath, "/") + "/api/permissions/remove_user_from_template"
			RawQuery.Add("templateId", templateID.(string))
		} else if templateName, ok := d.GetOk("template_name"); ok {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeBasePath, "/") + "/api/permissions/remove_user_from_template"
			RawQuery.Add("templateName", templateName.(string))
		} else {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeBasePath, "/") + "/api/permissions/remove_user"
		}

		for _, perm := range toRemovePermissions {
			RawQuery.Set("permission", perm)
			sonarQubeURL.RawQuery = RawQuery.Encode()

			resp, err := httpRequestHelper(
				m.(*ProviderConfiguration).httpClient,
				"POST",
				sonarQubeURL.String(),
				http.StatusNoContent,
				"resourceSonarqubePermissionsUpdate",
			)
			if err != nil {
				return fmt.Errorf("resourceSonarqubePermissionsUpdate: Error removing Sonarqube permissions: %+v", err)
			}
			defer resp.Body.Close()
		}

		RawQuery = url.Values{}

		if projectKey, ok := d.GetOk("project_key"); ok {
			RawQuery.Add("projectKey", projectKey.(string))
		}

		RawQuery.Add("login", loginName.(string))
		if templateID, ok := d.GetOk("template_id"); ok {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeBasePath, "/") + "/api/permissions/add_user_to_template"
			RawQuery.Add("templateId", templateID.(string))
		} else if templateName, ok := d.GetOk("template_name"); ok {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeBasePath, "/") + "/api/permissions/add_user_to_template"
			RawQuery.Add("templateName", templateName.(string))
		} else {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeBasePath, "/") + "/api/permissions/add_user"
		}

		for _, perm := range toAddPermissions {
			RawQuery.Add("permission", perm)
			sonarQubeURL.RawQuery = RawQuery.Encode()

			resp, err := httpRequestHelper(
				m.(*ProviderConfiguration).httpClient,
				"POST",
				sonarQubeURL.String(),
				http.StatusNoContent,
				"resourceSonarqubePermissionsUpdate",
			)
			if err != nil {
				return fmt.Errorf("resourceSonarqubePermissionsUpdate: Error adding Sonarqube permissions: %+v", err)
			}
			defer resp.Body.Close()
		}
	} else if groupName, ok := d.GetOk("group_name"); ok {
		toAddPermissions, toRemovePermissions := calculatePermissionChanges(currentPermissions, targetPermissions)

		RawQuery = url.Values{}

		if projectKey, ok := d.GetOk("project_key"); ok {
			RawQuery.Add("projectKey", projectKey.(string))
		}

		RawQuery.Add("groupName", groupName.(string))
		if templateID, ok := d.GetOk("template_id"); ok {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeBasePath, "/") + "/api/permissions/remove_group_from_template"
			RawQuery.Add("templateId", templateID.(string))
		} else if templateName, ok := d.GetOk("template_name"); ok {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeBasePath, "/") + "/api/permissions/remove_group_from_template"
			RawQuery.Add("templateName", templateName.(string))
		} else {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeBasePath, "/") + "/api/permissions/remove_group"
		}

		for _, perm := range toRemovePermissions {
			RawQuery.Set("permission", perm)
			sonarQubeURL.RawQuery = RawQuery.Encode()

			resp, err := httpRequestHelper(
				m.(*ProviderConfiguration).httpClient,
				"POST",
				sonarQubeURL.String(),
				http.StatusNoContent,
				"resourceSonarqubePermissionsUpdate",
			)
			if err != nil {
				return fmt.Errorf("resourceSonarqubePermissionsUpdate: Error removing Sonarqube permissions: %+v", err)
			}
			defer resp.Body.Close()
		}

		RawQuery = url.Values{}

		if projectKey, ok := d.GetOk("project_key"); ok {
			RawQuery.Add("projectKey", projectKey.(string))
		}

		RawQuery.Add("groupName", groupName.(string))
		if templateID, ok := d.GetOk("template_id"); ok {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeBasePath, "/") + "/api/permissions/add_group_to_template"
			RawQuery.Add("templateId", templateID.(string))
		} else if templateName, ok := d.GetOk("template_name"); ok {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeBasePath, "/") + "/api/permissions/add_group_to_template"
			RawQuery.Add("templateName", templateName.(string))
		} else {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeBasePath, "/") + "/api/permissions/add_group"
		}

		for _, perm := range toAddPermissions {
			RawQuery.Add("permission", perm)
			sonarQubeURL.RawQuery = RawQuery.Encode()

			resp, err := httpRequestHelper(
				m.(*ProviderConfiguration).httpClient,
				"POST",
				sonarQubeURL.String(),
				http.StatusNoContent,
				"resourceSonarqubePermissionsUpdate",
			)
			if err != nil {
				return fmt.Errorf("resourceSonarqubePermissionsUpdate: Error adding Sonarqube permissions: %+v", err)
			}
			defer resp.Body.Close()
		}
	} else {
		return fmt.Errorf("resourceSonarqubePermissionsUpdate: Didn't find any identification")
	}

	return resourceSonarqubePermissionsRead(d, m)
}

func resourceSonarqubePermissionsDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	permissions := expandPermissions(d)

	// build the base query
	RawQuery := url.Values{}

	// if the permissions should be applied to a project
	// we append the project_key to the request
	if projectKey, ok := d.GetOk("project_key"); ok {
		RawQuery.Add("projectKey", projectKey.(string))
	}

	// we use different API endpoints and request params
	// based on the target principal type (group or user)
	if _, ok := d.GetOk("login_name"); ok {
		// permission target is USER
		if templateID, ok := d.GetOk("template_id"); ok {
			// template user permission
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/remove_user_from_template"
			RawQuery.Add("templateId", templateID.(string))
		} else if templateName, ok := d.GetOk("template_name"); ok {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/remove_user_from_template"
			RawQuery.Add("templateName", templateName.(string))
		} else {
			// direct user permission
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/remove_user"
		}
		RawQuery.Add("login", d.Get("login_name").(string))
		sonarQubeURL.RawQuery = RawQuery.Encode()

	} else {
		// permission target is GROUP
		if templateID, ok := d.GetOk("template_id"); ok {
			// template group permission
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/remove_group_from_template"
			RawQuery.Add("templateId", templateID.(string))
		} else if templateName, ok := d.GetOk("template_name"); ok {
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/remove_group_from_template"
			RawQuery.Add("templateName", templateName.(string))
		} else {
			// direct group permission
			sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/permissions/remove_group"
		}
		RawQuery.Add("groupName", d.Get("group_name").(string))
		sonarQubeURL.RawQuery = RawQuery.Encode()
	}

	// loop through all permissions that should be applied
	for _, permission := range permissions {
		CurrentRawQuery := RawQuery
		CurrentRawQuery.Del("permission")
		CurrentRawQuery.Add("permission", permission)
		sonarQubeURL.RawQuery = CurrentRawQuery.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"POST",
			sonarQubeURL.String(),
			http.StatusNoContent,
			"resourceSonarqubePermissionsDelete",
		)
		if err != nil {
			return fmt.Errorf("error creating Sonarqube permission: %+v", err)
		}
		defer resp.Body.Close()
	}

	return nil
}

func expandPermissions(flatPermissions interface{}) []string {
	expandedPermissions := make([]string, 0)
	for _, permission := range flatPermissions.([]interface{}) {
		expandedPermissions = append(expandedPermissions, permission.(string))
	}

	return expandedPermissions
}

func flattenPermissions(input *[]string) []interface{} {
	flatPermissions := make([]interface{}, 0)
	if input == nil {
		return flatPermissions
	}

	for _, permission := range *input {
		flatPermissions = append(flatPermissions, permission)
	}

	return flatPermissions
}

func calculatePermissionChanges(current, target []string) (toAdd, toRemove []string) {
	currentMap := make(map[string]bool)
	for _, perm := range current {
		currentMap[strings.ToLower(perm)] = true
	}

	targetMap := make(map[string]bool)
	for _, perm := range target {
		targetMap[strings.ToLower(perm)] = true
	}

	for _, perm := range target {
		if !currentMap[strings.ToLower(perm)] {
			toAdd = append(toAdd, perm)
		}
	}

	for _, perm := range current {
		if !targetMap[strings.ToLower(perm)] {
			toRemove = append(toRemove, perm)
		}
	}

	return toAdd, toRemove
}
