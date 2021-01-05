package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_user", &resource.Sweeper{
		Name: "sonarqube_user",
		F:    testSweepSonarqubeUserSweeper,
	})
}

func testSweepSonarqubeUserSweeper(r string) error {
	return nil
}

func testAccSonarqubeUserLocalConfig(rnd string, name string) string {
	return fmt.Sprintf(`
		resource "sonarqube_user" "%[1]s" {
			login_name = "%[2]s"
			name       = "%[2]s"
			password   = "secret-sauce37!"
		}`, rnd, name)
}

func TestAccSonarqubeUserLocal(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_user." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeUserLocalConfig(rnd, "testAccSonarqubeUserLocal"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeUserLocal"),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeUserLocal"),
				),
			},
		},
	})
}

func testAccSonarqubeUserNotLocalConfig(rnd string, name string) string {
	return fmt.Sprintf(`
		resource "sonarqube_user" "%[1]s" {
			login_name = "%[2]s"
			name       = "%[2]s"
			email      = "terraform-test@sonarqube.com"
			is_local   = false
	    }`, rnd, name)
}

func TestAccSonarqubeUserNotLocal(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_user." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeUserNotLocalConfig(rnd, "testAccSonarqubeUserNotLocal"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeUserNotLocal"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeUserNotLocal"),
				),
			},
		},
	})
}
