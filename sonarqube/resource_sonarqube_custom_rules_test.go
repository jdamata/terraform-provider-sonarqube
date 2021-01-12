package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_custom_rule", &resource.Sweeper{
		Name: "sonarqube_custom_rule",
		F:    testSweepSonarqubeRuleSweeper,
	})
}

func testSweepSonarqubeRuleSweeper(r string) error {
	return nil
}

func testAccSonarqubeRuleBasicConfig(rnd string, key string, markdownDescription string, name string, severity string, status string, ruleType string, templateKey string) string {
	return fmt.Sprintf(`
		resource "sonarqube_custom_rule" "%[1]s" {
			key                  = "%[2]s"
			markdown_description = "%[3]s"
			name                 = "%[4]s"
			severity             = "%[5]s"
			status               = "%[6]s"
			rule_type            = "%[7]s"
			template_key         = "%[8]s"
		}`, rnd, key, markdownDescription, name, severity, status, ruleType, templateKey)
}

func TestAccSonarqubeRuleBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_custom_rule." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeRuleBasicConfig(rnd, "testAccSonarqubeRuleKey", "example description", "testAccSonarqubeRuleName", "INFO", "BETA", "CODE_SMELL", "flex:XPath"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubeRuleKey"),
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeRuleName"),
					resource.TestCheckResourceAttr(name, "severity", "INFO"),
					resource.TestCheckResourceAttr(name, "status", "BETA"),
					resource.TestCheckResourceAttr(name, "rule_type", "CODE_SMELL"),
					resource.TestCheckResourceAttr(name, "template_key", "flex:XPath"),
				),
			},
			{
				Config: testAccSonarqubeRuleBasicConfig(rnd, "testAccSonarqubeRuleKey1", "example description 1", "testAccSonarqubeRuleName1", "MINOR", "DEPRECATED", "BUG", "flex:XPath"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubeRuleKey1"),
					resource.TestCheckResourceAttr(name, "markdownDescription", "example description 1"),
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeRuleName1"),
					resource.TestCheckResourceAttr(name, "severity", "MINOR"),
					resource.TestCheckResourceAttr(name, "status", "DEPRECATED"),
					resource.TestCheckResourceAttr(name, "rule_type", "BUG"),
					resource.TestCheckResourceAttr(name, "template_key", "flex:XPath"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubeRuleKey1"),
					resource.TestCheckResourceAttr(name, "markdownDescription", "example description 1"),
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeRuleName1"),
					resource.TestCheckResourceAttr(name, "severity", "MINOR"),
					resource.TestCheckResourceAttr(name, "status", "DEPRECATED"),
					resource.TestCheckResourceAttr(name, "rule_type", "BUG"),
					resource.TestCheckResourceAttr(name, "template_key", "flex:XPath"),
				),
			},
		},
	})
}
