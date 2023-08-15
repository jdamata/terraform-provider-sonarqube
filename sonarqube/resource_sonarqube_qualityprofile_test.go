package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_qualityprofile", &resource.Sweeper{
		Name: "sonarqube_qualityprofile",
		F:    testSweepSonarqubeQualityProfileSweeper,
	})
}

func testSweepSonarqubeQualityProfileSweeper(r string) error {
	return nil
}

func testAccSonarqubeQualityProfileBasicConfig(rnd string, name string, language string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualityprofile" "%[1]s" {
			name     = "%[2]s"
			language = "%[3]s"
		}`, rnd, name, language)
}

func TestAccSonarqubeQualityProfileBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualityprofile." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualityProfileBasicConfig(rnd, "testAccSonarqubeQualityProfile", "js"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeQualityProfile"),
					resource.TestCheckResourceAttr(name, "language", "js"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeQualityProfile"),
					resource.TestCheckResourceAttr(name, "language", "js"),
					resource.TestCheckResourceAttrSet(name, "key"),
				),
			},
		},
	})
}
