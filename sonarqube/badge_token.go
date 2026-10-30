package sonarqube

import "os"

// Badge sensitivity is a schema feature, so it must be selected before provider
// configuration and kept consistent across all Terraform commands and aliases.
func badgeTokenSensitive() bool {
	return os.Getenv("SONARQUBE_BADGE_TOKEN_SENSITIVE") == "true"
}
