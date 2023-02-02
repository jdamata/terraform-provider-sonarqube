package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_github_binding", &resource.Sweeper{
		Name: "sonarqube_github_binding",
		F:    testSweepSonarqubeGithubBinding,
	})
}

// TODO: implement sweeper to clean up projects: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeGithubBinding(r string) error {
	return nil
}

func testAccSonarqubeGithubBindingName(rnd string, projName string, almSetting string) string {
	return fmt.Sprintf(`
		
		resource "sonarqube_alm_github" "%[1]s" {
			appid       = "12345"
			clientid    = "56789"
			clientsecret    = "secret"
			key    = "%[3]s"
			privatekey    = "myprivatekey"
			url    = "https://api.github.com"
			webhooksecret = "mysecret"
		}

		resource "sonarqube_project" "%[1]s" {
			name       = "%[2]s"
			project    = "%[2]s"
			visibility = "public"
		}
		resource "sonarqube_github_binding" "%[1]s" {
			almsetting   = "%[3]s"
			monorepo     = "no"
			project = sonarqube_project.%[1]s.project
			repository   = sonarqube_project.%[1]s.project
			summarycommentenabled = "true"
		    depends_on = [sonarqube_alm_github.%[1]s]
		}`, rnd, projName, almSetting)
}

func TestAccSonarqubeGithubBindingName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_github_binding." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeGithubBindingName(rnd, "testAccSonarqubeGithubBindingName", "github"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "almsetting", "github"),
				),
			},
			{
				Config: testAccSonarqubeGithubBindingName(rnd, "testAccSonarqubeGithubBindingName", "githubb"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "almsetting", "githubb"),
				),
			},
		},
	})
}
