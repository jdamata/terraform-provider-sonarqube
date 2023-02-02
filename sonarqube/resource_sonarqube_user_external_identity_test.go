package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_user_external_identity", &resource.Sweeper{
		Name: "sonarqube_user_external_identity",
		F:    testSweepSonarqubeUserExteernalIdentitySweeper,
	})
}

func testSweepSonarqubeUserExteernalIdentitySweeper(r string) error {
	return nil
}

func testAccSonarqubeUserExternalIdentityConfig(rnd string, login string, externalIdentity string, externalProvider string) string {
	return fmt.Sprintf(`
		resource "sonarqube_user" "%[1]s" {
			login_name = "%[2]s"
			name       = "Test User"
			email      = "terraform-test@sonarqube.com"
			password   = "secret-sauce"
		}

		resource "sonarqube_user_external_identity" "%[1]s" {
			login_name         = sonarqube_user.%[1]s.login_name
			external_identity  = "%[3]s"
			external_provider  = "%[4]s"
		}`, rnd, login, externalIdentity, externalProvider)
}

func TestAccSonarqubeUserExternalIdentity(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_user_external_identity." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeUserExternalIdentityConfig(rnd, "testAccSonarqubeUser", "terraform-test@sonarqube.com", "saml"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "login_name", "testAccSonarqubeUser"),
					resource.TestCheckResourceAttr(name, "external_identity", "terraform-test@sonarqube.com"),
					resource.TestCheckResourceAttr(name, "external_provider", "saml"),
				),
			},
			{
				Config: testAccSonarqubeUserExternalIdentityConfig(rnd, "testAccSonarqubeUser", "terraform-test2@sonarqube.com", "LDAP"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "login_name", "testAccSonarqubeUser"),
					resource.TestCheckResourceAttr(name, "external_identity", "terraform-test2@sonarqube.com"),
					resource.TestCheckResourceAttr(name, "external_provider", "LDAP"),
				),
			},
		},
	})
}
