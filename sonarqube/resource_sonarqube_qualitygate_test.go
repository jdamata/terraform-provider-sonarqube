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

func testAccSonarqubeQualitygateCopyConfig(rnd string, baseName string, conditionName string, threshold string, op string, copyName string) string {
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

	resource "sonarqube_qualitygate" "%[1]s" {
		depends_on = [sonarqube_qualitygate.%[2]s, sonarqube_qualitygate_condition.qualitygate_condition]
		name = "%[6]s"
		copy_from = "%[2]s"
	}`, rnd, baseName, conditionName, threshold, op, copyName)
}

func TestAccSonarqubeQualitygateCopy(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate." + rnd

	// Copy test variables
	baseGateName := "baseGate"
	baseGateResourceName := "sonarqube_qualitygate." + baseGateName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygateCopyConfig(rnd, "baseGate", "comment_lines_density", "68", "LT", "baseGateCopy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(baseGateResourceName, "conditions", name, "conditions"),
				),
			},
		},
	})
}
