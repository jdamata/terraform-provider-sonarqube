package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_group", &resource.Sweeper{
		Name: "sonarqube_group",
		F:    testSweepSonarqubeGroupSweeper,
	})
}

// TODO: implement sweeper to clean up groups: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeGroupSweeper(r string) error {
	return nil
}

func TestGroupSimple(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_group." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testSonarPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testGroupConfig(rnd, "group description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "description", "group description"),
				),
			},
		},
	})
}

func testGroupConfig(resourceID string, description string) string {
	return fmt.Sprintf(`
		resource "sonarqube_group" "%[1]s" {
		  name        = "%[1]s"
		  description = "%[2]s"
		}
		`, resourceID, description)
}
