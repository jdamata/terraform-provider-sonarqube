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

// User struct
type User struct {
	Login       string   `json:"login,omitempty"`
	Name        string   `json:"name,omitempty"`
	Email       string   `json:"email,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	IsActive    bool     `json:"active,omitempty"`
	IsLocal     bool     `json:"local,omitempty"`
}

// GetUser for unmarshalling response body where users are retured
type GetUser struct {
	Paging Paging `json:"paging"`
	Users  []User `json:"users"`
}

// CreateUserResponse struct
type CreateUserResponse struct {
	User User `json:"user"`
}

// Returns the resource represented by this file.
func resourceSonarqubeUser() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Sonarqube User resource. This can be used to manage Sonarqube Users.",
		Create:      resourceSonarqubeUserCreate,
		Read:        resourceSonarqubeUserRead,
		Update:      resourceSonarqubeUserUpdate,
		Delete:      resourceSonarqubeUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeUserImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The login name of the User to create. Changing this forces a new resource to be created.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the User to create. Changing this forces a new resource to be created.",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email of the User to create.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password of User to create. This is only used if the user is of type `local`.",
			},
			"is_local": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "`True` if the User should be of type `local`. Defaults to `true`.",
			},
		},
	}
}

func resourceSonarqubeUserCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/users/create"

	isLocal := d.Get("is_local").(bool)

	rawQuery := url.Values{
		"login": []string{d.Get("login_name").(string)},
		"name":  []string{d.Get("name").(string)},
		"local": []string{strconv.FormatBool(isLocal)},
	}

	if password, ok := d.GetOk("password"); ok {
		rawQuery.Add("password", password.(string))
	}

	if email, ok := d.GetOk("email"); ok {
		rawQuery.Add("email", email.(string))
	}

	sonarQubeURL.RawQuery = rawQuery.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeUserCreate",
	)
	if err != nil {
		return fmt.Errorf("error creating Sonarqube user: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	userResponse := CreateUserResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeUserCreate: Failed to decode json into struct: %+v", err)
	}

	if userResponse.User.Login != "" {
		d.SetId(userResponse.User.Login)
	} else {
		return fmt.Errorf("resourceSonarqubeUserCreate: Create response didn't contain the user login")
	}

	return resourceSonarqubeUserRead(d, m)
}

func resourceSonarqubeUserRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/users/search"

	sonarQubeURL.RawQuery = url.Values{
		"ps": []string{"500"},
		"q":  []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeUserRead",
	)
	if err != nil {
		return fmt.Errorf("error reading Sonarqube user: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	userResponse := GetUser{}
	err = json.NewDecoder(resp.Body).Decode(&userResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeUserCreate: Failed to decode json into struct: %+v", err)
	}

	// Loop over all users to see if the current user exists.
	for _, value := range userResponse.Users {
		if d.Id() == value.Login {
			d.SetId(value.Login)
			errs := []error{}
			errs = append(errs, d.Set("login_name", value.Login))
			errs = append(errs, d.Set("name", value.Name))
			errs = append(errs, d.Set("email", value.Email))
			errs = append(errs, d.Set("is_local", value.IsLocal))
			return errors.Join(errs...)
		}
	}

	return fmt.Errorf("resourceSonarqubeUserRead: Failed to find user: %+v", d.Id())
}

func resourceSonarqubeUserUpdate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURLSubPath := sonarQubeURL.Path
	// handle default updates (api/users/update)
	if d.HasChange("email") {
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURLSubPath, "/") + "/api/users/update"
		sonarQubeURL.RawQuery = url.Values{
			"login": []string{d.Id()},
			"email": []string{d.Get("email").(string)},
		}.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"POST",
			sonarQubeURL.String(),
			http.StatusOK,
			"resourceSonarqubeUserUpdate",
		)
		if err != nil {
			return fmt.Errorf("error updating Sonarqube user: %+v", err)
		}
		defer resp.Body.Close()
	}

	// handle password updates (api/users/change_password)
	if d.HasChange("password") {
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURLSubPath, "/") + "/api/users/change_password"
		oldPassword, newPassword := d.GetChange("password")
		sonarQubeURL.RawQuery = url.Values{
			"login":            []string{d.Id()},
			"password":         []string{newPassword.(string)},
			"previousPassword": []string{oldPassword.(string)},
		}.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"POST",
			sonarQubeURL.String(),
			http.StatusNoContent,
			"resourceSonarqubeUserUpdate",
		)
		if err != nil {
			return fmt.Errorf("error updating Sonarqube user: %+v", err)
		}
		defer resp.Body.Close()
	}

	return resourceSonarqubeUserRead(d, m)
}

func resourceSonarqubeUserDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/users/deactivate"
	sonarQubeURL.RawQuery = url.Values{
		"login":     []string{d.Id()},
		"anonymize": []string{strconv.FormatBool(m.(*ProviderConfiguration).sonarQubeAnonymizeUsers)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeUserDelete",
	)
	if err != nil {
		return fmt.Errorf("error deleting (deactivating) Sonarqube user: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeUserImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeUserRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
