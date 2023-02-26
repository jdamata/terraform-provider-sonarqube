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

func testAccSonarqubeQualitygateBasicConfig(rnd string, name string, is_default string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualitygate" "%[1]s" {
			name = "%[2]s"
			is_default = "%[3]s"
		}`, rnd, name, is_default)
}

func testAccSonarqubeQualitygateCopyConfig(rnd string, base_name string, condition_name string, threshold string, op string, copy_name string) string {
	return fmt.Sprintf(`
	resource "sonarqube_qualitygate" "%[2]s" {
		name = "%[2]s"
	}
	
	resource "sonarqube_qualitygate_condition" "qualitygate_condition" {
		gatename  = sonarqube_qualitygate.%[2]s.id
		metric    = "%[3]s"
		threshold = "%[4]s"
		op        = "%[5]s"
	}

	resource "sonarqube_qualitygate" "main" {
		name = "%[6]s"
		copy_from = "%[2]s"
	}`, rnd, base_name, condition_name, threshold, op, copy_name)
}

func TestAccSonarqubeQualitygateBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygateBasicConfig(rnd, "testAccSonarqubeQualitygate", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeQualitygate"),
					resource.TestCheckResourceAttr(name, "is_default", "true"),
				),
			},
			{
				Config: testAccSonarqubeQualitygateCopyConfig(rnd, "baseGate", "test_condition", "68", "GT", "baseGateCopy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "baseGateCopy"),
					resource.TestCheckResourceAttr(name, "is_default", "false"),
					resource.TestCheckResourceAttr(name, "condition.6.metric", "test_condition"),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_default"},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeQualitygate"),
					resource.TestCheckResourceAttr(name, "is_default", "true"),
				),
			},
		},
	})
}
