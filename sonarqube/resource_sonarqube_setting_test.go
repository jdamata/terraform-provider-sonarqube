package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_setting", &resource.Sweeper{
		Name: "sonarqube_setting",
		F:    testSweepSonarqubeSettingSweeper,
	})
}

func testSweepSonarqubeSettingSweeper(r string) error {
	return nil
}

func testAccSonarqubeSettingBasicConfig(rnd string, key string, value string) string {
	return fmt.Sprintf(`
		resource "sonarqube_setting" "%[1]s" {
			key = "%[2]s"
			value = "%[3]s"
		}`, rnd, key, value)
}

func TestAccSonarqubeSettingBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_setting." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeSettingBasicConfig(rnd, "sonar.demo", "sonarqube@example.org"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "sonar.demo"),
					resource.TestCheckResourceAttr(name, "value", "sonarqube@example.org"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "sonar.demo"),
					resource.TestCheckResourceAttr(name, "value", "sonarqube@example.org"),
				),
			},
		},
	})
}
