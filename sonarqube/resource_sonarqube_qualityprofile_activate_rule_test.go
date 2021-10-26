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

func testAccSonarqubeQualityprofileActivateRuleBasicConfig(rnd string, name string, language string, severity string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualityprofile" "%[1]s" {
			name     = "%[2]s"
			language = "%[3]s"
		}

		resource "sonarqube_rule" "%[1]s" {
			custom_key = "%[2]s"
			markdown_description = "My rule"
			name = "%[2]s" 
			severity = "%[4]s"
			template_key = "xml:XPathCheck"
			type = "VULNERABILITY"
		}

		resource "sonarqube_qualityprofile_activate_rule" "%[1]s" {
			key = sonarqube_qualityprofile.%[1]s.key
			rule = sonarqube_rule.%[1]s.id
			severity = "%[4]s"
		}`, rnd, name, language, severity)
}

func TestAccSonarqubeQualityprofileActivateRuleBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_quality_profile_activate_rule." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{

			{
				Config: testAccSonarqubeQualityprofileActivateRuleBasicConfig(rnd, "testProfile", "xml", "BLOCKER"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testProfile"),
					resource.TestCheckResourceAttr(name, "language", "xml"),
					resource.TestCheckResourceAttr(name, "custom_key", "testProfile"),
					resource.TestCheckResourceAttr(name, "severity", "BLOCKER"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testProfile"),
					resource.TestCheckResourceAttr(name, "language", "xml"),
					resource.TestCheckResourceAttr(name, "custom_key", "testProfile"),
					resource.TestCheckResourceAttr(name, "severity", "BLOCKER"),
				),
			},
		},
	})
}
