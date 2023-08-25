package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccSonarqubeProjectDataSourceConfig(rnd string, name string, project string, visibility string) string {
	return fmt.Sprintf(`
		resource "sonarqube_project" "%[1]s" {
		  name       = "%[2]s"
		  project    = "%[3]s"
		  visibility = "%[4]s"
		}
		data "sonarqube_project" "%[1]s" {
			project = sonarqube_project.%[1]s.id
		}
		`, rnd, name, project, visibility)
}

func TestAccSonarqubeProjectDataSource(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "data.sonarqube_project." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeProjectDataSourceConfig(rnd, "testAccSonarqubeProjectDataSource", "testAccSonarqubeProjectDataSource", "public"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeProjectDataSource"),
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProjectDataSource"),
					resource.TestCheckResourceAttr(name, "visibility", "public"),
				),
			},
		},
	})
}
