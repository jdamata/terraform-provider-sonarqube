package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

func testAccSonarqubeGroupBasicConfig(rnd string, name string, description string) string {
	return fmt.Sprintf(`
		resource "sonarqube_group" "%[1]s" {
		  name        = "%[2]s"
		  description = "%[3]s"
		}
		`, rnd, name, description)
}

func TestAccSonarqubeGroupBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_group." + rnd
	groupName := "testAccSonarqubeGroup" + rnd
	groupDescription := "testAccSonarqubeDescription" + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeGroupBasicConfig(rnd, groupName, groupDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", groupName),
					resource.TestCheckResourceAttr(name, "description", groupDescription),
				),
			},
		},
	})
}

func TestAccSonarqubeGroupUpdate(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_group." + rnd
	groupName := "testAccSonarqubeGroup" + rnd
	updatedGroupName := "testAccSonarqubeGroupUpdated" + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeGroupBasicConfig(rnd, groupName, "group description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", groupName),
					resource.TestCheckResourceAttr(name, "description", "group description"),
				),
			},
			{
				Config: testAccSonarqubeGroupBasicConfig(rnd, groupName, "group description 2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", groupName),
					resource.TestCheckResourceAttr(name, "description", "group description 2"),
				),
			},
			{
				Config: testAccSonarqubeGroupBasicConfig(rnd, updatedGroupName, "group description 3"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", updatedGroupName),
					resource.TestCheckResourceAttr(name, "description", "group description 3"),
				),
			},
		},
	})
}

func TestAccSonarqubeGroupImport(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_group." + rnd
	groupName := "testAccSonarqubeGroup" + rnd
	groupDescription := "testAccSonarqubeGroupDescription" + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeGroupBasicConfig(rnd, groupName, groupDescription),
			},
			{
				ResourceName: name,
				ImportState:  true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", groupName),
					resource.TestCheckResourceAttr(name, "description", groupDescription),
				),
			},
		},
	})
}
