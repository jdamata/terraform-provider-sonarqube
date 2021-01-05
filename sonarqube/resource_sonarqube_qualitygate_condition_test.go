package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_qualitygate_condition", &resource.Sweeper{
		Name: "sonarqube_qualitygate_condition",
		F:    testSweepSonarqubeQualitygateConditionSweeper,
	})
}

// TODO: implement sweeper to clean up projects: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeQualitygateConditionSweeper(r string) error {
	return nil
}

func testAccSonarqubeQualitygateConditionBasicConfig(id string, name string, metric string, threshold string, op string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualitygate" "%[1]s" {
			name = "%[2]s"
		}

		resource "sonarqube_qualitygate_condition" "%[1]s" {
			gatename  = sonarqube_qualitygate.%[1]s.name
			metric    = "%[3]s"
			threshold = "%[4]s"
			op        = "%[5]s"
		}		
		`, id, name, metric, threshold, op)
}

func TestAccSonarqubeQualitygateConditionBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate_condition." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygateConditionBasicConfig(rnd, "testAccSonarqubeQualitygateCondition", "vulnerabilities", "10", "GT"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "metric", "vulnerabilities"),
					resource.TestCheckResourceAttr(name, "threshold", "10"),
					resource.TestCheckResourceAttr(name, "op", "GT"),
				),
			},
		},
	})
}
