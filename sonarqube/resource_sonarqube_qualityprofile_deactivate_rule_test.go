package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_qualityprofile_deactivate_rule", &resource.Sweeper{
		Name: "sonarqube_qualityprofile_deactivate_rule",
		F:    testSweeepSonarqibeQualityprofileDeactivateRuleSweeper,
	})
}

func testSweeepSonarqibeQualityprofileDeactivateRuleSweeper(r string) error {
	return nil
}

func testAccSonarqubeQualityprofileDeactivateRuleBasicConfig(rnd string, name string, key string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualityprofile" "%[1]s" {
			name     = "%[2]s"
			language = "xml"
		}

		data "sonarqube_rule" "%[1]s" {
		  key = "squid:forbidSonar"
		}

		resource "sonarqube_qualityprofile_deactivate_rule" "%[1]s" {
			key = sonarqube_qualityprofile.%[1]s.key
			rule = data.sonarqube_rule.%[1]s.key
		}`, rnd, name, key)
}

func TestAccSonarqubeQualityprofileDeactivateRuleBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualityprofile_deactivate_rule." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualityprofileDeactivateRuleBasicConfig(rnd, "testProfile", "xml:S3420"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "key"),
					resource.TestCheckResourceAttrSet(name, "rule"),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key", "rule"},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "key"),
					resource.TestCheckResourceAttrSet(name, "rule"),
				),
			},
		},
	})
}
