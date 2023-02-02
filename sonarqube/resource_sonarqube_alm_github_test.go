package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_alm_github", &resource.Sweeper{
		Name: "sonarqube_alm_github",
		F:    testSweepSonarqubeAlmGithub,
	})
}

// TODO: implement sweeper to clean up projects: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeAlmGithub(r string) error {
	return nil
}

func testAccSonarqubeAlmGithubName(rnd string, name string, appId string, clientId string) string {
	return fmt.Sprintf(`
		
		resource "sonarqube_alm_github" "%[1]s" {
			appid       = "%[3]s"
			clientid    = "%[4]s"
			clientsecret    = "secret"
			key    = "%[2]s"
			privatekey    = "myprivatekey"
			url    = "https://api.github.com"
			webhooksecret = "mysecret"
		}`, rnd, name, appId, clientId)
}

func TestAccSonarqubeAlmGithubName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_alm_github." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeAlmGithubName(rnd, "testAccSonarqubeAlmGithubName", "123456", "234567"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubeAlmGithubName"),
					resource.TestCheckResourceAttr(name, "appid", "123456"),
					resource.TestCheckResourceAttr(name, "clientid", "234567"),
				),
			},
			{
				Config: testAccSonarqubeAlmGithubName(rnd, "testAccSonarqubeAlmGithubNameUpdate", "654321", "765432"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubeAlmGithubNameUpdate"),
					resource.TestCheckResourceAttr(name, "appid", "654321"),
					resource.TestCheckResourceAttr(name, "clientid", "765432"),
				),
			},
		},
	})
}
