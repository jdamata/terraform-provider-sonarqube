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

func testAccSonarqubeRuleBasicConfig(rnd string, key string, markdown_description string, name string) string {
	return fmt.Sprintf(`
		resource "sonarqube_rule" "%[1]s" {
			key = "%[2]s"
			markdown_description = "%[3]s"
			name = "%[4]s"
		}`, rnd, key, markdown_description, name)
}

func TestAccSonarqubeRuleBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_rule." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeRuleBasicConfig(rnd, "key", "markdown_description", "name"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "key"),
					resource.TestCheckResourceAttr(name, "markdown_description", "markdown_description"),
					resource.TestCheckResourceAttr(name, "name", "name"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "key"),
					resource.TestCheckResourceAttr(name, "markdown_description", "markdown_description"),
					resource.TestCheckResourceAttr(name, "params", "key1=v1"),
					resource.TestCheckResourceAttr(name, "prevent_reactivation", "true"),
					resource.TestCheckResourceAttr(name, "severity", "INFO"),
					resource.TestCheckResourceAttr(name, "status", "BETA"),
				),
			},
		},
	})
}