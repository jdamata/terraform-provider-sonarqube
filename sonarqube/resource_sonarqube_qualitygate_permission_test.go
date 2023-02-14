package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_group_qualitygate_permission", &resource.Sweeper{
		Name: "sonarqube_group_qualitygate_permission",
		F:    testSweepSonarqubeQualitygatePermissionSweeper,
	})
}

func testSweepSonarqubeQualitygatePermissionSweeper(r string) error {
	return nil
}

func testAccPreCheckQualityGatePermissionFeature(t *testing.T) {
	sonarQubeVersion := testAccProvider.Meta().(*ProviderConfiguration).sonarQubeVersion

	minimumVersion, _ := version.NewVersion("9.2")
	if sonarQubeVersion.LessThan(minimumVersion) {
		t.Skipf("Skipping test of unsupported feature")
	}
}

func testAccSonarqubeQualitygatePermissionUserConfig(rnd string, gateName string, loginName string) string {
	return fmt.Sprintf(`
		resource "sonarqube_user" "%[1]s_user" {
			login_name = "%[3]s"
			name       = "Test User"
			email      = "terraform-test@sonarqube.com"
			password   = "secret-sauce!"
		}

		resource "sonarqube_qualitygate" "%[1]s_gate" {
			name = "%[2]s"
		}

		resource "sonarqube_qualitygate_permission" "%[1]s" {
			gatename = sonarqube_qualitygate.%[1]s_gate.name
			type     = "user"
			subject  = sonarqube_user.%[1]s_user.login_name
		}
		`, rnd, gateName, loginName)
}

func testAccSonarqubeQualitygatePermissionGroupConfig(rnd string, gateName string, groupName string) string {
	return fmt.Sprintf(`
		resource "sonarqube_group" "%[1]s_group" {
			name = "%[3]s"
		}

		resource "sonarqube_qualitygate" "%[1]s_gate" {
			name = "%[2]s"
		}

		resource "sonarqube_qualitygate_permission" "%[1]s" {
			gatename = sonarqube_qualitygate.%[1]s_gate.name
			type     = "group"
			subject  = sonarqube_group.%[1]s_group.name
		}
		`, rnd, gateName, groupName)
}

func TestAccSonarqubeQualitygatePermissionBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate_permission." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckQualityGatePermissionFeature(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygatePermissionUserConfig(rnd, "testAccSonarqubeQualtiyGate", "testAccSonarqubeUser"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "gatename", "testAccSonarqubeQualtiyGate"),
					resource.TestCheckResourceAttr(name, "type", "user"),
					resource.TestCheckResourceAttr(name, "subject", "testAccSonarqubeUser"),
				),
			},
			{
				Config: testAccSonarqubeQualitygatePermissionGroupConfig(rnd, "testAccSonarqubeQualtiyGate", "testAccSonarqubeGroup"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "gatename", "testAccSonarqubeQualtiyGate"),
					resource.TestCheckResourceAttr(name, "type", "group"),
					resource.TestCheckResourceAttr(name, "subject", "testAccSonarqubeGroup"),
				),
			},
		},
	})
}
