package sonarqube

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// GetTokens struct
type GetTokens struct {
	Login  string  `json:"login,omitempty"`
	Tokens []Token `json:"userTokens,omitempty"`
}

// Token struct
type Token struct {
	Login          string       `json:"login,omitempty"`
	Name           string       `json:"name,omitempty"`
	Token          string       `json:"token,omitempty"`
	ExpirationDate string       `json:"expirationDate,omitempty"`
	Type           string       `json:"type,omitempty"`
	CreatedAt      string       `json:"createdAt,omitempty"`
	IsExpired      bool         `json:"isExpired,omitempty"`
	Project        TokenProject `json:"project,omitempty"`
}

type TokenProject struct {
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}

// Token types
type TokenType string

const (
	UserToken            TokenType = "USER_TOKEN"
	GlobalAnalysisToken  TokenType = "GLOBAL_ANALYSIS_TOKEN"
	ProjectAnalysisToken TokenType = "PROJECT_ANALYSIS_TOKEN"
)

// Returns the resource represented by this file.
func resourceSonarqubeUserToken() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Sonarqube User token resource. This can be used to manage Sonarqube User tokens.",
		Create:      resourceSonarqubeUserTokenCreate,
		Read:        resourceSonarqubeUserTokenRead,
		Update:      resourceSonarqubeUserTokenRead,
		Delete:      resourceSonarqubeUserTokenDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeUserTokenImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 100)),
				Description:      "The name of the Token to create. Changing this forces a new resource to be created.",
			},
			"login_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The login name of the User for which the token should be created. If not set, the token is created for the authenticated user. Changing this forces a new resource to be created.",
			},
			"expiration_date": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The expiration date of the token being generated, in ISO 8601 format (YYYY-MM-DD). If not set, default to no expiration.",
			},
			"metadata_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Manage only the metadata and revocation of an already issued token. Creation fails if the named token does not exist or its type or project does not match. The token value is never stored in state in this mode. Changing this setting does not recreate the token, and disabling it cannot recover an existing token value.",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The token value.",
			},
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          UserToken,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(UserToken), string(GlobalAnalysisToken), string(ProjectAnalysisToken)}, false)),
				Description:      "The kind of Token to create. Changing this forces a new resource to be created. Possible values are USER_TOKEN, GLOBAL_ANALYSIS_TOKEN, or PROJECT_ANALYSIS_TOKEN. Defaults to USER_TOKEN. If set to PROJECT_ANALYSIS_TOKEN, then the project_key must also be specified.",
			},
			"project_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The key of the only project that can be analyzed by the PROJECT_ANALYSIS TOKEN being created. Changing this forces a new resource to be created.",
			},
		},
	}
}

func resourceSonarqubeUserTokenCreate(d *schema.ResourceData, m interface{}) error {
	if d.Get("metadata_only").(bool) {
		tokens, err := readSonarqubeUserTokens(m, d.Get("login_name").(string))
		if err != nil {
			return err
		}
		for _, token := range tokens.Tokens {
			if token.Name != d.Get("name").(string) {
				continue
			}
			if err := validateTokenMetadata(d, token); err != nil {
				return err
			}
			d.SetId(fmt.Sprintf("%s/%s", tokens.Login, token.Name))
			return resourceSonarqubeUserTokenRead(d, m)
		}
		return fmt.Errorf("token %q must already exist when metadata_only is enabled", d.Get("name").(string))
	}
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/user_tokens/generate"

	tokenType := TokenType(d.Get("type").(string))
	rawQuery := url.Values{
		"name": []string{d.Get("name").(string)},
		"type": []string{string(tokenType)},
	}

	if tokenType == UserToken {
		loginName := d.Get("login_name").(string)
		if loginName != "" {
			rawQuery.Add("login", loginName)
		}
	} else if tokenType == ProjectAnalysisToken {
		projectKey := d.Get("project_key").(string)
		if projectKey == "" {
			return fmt.Errorf("resourceSonarqubeUserTokenCreate: 'project_key' must be configured when the token 'type' is %s", ProjectAnalysisToken)
		}
		rawQuery.Add("projectKey", projectKey)
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
			if err := d.Set("token", tokenResponse.Token); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("resourceSonarqubeUserTokenCreate: Create response didn't contain the token")
		}
	} else {
		return fmt.Errorf("resourceSonarqubeUserTokenCreate: Create response didn't contain the user login")
	}

	return resourceSonarqubeUserTokenRead(d, m)
}

