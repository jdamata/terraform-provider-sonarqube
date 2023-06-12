package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_qualitygate_project_association", &resource.Sweeper{
		Name: "sonarqube_qualitygate_project_association",
		F:    testSweepSonarqubeQualitygateProjectAssociationSweeper,
	})
}

// TODO: implement sweeper to clean up projects: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeQualitygateProjectAssociationSweeper(r string) error {
	return nil
}

func testAccSonarqubeQualitygateProjectAssociationGateName(rnd string, name string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualitygate" "%[1]s" {
			name = "%[2]s"

			condition {
				metric    = "new_coverage"
				op        = "LT"
				threshold = "30"
			}
		}

		resource "sonarqube_project" "%[1]s" {
			name       = "%[2]s"
			project    = "%[2]s"
			visibility = "public"
		}

		resource "sonarqube_qualitygate_project_association" "%[1]s" {
			gatename   = sonarqube_qualitygate.%[1]s.name
			projectkey = sonarqube_project.%[1]s.project
		}`, rnd, name)
}

func TestAccSonarqubeQualitygateProjectAssociationGateName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate_project_association." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygateProjectAssociationGateName(rnd, "testAccSonarqubeProjectAssociationGateName"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "gatename", "testAccSonarqubeProjectAssociationGateName"),
					resource.TestCheckResourceAttr(name, "projectkey", "testAccSonarqubeProjectAssociationGateName"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "gatename", "testAccSonarqubeProjectAssociationGateName"),
					resource.TestCheckResourceAttr(name, "projectkey", "testAccSonarqubeProjectAssociationGateName"),
				),
			},
		},
	})
}
