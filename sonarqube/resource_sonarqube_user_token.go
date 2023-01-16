package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetTokens struct
type GetTokens struct {
	Login  string  `json:"login,omitempty"`
	Tokens []Token `json:"userTokens,omitempty"`
}

// Token struct
type Token struct {
	Login          string `json:"login,omitempty"`
	Name           string `json:"name,omitempty"`
	Token          string `json:"token,omitempty"`
	ExpirationDate string `json:"expirationDate,omitempty"`
}

// Returns the resource represented by this file.
func resourceSonarqubeUserToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeUserTokenCreate,
		Read:   resourceSonarqubeUserTokenRead,
		Delete: resourceSonarqubeUserTokenDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceSonarqubeUserTokenCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/user_tokens/generate"

	rawQuery := url.Values{
		"login": []string{d.Get("login_name").(string)},
		"name":  []string{d.Get("name").(string)},
	}

	if _, ok := d.GetOk("expiration_date"); ok {
		rawQuery.Add("expirationDate", d.Get("expiration_date").(string))
	}

	sonarQubeURL.RawQuery = rawQuery.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeUserTokenCreate",
	)
	if err != nil {
		return fmt.Errorf("error creating Sonarqube user token: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	tokenResponse := Token{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeUserTokenCreate: Failed to decode json into struct: %+v", err)
	}

	if tokenResponse.Login != "" {
		// the ID consists of the login_name and the token name (foo/bar)
		d.SetId(fmt.Sprintf("%s/%s", d.Get("login_name").(string), d.Get("name").(string)))
		// we set the token value here as the API wont return it later
		if tokenResponse.Token != "" {
			d.Set("token", tokenResponse.Token)
		} else {
			return fmt.Errorf("resourceSonarqubeUserTokenCreate: Create response didn't contain the token")
		}
	} else {
		return fmt.Errorf("resourceSonarqubeUserTokenCreate: Create response didn't contain the user login")
	}

	return resourceSonarqubeUserTokenRead(d, m)
}

func resourceSonarqubeUserTokenRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/user_tokens/search"

	// split d.Id into login_name and the token name (foo/bar)
	login := strings.Split(d.Id(), "/")
	sonarQubeURL.RawQuery = url.Values{
		"login": []string{login[0]},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeUserTokenRead",
	)
	if err != nil {
		return fmt.Errorf("error reading Sonarqube user tokens: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	getTokensResponse := GetTokens{}
	err = json.NewDecoder(resp.Body).Decode(&getTokensResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeUserTokenCreate: Failed to decode json into struct: %+v", err)
	}

	// Loop over all user token to see if the current token exists.
	if getTokensResponse.Tokens != nil {
		for _, value := range getTokensResponse.Tokens {
			if d.Get("name").(string) == value.Name {
				d.SetId(fmt.Sprintf("%s/%s", d.Get("login_name").(string), d.Get("name").(string)))
				d.Set("login_name", getTokensResponse.Login)
				d.Set("name", value.Name)
				d.Set("expiration_date", value.ExpirationDate)
				return nil
			}
		}
	}

	return fmt.Errorf("resourceSonarqubeUserTokenCreate: Failed to find user token: %+v", d.Id())
}

func resourceSonarqubeUserTokenDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/user_tokens/revoke"
	sonarQubeURL.RawQuery = url.Values{
		"login": []string{d.Get("login_name").(string)},
		"name":  []string{d.Get("name").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeUserTokenDelete",
	)
	if err != nil {
		return fmt.Errorf("error deleting Sonarqube user token: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}
