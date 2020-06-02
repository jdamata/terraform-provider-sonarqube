package sonarqube

import (
	"encoding/binary"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	log "github.com/sirupsen/logrus"
)

var sonarqubeProvider *schema.Provider

// Provider for sonarqube
func Provider() terraform.ResourceProvider {
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
			"scheme": {
				Type:        schema.TypeString,
				Default:     "https",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SONAR_SCHEME", "SONARQUBE_SCHEME"}, nil),
				Optional:    true,
			},
		},
		// Add the resources supported by this provider to this map.
		ResourcesMap: map[string]*schema.Resource{
			"sonarqube_group":                           resourceSonarqubeGroup(),
			"sonarqube_permissions":                     resourceSonarqubePermissions(),
			"sonarqube_plugin":                          resourceSonarqubePlugin(),
			"sonarqube_project":                         resourceSonarqubeProject(),
			"sonarqube_qualitygate":                     resourceSonarqubeQualityGate(),
			"sonarqube_qualitygate_condition":           resourceSonarqubeQualityGateCondition(),
			"sonarqube_qualitygate_project_association": resourceSonarqubeQualityGateProjectAssociation(),
		},
		ConfigureFunc: configureProvider,
	}
	return sonarqubeProvider
}

//ProviderConfiguration contains the sonarqube providers configuration
type ProviderConfiguration struct {
	httpClient   *http.Client
	sonarQubeURL url.URL
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	client := &http.Client{}

	sonarQubeURL := url.URL{
		Scheme:     d.Get("scheme").(string),
		Host:       d.Get("host").(string),
		User:       url.UserPassword(d.Get("user").(string), d.Get("pass").(string)),
		ForceQuery: true,
	}

	// Check that the sonarqube api is available and a supported version
	err := sonarqubeHealth(*client, sonarQubeURL)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &ProviderConfiguration{
		httpClient:   client,
		sonarQubeURL: sonarQubeURL,
	}, nil
}

func sonarqubeHealth(client http.Client, sonarqube url.URL) error {
	// Make request to sonarqube version endpoint
	sonarqube.Path = "api/server/version"
	req, err := http.NewRequest("GET", sonarqube.String(), http.NoBody)
	if err != nil {
		log.Error(err)
		return errors.New("Unable to construct sonarqube version request")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return errors.New("Unable to reach sonarqube")
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		return errors.New("Sonarqube version api did not return a 200")
	}

	// Read in the response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return errors.New("Failed to parse response body on GET sonarqube version api")
	}

	// Convert response to a int.
	version := binary.BigEndian.Uint64(bodyBytes)
	if version < 8 {
		log.Error(err)
		return errors.New("Unsupported version of sonarqube. Minimum supported version is 8")
	}

	return nil
}
