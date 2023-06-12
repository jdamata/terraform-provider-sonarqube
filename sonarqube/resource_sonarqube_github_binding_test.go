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
func testAccPreCheckGithubBindingSupport(t *testing.T) {
	if err := checkGithubBindingSupport(testAccProvider.Meta().(*ProviderConfiguration)); err != nil {
		t.Skipf("Skipping test of unsupported feature (GitHub Binding)")
	}
}

func testAccSonarqubeGithubBindingName(rnd string, projName string, almSetting string, repoName string) string {
	return fmt.Sprintf(`
		
		resource "sonarqube_alm_github" "%[1]s" {
			app_id       = "12345"
			client_id    = "56789"
			client_secret    = "secret"
			key    = "%[3]s"
			private_key    = "myprivate_key"
			url    = "https://api.github.com"
			webhook_secret = "mysecret"
		}

		resource "sonarqube_project" "%[1]s" {
			name       = "%[2]s"
			project    = "%[2]s"
			visibility = "public"
		}
		resource "sonarqube_github_binding" "%[1]s" {
			alm_setting   = "%[3]s"
			monorepo     = "false"
			project = sonarqube_project.%[1]s.project
			repository   = "%[4]s"
			summary_comment_enabled = "true"
		    depends_on = [sonarqube_alm_github.%[1]s]
		}`, rnd, projName, almSetting, repoName)
}

func TestAccSonarqubeGithubBindingName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_github_binding." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckGithubBindingSupport(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeGithubBindingName(rnd, "testAccSonarqubeGithubBindingName", "github", "testAccSonarqubeGithubBindingName"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "github"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeGithubBindingName"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "github"),
				),
			},
			{
				Config: testAccSonarqubeGithubBindingName(rnd, "testAccSonarqubeGithubBindingName", "githubb", "testAccSonarqubeGithubBindingName"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "githubb"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "githubb"),
				),
			},
			{
				Config: testAccSonarqubeGithubBindingName(rnd, "testAccSonarqubeGithubBindingName", "GitHub", "org/testAccSonarqubeGithubBindingName"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "org/testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "GitHub"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "org/testAccSonarqubeGithubBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "GitHub"),
				),
			},
		},
	})
}