func resourceSonarqubeUserTokenRead(d *schema.ResourceData, m interface{}) error {
	login_name := strings.Split(d.Id(), "/")
	getTokensResponse, err := readSonarqubeUserTokens(m, login_name[0])
	if err != nil {
		return err
	}
	if d.Get("metadata_only").(bool) {
		if err := d.Set("token", nil); err != nil {
			return err
		}
	}

	// Loop over all user token to see if the current token exists.
	if getTokensResponse.Tokens != nil {
		for _, value := range getTokensResponse.Tokens {
			if login_name[1] == value.Name {
				if d.Get("metadata_only").(bool) {
					if err := validateTokenMetadata(d, value); err != nil {
						return err
					}
				}
				d.SetId(fmt.Sprintf("%s/%s", getTokensResponse.Login, value.Name))
				errs := []error{}
				// Set login_name only if it's a token for another user (not the authenticated one),
				// or if this is an import and only the ID is provided (i.e., name is not set).
				if d.Get("login_name").(string) != "" || d.Get("name").(string) == "" {
					errs = append(errs, d.Set("login_name", getTokensResponse.Login))
				}
				if d.Get("type").(string) != value.Type {
					errs = append(errs, d.Set("type", value.Type))
				}
				errs = append(errs, d.Set("name", value.Name))
				if value.ExpirationDate != "" {
					dateReceived, errTimeParse := time.Parse("2006-01-02T15:04:05-0700", value.ExpirationDate)
					if errTimeParse != nil {
						return fmt.Errorf("resourceSonarqubeUserTokenRead: Failed to parse ExpirationDate: %+v", errTimeParse)
					}
					errs = append(errs, d.Set("expiration_date", dateReceived.Format("2006-01-02")))
				}
				return errors.Join(errs...)
			}
		}
	}

	// Token was revoked outside Terraform. Clear the ID so the next plan recreates it
	// instead of failing on refresh.
	d.SetId("")
	return nil
}

func resourceSonarqubeUserTokenDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/user_tokens/revoke"
	rawQuery := url.Values{
		"name": []string{d.Get("name").(string)},
	}
	login := d.Get("login_name").(string)
	if login != "" {
		rawQuery.Add("login", login)
	}

	sonarQubeURL.RawQuery = rawQuery.Encode()

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

func resourceSonarqubeUserTokenImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeUserTokenRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func readSonarqubeUserTokens(m interface{}, login string) (*GetTokens, error) {
	conf := m.(*ProviderConfiguration)
	endpoint := conf.sonarQubeURL
	endpoint.Path = strings.TrimSuffix(endpoint.Path, "/") + "/api/user_tokens/search"
	if login != "" {
		endpoint.RawQuery = url.Values{"login": []string{login}}.Encode()
	}
	response, err := httpRequestHelper(conf.httpClient, "GET", endpoint.String(), http.StatusOK, "resourceSonarqubeUserTokenRead")
	if err != nil {
		return nil, fmt.Errorf("error reading Sonarqube user tokens: %w", err)
	}
	defer response.Body.Close()
	var tokens GetTokens
	if err := json.NewDecoder(response.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("resourceSonarqubeUserTokenRead: Failed to decode json into struct: %w", err)
	}
	return &tokens, nil
}

func validateTokenMetadata(d *schema.ResourceData, token Token) error {
	if expected, ok := d.GetOk("expiration_date"); ok {
		expiration, err := time.Parse("2006-01-02T15:04:05-0700", token.ExpirationDate)
		if err != nil || expiration.Format("2006-01-02") != expected.(string) {
			return fmt.Errorf("token %q does not have the expected expiration date", token.Name)
		}
	}
	if token.Type != d.Get("type").(string) {
		return fmt.Errorf("token %q has type %q, expected %q", token.Name, token.Type, d.Get("type"))
	}
	if token.Type == string(ProjectAnalysisToken) && token.Project.Key != d.Get("project_key").(string) {
		return fmt.Errorf("token %q belongs to project %q, expected %q", token.Name, token.Project.Key, d.Get("project_key"))
	}
	return nil
}
