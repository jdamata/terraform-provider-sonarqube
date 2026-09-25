package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccSonarqubeSecuredSettingConfig(rnd, key, value string) string {
	return fmt.Sprintf(`
resource "sonarqube_secured_setting" "%[1]s" {
  key   = "%[2]s"
  value = "%[3]s"
}
`, rnd, key, value)
}

func testAccSonarqubeSecuredSettingProjectConfig(rnd, projectKey, key, value string) string {
	return fmt.Sprintf(`
resource "sonarqube_project" "%[1]s" {
  name       = "%[2]s"
  project    = "%[2]s"
  visibility = "public"
}

resource "sonarqube_secured_setting" "%[1]s" {
  key       = "%[3]s"
  value     = "%[4]s"
  component = sonarqube_project.%[1]s.project
}
`, rnd, projectKey, key, value)
}

func TestAccSonarqubeSecuredSettingBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	resourceName := "sonarqube_secured_setting." + rnd
	key := "sonar.auth.github.clientSecret.secured"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeSecuredSettingConfig(rnd, key, "initial-secret-value"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key", key),
					// value is write-only — not checked here
				),
			},
			// Update: value change is applied silently (re-posts to API)
			{
				Config: testAccSonarqubeSecuredSettingConfig(rnd, key, "updated-secret-value"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key", key),
				),
			},
			// Import: value is excluded because the API never returns it
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
		},
	})
}

func TestAccSonarqubeSecuredSettingProjectScoped(t *testing.T) {
	rnd := generateRandomResourceName()
	resourceName := "sonarqube_secured_setting." + rnd
	projectKey := "testAccSecuredSettingProject" + rnd
	// Use an arbitrary .secured-suffix key — built-in ones (sonar.auth.*.secured) are
	// global-only and cannot be set at project level.
	key := "test.acc.custom.password.secured"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeSecuredSettingProjectConfig(rnd, projectKey, key, "project-secret"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "component", projectKey),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateIdFunc:       testAccSonarqubeSecuredSettingImportID(resourceName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
		},
	})
}

func testAccSonarqubeSecuredSettingImportID(resourceNode string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceNode]
		if !ok {
			return "", fmt.Errorf("Resource node not found: %s", resourceNode)
		}
		key := rs.Primary.Attributes["key"]
		component := rs.Primary.Attributes["component"]
		if component != "" {
			return key + ":" + component, nil
		}
		return key, nil
	}
}
