package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_qualityprofile_project_association", &resource.Sweeper{
		Name: "sonarqube_qualityprofile_project_association",
		F:    testSweepSonarqubeQualityProfileProjectAssociationSweeper,
	})
}

func testSweepSonarqubeQualityProfileProjectAssociationSweeper(r string) error {
	return nil
}

func testAccSonarqubeQualityProfileProjectAssociationBasicConfig(rnd string, name string, language string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualityprofile" "%[1]s" {
			name     = "%[2]s"
			language = "%[3]s"
		}

		resource "sonarqube_project" "%[1]s" {
			name       = "%[2]s"
			project    = "%[2]s"
			visibility = "public" 
		}

		resource "sonarqube_qualityprofile_project_association" "%[1]s" {
			quality_profile = sonarqube_qualityprofile.%[1]s.name
			project         = sonarqube_project.%[1]s.name
			language        = "%[3]s"
		}`, rnd, name, language)
}

func TestAccSonarqubeQualityProfileProjectAssociationBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualityprofile_project_association." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualityProfileProjectAssociationBasicConfig(rnd, "testAccSonarqubeProfileProjectAssociation", "js"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "quality_profile", "testAccSonarqubeProfileProjectAssociation"),
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProfileProjectAssociation"),
					resource.TestCheckResourceAttr(name, "language", "js"),
				),
			},
		},
	})
}
