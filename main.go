package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/jdamata/terraform-provider-sonarqube/sonarqube"
)

func main() {

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(
		&plugin.ServeOpts{
			Debug:        debug,
			ProviderAddr: "registry.terraform.io/jdamata/sonarqube",
			ProviderFunc: sonarqube.Provider,
		},
	)
}
