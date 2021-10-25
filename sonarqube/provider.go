package sonarqube

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-cleanhttp"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var sonarqubeProvider *schema.Provider

// Provider for sonarqube
func Provider() *schema.Provider {
	sonarqubeProvider = &schema.Provider{
		// Provider configuration
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SONAR_USER", "SONARQUBE_USER"}, nil),
				Required:    true,
			},
			"pass": {
				Type:        schema.TypeString,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SONAR_PASS", "SONARQUBE_PASS"}, nil),
				Required:    true,
			},
			"host": {
				Type:        schema.TypeString,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SONAR_HOST", "SONARQUBE_HOST"}, nil),
				Required:    true,
			},
			"installed_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_insecure_skip_verify": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Allows ignoring insecure certificates when set to true. Defaults to false. Disabling TLS verification is dangerous and should only be done for local testing.",
				Default:     false,
			},
		},
		// Add the resources supported by this provider to this map.
		ResourcesMap: map[string]*schema.Resource{
			"sonarqube_group":                              resourceSonarqubeGroup(),
			"sonarqube_permission_template":                resourceSonarqubePermissionTemplate(),
			"sonarqube_permissions":                        resourceSonarqubePermissions(),
			"sonarqube_plugin":                             resourceSonarqubePlugin(),
			"sonarqube_project":                            resourceSonarqubeProject(),
			"sonarqube_qualityprofile":                     resourceSonarqubeQualityProfile(),
			"sonarqube_qualityprofile_project_association": resourceSonarqubeQualityProfileProjectAssociation(),
			"sonarqube_qualitygate":                        resourceSonarqubeQualityGate(),
			"sonarqube_qualitygate_condition":              resourceSonarqubeQualityGateCondition(),
			"sonarqube_qualitygate_project_association":    resourceSonarqubeQualityGateProjectAssociation(),
			"sonarqube_user":                               resourceSonarqubeUser(),
			"sonarqube_user_token":                         resourceSonarqubeUserToken(),
			"sonarqube_webhook":                            resourceSonarqubeWebhook(),
			"sonarqube_rule":                               resourceSonarqubeRule(),
			"sonarqube_qualityprofile_activate_rule":       resourceSonarqubeQualityProfileRule(),
		},
		ConfigureFunc: configureProvider,
	}
	return sonarqubeProvider
}

//ProviderConfiguration contains the sonarqube providers configuration
type ProviderConfiguration struct {
	httpClient       *retryablehttp.Client
	sonarQubeURL     url.URL
	sonarQubeVersion *version.Version
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: d.Get("tls_insecure_skip_verify").(bool),
	}

	client := retryablehttp.NewClient()
	client.HTTPClient.Transport = transport

	host, err := url.Parse(d.Get("host").(string))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse sonarqube host: %+v", err)
	}

	sonarQubeURL := url.URL{
		Scheme:     host.Scheme,
		Host:       host.Host,
		User:       url.UserPassword(d.Get("user").(string), d.Get("pass").(string)),
		ForceQuery: true,
	}

	var installedVersion *version.Version
	if v, ok := d.GetOk("installed_version"); ok {
		installedVersion, err = version.NewVersion(v.(string))
		if err != nil {
			return nil, err
		}
	} else {
		// Check that the sonarqube api is available, get version
		installedVersion, err = sonarqubeHealth(client, sonarQubeURL)
		if err != nil {
			return nil, err
		}
	}

	minimumVersion, _ := version.NewVersion("7.9")
	if installedVersion.LessThan(minimumVersion) {
		return nil, fmt.Errorf("Unsupported version of sonarqube. Minimum supported version is %+v. Running version is %+v", minimumVersion, installedVersion)
	}

	return &ProviderConfiguration{
		httpClient:       client,
		sonarQubeURL:     sonarQubeURL,
		sonarQubeVersion: installedVersion,
	}, nil
}

func sonarqubeHealth(client *retryablehttp.Client, sonarqube url.URL) (*version.Version, error) {
	// Make request to sonarqube version endpoint
	sonarqube.Path = "api/server/version"
	req, err := retryablehttp.NewRequest("GET", sonarqube.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("Unable to construct sonarqube version request: %+v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to reach sonarqube: %+v", err)
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Sonarqube version api did not return a 200: %+v", err)
	}

	// Read in the response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse response body on GET sonarqube version api: %+v", err)
	}

	// Convert response to a int.
	bodyString := string(bodyBytes)
	installedVersion, err := version.NewVersion(bodyString)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert sonarqube version to a version: %+v", err)
	}

	return installedVersion, nil
}
