package sonarqube

import (
	"fmt"
	"strings"
)

// Community edition has no DevOps platform bindings unless a plugin such as the Community Branch
// Plugin adds the endpoints, so the provider calls them and only explains a failure.
func bindingNotSupportedHint(conf *ProviderConfiguration, platform string, err error) error {
	if strings.ToLower(conf.sonarQubeEdition) != "community" {
		return err
	}
	return fmt.Errorf("%s bindings need the Developer edition or above, or a plugin that provides them such as the Community Branch Plugin (running SonarQube %s %s): %w", platform, conf.sonarQubeEdition, conf.sonarQubeVersion, err)
}
