package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_qualitygate", &resource.Sweeper{
		Name: "sonarqube_qualitygate",
		F:    testSweepSonarqubeQualitygateSweeper,
	})
}

func testSweepSonarqubeQualitygateSweeper(r string) error {
	return nil
}

func testAccSonarqubeQualitygateBasicConfig(rnd string, name string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualitygate" "%[1]s" {
			name = "%[2]s"
		}`, rnd, name)
}

func TestAccSonarqubeQualitygateBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygateBasicConfig(rnd, "testAccSonarqubeQualitygate"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeQualitygate"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeQualitygate"),
				),
			},
		},
	})
}
