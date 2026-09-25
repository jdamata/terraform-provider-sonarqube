package sonarqube

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccSonarqubeQualityprofileActiveRulesDataSourceConfig(rnd string, name string, key string, severity string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualityprofile" "%[1]s" {
			name     = "%[2]s"
			language = "xml"
		}

		resource "sonarqube_rule" "%[1]s" {
			custom_key = "%[3]s"
			markdown_description = "My rule"
			name = "%[3]s"
			severity = "%[4]s"
			template_key = "xml:XPathCheck"
			type = "VULNERABILITY"
		}

		resource "sonarqube_qualityprofile_activate_rule" "%[1]s" {
			key = sonarqube_qualityprofile.%[1]s.key
			rule = sonarqube_rule.%[1]s.id
			severity = "%[4]s"
		}

		data "sonarqube_qualityprofile_active_rules" "%[1]s" {
			key = sonarqube_qualityprofile.%[1]s.key
			depends_on = [sonarqube_qualityprofile_activate_rule.%[1]s]
		}`, rnd, name, key, severity)
}

func TestAccSonarqubeQualityprofileActiveRulesDataSource(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "data.sonarqube_qualityprofile_active_rules." + rnd
	expectedRuleKey := "xml:activeRules"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualityprofileActiveRulesDataSourceConfig(rnd, "testProfileActiveRules", "activeRules", "BLOCKER"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRulePresentInList(name, expectedRuleKey),
				),
			},
		},
	})
}

func testAccCheckRulePresentInList(resourceName string, expectedKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		prefix := "rules."
		for k, v := range rs.Primary.Attributes {
			if strings.HasPrefix(k, prefix) && strings.HasSuffix(k, ".key") && v == expectedKey {
				return nil
			}
		}

		return fmt.Errorf("rule key %s not found in %s", expectedKey, resourceName)
	}
}
