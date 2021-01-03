package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_project", &resource.Sweeper{
		Name: "sonarqube_project",
		F:    testSweepSonarqubeProjectSweeper,
	})
}

// TODO: implement sweeper to clean up projects: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeProjectSweeper(r string) error {
	return nil
}

func TestAccSonarqubeProject_basic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_project." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeProjectConfig(rnd, "my_project", "public"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "my_project"),
					resource.TestCheckResourceAttr(name, "visibility", "public"),
				),
			},
		},
	})
}

func testAccSonarqubeProjectConfig(name string, project string, visibility string) string {
	return fmt.Sprintf(`
		resource "sonarqube_project" "%[1]s" {
		  name       = "%[1]s"
		  project    = "%[2]s"
		  visibility = "%[3]s"
		}
		`, name, project, visibility)
}
