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

func testAccSonarqubePortfolioDataSourceConfigTags(rnd string, key string, name string, description string, visibility string, tags []string) string {
	formattedTags := generateHCLList(tags)
	return fmt.Sprintf(`
		resource "sonarqube_portfolio" "%[1]s" {
			key       = "%[2]s"
			name    = "%[3]s"
			description = "%[4]s"
			visibility = "%[5]s"
			selection_mode = "TAGS"
			tags = %[6]s // Note that the "" should be missing since this is a list
		}
		data "sonarqube_portfolio" "%[1]s" {
			key = sonarqube_portfolio.%[1]s.id
		}
		`, rnd, key, name, description, visibility, formattedTags)
}

func testAccSonarqubePortfolioDataSourceConfigRegex(rnd string, key string, name string, description string, visibility string, regexp string) string {
	return fmt.Sprintf(`
		resource "sonarqube_portfolio" "%[1]s" {
			key       = "%[2]s"
			name    = "%[3]s"
			description = "%[4]s"
			visibility = "%[5]s"
			selection_mode = "REGEX"
			rexep = "%[6]s"
		}
		data "sonarqube_portfolio" "%[1]s" {
			key = sonarqube_portfolio.%[1]s.id
		}
		`, rnd, key, name, description, visibility, regexp)
}

func TestAccSonarqubePortfolioDataSource(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "data.sonarqube_portfolio." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckPortfolioSupport(t) },
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

func TestAccSonarqubePortfolioDataSourceTags(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "data.sonarqube_portfolio." + rnd
	tags := []string{"tag1", "tag2"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckPortfolioSupport(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePortfolioDataSourceConfigTags(rnd, "testAccSonarqubePortfolioDataSourceKey", "testAccSonarqubePortfolioDataSourceName", "testAccSonarqubePortfolioDataSourceDescription", "public", tags),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioDataSourceKey"),
					resource.TestCheckResourceAttr(name, "selection_mode", "TAGS"),
					resource.TestCheckResourceAttr(name, "tags.0", tags[0]),
					resource.TestCheckResourceAttr(name, "tags.1", tags[1]),
				),
			},
		},
	})
}

func TestAccSonarqubePortfolioDataSourceRegex(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "data.sonarqube_portfolio." + rnd
	exampleRegex := "myExampleRegex"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckPortfolioSupport(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePortfolioDataSourceConfigRegex(rnd, "testAccSonarqubePortfolioDataSourceKey", "testAccSonarqubePortfolioDataSourceName", "testAccSonarqubePortfolioDataSourceDescription", "public", exampleRegex),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioDataSourceKey"),
					resource.TestCheckResourceAttr(name, "selection_mode", "REGEXP"),
					resource.TestCheckResourceAttr(name, "regexp", exampleRegex),
				),
			},
		},
	})
}
