package sonarqube

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Returns the resource represented by this file.
func resourceSonarqubeUserExternalIdentity() *schema.Resource {
	return &schema.Resource{
		Description: `Updates the _external identity_ of a _non local_ Sonarqube User. This can be used to set the _Identity Provider_ which should be used to
authenticate a specific user.

The Sonarqube API currently does not provide an endpoint to read the _external identity_ setting of an user.`,
		Create: resourceSonarqubeUserExternalIdentityCreate,
		Read:   resourceSonarqubeUserExternalIdentityRead,
		Delete: resourceSonarqubeUserExternalIdentityDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The login name of the User to update. Changing this forces a new resource to be created.",
			},
			"external_identity": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The identifier of the User used by the Authentication Provider. Changing this forces a new resource to be created.",
			},
			"external_provider": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The key of the Authentication Provider. The Authentication Provider must be activated on Sonarqube. Changing this forces a new resource to be created.",
			},
		},
	}
}

func resourceSonarqubeUserExternalIdentityCreate(d *schema.ResourceData, m interface{}) error {
	isLocal, err := isLocal(d.Get("login_name").(string), m)
	if err != nil {
		return fmt.Errorf("error updating Sonarqube user: %+v", err)
	}
	if isLocal {
		return fmt.Errorf("error setting external identity: Sonarqube user '%+v' is not 'external'", d.Get("login_name").(string))
	}

	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/users/update_identity_provider"

	rawQuery := url.Values{
		"login":               []string{d.Get("login_name").(string)},
		"newExternalIdentity": []string{d.Get("external_identity").(string)},
		"newExternalProvider": []string{d.Get("external_provider").(string)},
	}

	sonarQubeURL.RawQuery = rawQuery.Encode()

	_, err = httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeUserExternalIdentityCreate",
	)
	if err != nil {
		return fmt.Errorf("error updating Sonarqube user: %+v", err)
	}

	d.SetId(d.Get("login_name").(string))
	errs := []error{}
	errs = append(errs, d.Set("external_identity", d.Get("external_identity").(string)))
	errs = append(errs, d.Set("external_provider", d.Get("external_provider").(string)))

	return errors.Join(errs...)
}

func resourceSonarqubeUserExternalIdentityRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSonarqubeUserExternalIdentityDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func isLocal(login string, m interface{}) (bool, error) {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/users/search"

	sonarQubeURL.RawQuery = url.Values{
		"q": []string{login},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeUserExternalIdentity",
	)
	if err != nil {
		return false, fmt.Errorf("error reading Sonarqube user: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	userResponse := GetUser{}
	err = json.NewDecoder(resp.Body).Decode(&userResponse)
	if err != nil {
		return false, fmt.Errorf("failed to decode json into struct: %+v", err)
	}

	// Loop over all users to find the requested user
	for _, value := range userResponse.Users {
		if login == value.Login {
			return value.IsLocal, nil
		}
	}

	// User not found in response
	return false, fmt.Errorf("failed to find user: %+v", login)
}
