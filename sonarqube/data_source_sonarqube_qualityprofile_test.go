package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccSonarqubeQualityProfileDataSourceConfig(rnd string, name string, language string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualityprofile" "%[1]s" {
			name     = "%[2]s"
			language = "%[3]s"
		}
	
		data "sonarqube_qualityprofile" "%[1]s" {
			key = sonarqube_qualityprofile.%[1]s.id
		}`, rnd, name, language)
}

func testAccSonarqubeQualityProfileDataSourceConfigNoId(rnd string, name string, language string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualityprofile" "%[1]s" {
			name     = "%[2]s"
			language = "%[3]s"
		}
	
		data "sonarqube_qualityprofile" "%[1]s" {
			name = sonarqube_qualityprofile.%[1]s.name
			language = sonarqube_qualityprofile.%[1]s.language
		}`, rnd, name, language)
}

func TestAccSonarqubeQualityProfileDataSource(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "data.sonarqube_qualityprofile." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualityProfileDataSourceConfig(rnd, "testAccSonarqubeQualityProfileDataSource", "js"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeQualityProfileDataSource"),
					resource.TestCheckResourceAttr(name, "language", "js"),
				),
			},
		},
	})
}

func TestAccSonarqubeQualityProfileDataSourceNoId(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "data.sonarqube_qualityprofile." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualityProfileDataSourceConfigNoId(rnd, "testAccSonarqubeQualityProfileDataSourceNoId", "cs"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeQualityProfileDataSourceNoId"),
					resource.TestCheckResourceAttr(name, "language", "cs"),
				),
			},
		},
	})
}
