package main

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		// Provider configuration
		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SONAR_USER", "SONARQUBE_USER"}, nil),
			},
			"pass": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SONAR_PASS", "SONARQUBE_PASS"}, nil),
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SONAR_URL", "SONARQUBE_URL"}, nil),
			},
		},
		// Add the resources supported by this provider to this map.
		ResourcesMap: map[string]*schema.Resource{
			"sonarqube_qualitygate": qualityGate(),
		},
		ConfigureFunc: providerConfigure,
	}
}

//ProviderConfiguration contains the formatted credentials for comminication with the sonarqube api
type ProviderConfiguration struct {
	httpClient *http.Client
	sonarURL   string
	sonarUser  string
	sonarPass  string
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := &http.Client{}
	user := d.Get("user").(string)
	pass := d.Get("pass").(string)
	url := d.Get("url").(string)

	return &ProviderConfiguration{
		httpClient: client,
		sonarURL:   url,
		sonarUser:  user,
		sonarPass:  pass,
	}, nil
}
