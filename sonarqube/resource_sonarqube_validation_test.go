package sonarqube

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestValidationVisibilityRejectsInternal verifies that the ValidateFunc on
// sonarqube_project.visibility rejects values outside ["public", "private"].
func TestValidationVisibilityRejectsInternal(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "sonarqube_project" "test" {
  name       = "test"
  project    = "test"
  visibility = "internal"
}
`,
				ExpectError: regexp.MustCompile(`expected visibility to be one of`),
			},
		},
	})
}

// TestValidationInstalledVersionRejectsNonSemver verifies that the ValidateFunc
// on the provider's installed_version attribute rejects non-semver strings.
func TestValidationInstalledVersionRejectsNonSemver(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
provider "sonarqube" {
  host              = "http://localhost:9001"
  token             = "fake-token"
  installed_version = "not-a-version"
}
resource "sonarqube_project" "test" {
  name    = "test"
  project = "test"
}
`,
				ExpectError: regexp.MustCompile(`is not a valid version string`),
			},
		},
	})
}

// TestValidationProviderRejectsBothTokenAndPass verifies ExactlyOneOf on
// the provider's token/pass attributes: both set must be rejected.
func TestValidationProviderRejectsBothTokenAndPass(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
provider "sonarqube" {
  host  = "http://localhost:9001"
  token = "my-token"
  user  = "admin"
  pass  = "my-pass"
}
resource "sonarqube_project" "test" {
  name    = "test"
  project = "test"
}
`,
				ExpectError: regexp.MustCompile(`only one of .pass,token. can be specified`),
			},
		},
	})
}

// TestValidationProviderRejectsNeitherTokenNorPass verifies ExactlyOneOf on
// the provider's token/pass attributes: neither set must be rejected.
// Env vars are blanked so DefaultFunc cannot supply a fallback.
func TestValidationProviderRejectsNeitherTokenNorPass(t *testing.T) {
	// Blank all auth env vars so DefaultFunc returns nil for both token and pass.
	for _, v := range []string{
		"SONAR_TOKEN", "SONARQUBE_TOKEN",
		"SONAR_USER", "SONARQUBE_USER",
		"SONAR_PASS", "SONARQUBE_PASS",
	} {
		t.Setenv(v, "")
	}

	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
provider "sonarqube" {
  host = "http://localhost:9001"
}
resource "sonarqube_project" "test" {
  name    = "test"
  project = "test"
}
`,
				ExpectError: regexp.MustCompile(`one of .pass,token. must be specified`),
			},
		},
	})
}
