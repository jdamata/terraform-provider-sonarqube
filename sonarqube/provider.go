package sonarqube

import (
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
			"sonarqube_project":                         resourceSonarqubeProject(),
			"sonarqube_qualitygate":                     resourceSonarqubeQualityGate(),
			"sonarqube_qualitygate_condition":           resourceSonarqubeQualityGateCondition(),
			"sonarqube_qualitygate_project_association": resourceSonarqubeQualityGateProjectAssociation(),
		},
		ConfigureFunc: providerConfigure,
	}
	return sonarqubeProvider
}

//ProviderConfiguration contains the sonarqube providers configuration
type ProviderConfiguration struct {
	httpClient   *http.Client
	sonarQubeURL url.URL
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := &http.Client{}

	sonarQubeURL := url.URL{
		Scheme: d.Get("scheme").(string),
		Host:   d.Get("host").(string),
		User:   url.UserPassword(d.Get("user").(string), d.Get("pass").(string)),
	}

	return &ProviderConfiguration{
		httpClient:   client,
		sonarQubeURL: sonarQubeURL,
	}, nil
}
