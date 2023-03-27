package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccSonarqubePortfolioDataSourceConfig(rnd string, key string, name string, description string, visibility string) string {
	return fmt.Sprintf(`
		resource "sonarqube_portfolio" "%[1]s" {
			key       = "%[2]s"
			name    = "%[3]s"
			description = "%[4]s"
			visibility = "%[5]s"
		}
		data "sonarqube_portfolio" "%[1]s" {
			key = sonarqube_portfolio.%[1]s.id
		}
		`, rnd, key, name, description, visibility)
}

func TestAccSonarqubePortfolioDataSource(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "data.sonarqube_portfolio." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePortfolioDataSourceConfig(rnd, "testAccSonarqubePortfolioDataSourceKey", "testAccSonarqubePortfolioDataSourceName", "testAccSonarqubePortfolioDataSourceDescription", "public"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioDataSourceKey"),
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubePortfolioDataSourceName"),
					resource.TestCheckResourceAttr(name, "description", "testAccSonarqubePortfolioDataSourceDescription"),
					resource.TestCheckResourceAttr(name, "visibility", "public"),
				),
			},
		},
	})
}
