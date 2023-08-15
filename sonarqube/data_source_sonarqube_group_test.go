package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccSonarqubeGroupDataSourceConfig(rnd string, name string, description string) string {
	return fmt.Sprintf(`
		resource "sonarqube_group" "%[1]s" {
			name        = "%[2]s"
			description = "%[3]s"
		}
		
		data "sonarqube_group" "%[1]s" {
			name = sonarqube_group.%[1]s.name
		}`, rnd, name, description)
}

func TestAccSonarqubeGroupDataSource(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "data.sonarqube_group." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeGroupDataSourceConfig(rnd, "testAccSonarqubeGroupDataSource", "Terraform Test Group Data-source"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeGroupDataSource"),
					resource.TestCheckResourceAttr(name, "description", "Terraform Test Group Data-source"),
				),
			},
		},
	})
}
