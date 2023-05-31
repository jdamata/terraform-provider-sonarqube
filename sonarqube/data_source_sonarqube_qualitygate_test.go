package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccSonarqubeQualityGateDataSourceConfig(rnd string, name string, copy_from string, metric string, op string, threshold string) string {
	if len(copy_from) > 0 {
		return fmt.Sprintf(`
		resource "sonarqube_qualitygate" "%[1]s" {
			name      = "%[2]s"
			copy_from = "%[3]s"
		}
	
		data "sonarqube_qualitygate" "%[1]s" {
			name = sonarqube_qualitygate.%[1]s.name
		}`, rnd, name, copy_from)

	} else {

		return fmt.Sprintf(`
		resource "sonarqube_qualitygate" "%[1]s" {
			name      = "%[2]s"
			condition {
				metric    = "%[3]s"
				op        = "%[4]s"
				threshold = "%[5]s"
			}
		}
	
		data "sonarqube_qualitygate" "%[1]s" {
			name = sonarqube_qualitygate.%[1]s.name
		}`, rnd, name, metric, op, threshold)
	}
}

func TestAccSonarqubeQualityGateDataSource(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "data.sonarqube_qualitygate." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// QualityGate with copy_from
			{
				Config: testAccSonarqubeQualityGateDataSourceConfig(rnd, "testAccSonarqubeQualityGateDataSourceCopy", "Sonar way", "", "", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeQualityGateDataSourceCopy"),
				),
			},
			// QualityGate with condition
			{
				Config: testAccSonarqubeQualityGateDataSourceConfig(rnd, "testAccSonarqubeQualityGateDataSourceCondition", "", "new_coverage", "LT", "50"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeQualityGateDataSourceCondition"),
				),
			},
		},
	})
}
