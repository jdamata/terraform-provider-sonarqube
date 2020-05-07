package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/jdamata/terraform-provider-sonarqube/sonarqube"
)

func main() {
	plugin.Serve(
		&plugin.ServeOpts{
			ProviderFunc: sonarqube.Provider,
		},
	)
}
