package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_rule", &resource.Sweeper{
		Name: "sonarqube_rule",
		F:    testSweepSonarqubeRuleSweeper,
	})
}

func testSweepSonarqubeRuleSweeper(r string) error {
	return nil
}

func testAccSonarqubeRuleBasicConfig(rnd string, custom_key string, markdown_description string, name string, template_key string, severity string, status string, type_p string) string {
	return fmt.Sprintf(`
		resource "sonarqube_rule" "%[1]s" {
			custom_key = "%[2]s"
			markdown_description = "%[3]s"
			name = "%[4]s"
			template_key = "%[5]s"
			severity = "%[6]s"
			status = "%[7]s"
			type = "%[8]s"
		}`, rnd, custom_key, markdown_description, name, template_key, severity, status, type_p)
}

func TestAccSonarqubeRuleBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_rule." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeRuleBasicConfig(rnd, "key", "markdown_description", "name", "xml:XPathCheck", "INFO", "BETA", "VULNERABILITY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "custom_key", "key"),
					resource.TestCheckResourceAttr(name, "markdown_description", "markdown_description"),
					resource.TestCheckResourceAttr(name, "name", "name"),
					resource.TestCheckResourceAttr(name, "template_key", "xml:XPathCheck"),
					resource.TestCheckResourceAttr(name, "severity", "INFO"),
					resource.TestCheckResourceAttr(name, "status", "BETA"),
					resource.TestCheckResourceAttr(name, "type", "VULNERABILITY"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "custom_key", "key"),
					resource.TestCheckResourceAttr(name, "markdown_description", "markdown_description"),
					resource.TestCheckResourceAttr(name, "name", "name"),
					resource.TestCheckResourceAttr(name, "template_key", "xml:XPathCheck"),
					resource.TestCheckResourceAttr(name, "severity", "INFO"),
					resource.TestCheckResourceAttr(name, "status", "BETA"),
					resource.TestCheckResourceAttr(name, "type", "VULNERABILITY"),
				),
			},
		},
	})
}
