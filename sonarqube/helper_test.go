package sonarqube

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func generateRandomResourceName() string {
	return acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
}

func generateHCLList(s []string) string {
	semiformat := fmt.Sprintf("%+q", s)      // Turn the slice into a string that looks like ["one" "two" "three"]
	tokens := strings.Split(semiformat, " ") // Split this string by spaces
	return strings.Join(tokens, ", ")        // Join the Slice together (that was split by spaces) with commas
}

// Skip acceptance test if sonarqube version is less than 8
func skipFuncSonarVersion_7(m interface{}) func() (bool, error) {
	return func() (bool, error) {
		sonarQubeVersion := m.(*ProviderConfiguration).sonarQubeVersion
		if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
			return false, nil
		}
		return true, nil
	}
}

// Skip acceptance test if sonarqube version is greater than 8
func skipFuncSonarVersion_8(m interface{}) func() (bool, error) {
	return func() (bool, error) {
		sonarQubeVersion := m.(*ProviderConfiguration).sonarQubeVersion
		if version, _ := version.NewVersion("8.0"); sonarQubeVersion.LessThanOrEqual(version) {
			return true, nil
		}
		return false, nil
	}
}
