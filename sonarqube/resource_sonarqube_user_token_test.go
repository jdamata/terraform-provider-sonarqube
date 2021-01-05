package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_user_token", &resource.Sweeper{
		Name: "sonarqube_user_token",
		F:    testSweepSonarqubeUserTokenSweeper,
	})
}

func testSweepSonarqubeUserTokenSweeper(r string) error {
	return nil
}

func testAccSonarqubeUserTokenBasicConfig(rnd string, name string) string {
	return fmt.Sprintf(`
		resource "sonarqube_user" "%[1]s" {
			login_name = "%[2]s"
			name       = "%[2]s"
			password   = "secret-sauce37!"
		}
		resource "sonarqube_user_token" "%[1]s" {
			login_name = sonarqube_user.%[1]s.login_name
			name       = "%[2]s"
		}`, rnd, name)
}

func TestAccSonarqubeUserTokenBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_user_token." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeUserTokenBasicConfig(rnd, "testAccSonarqubeUserToken"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeUserToken"),
				),
			},
		},
	})
}
