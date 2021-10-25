package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_qualityprofile_activate_rule", &resource.Sweeper{
		Name: "sonarqube_qualityprofile_activate_rule",
		F:    testSweeepSonarqibeQualityprofileActivateRuleSweeper,
	})
}

func testSweeepSonarqibeQualityprofileActivateRuleSweeper(r string) error {
	return nil
}

func testAccSonarqubeQualityprofileActivateRuleBasicConfig(rnd string, name string, description string, language string, severity string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualityprofile" "%[1]s" {
			name     = "%[2]s"
			language = "%[4]s"
		}

		resource "sonarqube_rule" "%[1]s" {
			key = "%[2]s"
			markdown_description = "%[3]s"
			name = "%[2]s"
		}

		resource "sonarqube_qualityprofile_activate_rule" "%[1]s" {
			key = "sonarqube_qualityprofile.%[1]s.key"
			rule = "sonarqube_rule.%[1]s.id"
			severity = "%[5]s"
		}`, rnd, name, description, language, severity)
}

func TestAccSonarqubeQualityprofileActivateRuleBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_quality_profile_activate_rule" + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualityprofileActivateRuleBasicConfig(rnd, "name", "description", "java", "BLOCKER"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "name"),
					resource.TestCheckResourceAttr(name, "language", "java"),
					resource.TestCheckResourceAttr(name, "key", "name"),
					resource.TestCheckResourceAttr(name, "markdown_description", "description"),
					resource.TestCheckResourceAttr(name, "severity", "BLOCKER"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "name"),
					resource.TestCheckResourceAttr(name, "language", "java"),
					resource.TestCheckResourceAttr(name, "key", "false"),
					resource.TestCheckResourceAttr(name, "markdown_description", "description"),
					resource.TestCheckResourceAttr(name, "severity", "BLOCKER"),
				),
			},
		},
	})
}
