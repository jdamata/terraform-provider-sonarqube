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

/*func testAccSonarqubeQualitygateProjectAssociationGateID(rnd string, name string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualitygate" "%[1]s" {
			name = "%[2]s"
		}

		resource "sonarqube_project" "%[1]s" {
			name       = "%[2]s"
			project    = "%[2]s"
			visibility = "public"
		}

		resource "sonarqube_qualitygate_project_association" "%[1]s" {
			gateid     = sonarqube_qualitygate.%[1]s.id
			projectkey = sonarqube_project.%[1]s.project
		}`, rnd, name)
}

func TestAccSonarqubeQualitygateProjectAssociationGateID(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate_project_association." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygateProjectAssociationGateID(rnd, "testAccSonarqubeProjectAssociationGateID"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "gateid"),
					resource.TestCheckResourceAttr(name, "projectkey", "testAccSonarqubeProjectAssociationGateID"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				// If sonarqube version is > 8.0, this test case will fail because of a missing gateid
				ExpectError: regexp.MustCompile("Error: API returned an error: No quality gate has been found for name"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "gateid"),
					resource.TestCheckResourceAttr(name, "projectkey", "testAccSonarqubeProjectAssociationGateID"),
				),
			},
		},
	})
}*/
