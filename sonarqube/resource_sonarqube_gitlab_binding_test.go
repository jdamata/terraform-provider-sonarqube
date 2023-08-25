package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_gitlab_binding", &resource.Sweeper{
		Name: "sonarqube_gitlab_binding",
		F:    testSweepSonarqubeGitlabBinding,
	})
}

// TODO: implement sweeper to clean up projects: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeGitlabBinding(r string) error {
	return nil
}
func testAccPreCheckGitlabBindingSupport(t *testing.T) {
	if err := checkGitlabBindingSupport(testAccProvider.Meta().(*ProviderConfiguration)); err != nil {
		t.Skipf("Skipping test of unsupported feature (Gitlab Binding)")
	}
}

func testAccSonarqubeGitlabBindingName(rnd string, projName string, almSetting string, repoName string) string {
	return fmt.Sprintf(`

        resource "sonarqube_alm_gitlab" "%[1]s" {
            personal_access_token       = "123456"
            key    = "%[3]s"
            url    = "https://gitlab.com/api/v4"
        }

        resource "sonarqube_project" "%[1]s" {
            name       = "%[2]s"
            project    = "%[2]s"
            visibility = "public"
        }

        resource "sonarqube_gitlab_binding" "%[1]s" {
            alm_setting   = "%[3]s"
            monorepo     = "false"
            project = sonarqube_project.%[1]s.project
            repository   = "%[4]s"
            depends_on = [sonarqube_alm_gitlab.%[1]s]
        }`, rnd, projName, almSetting, repoName)
}

func TestAccSonarqubeGitlabBindingName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_gitlab_binding." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckGitlabBindingSupport(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeGitlabBindingName(rnd, "testAccSonarqubeGitlabBindingName", "gitlab", "testAccSonarqubeGitlabBindingName"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeGitlabBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeGitlabBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "gitlab"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeGitlabBindingName"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeGitlabBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeGitlabBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "gitlab"),
				),
			},
			{
				Config: testAccSonarqubeGitlabBindingName(rnd, "testAccSonarqubeGitlabBindingName", "gitlabb", "testAccSonarqubeGitlabBindingName"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeGitlabBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeGitlabBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "gitlabb"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeGitlabBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeGitlabBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "gitlabb"),
				),
			},
			{
				Config: testAccSonarqubeGitlabBindingName(rnd, "testAccSonarqubeGitlabBindingName", "GitLab", "org/testAccSonarqubeGitlabBindingName"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeGitlabBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "org/testAccSonarqubeGitlabBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "GitLab"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeGitlabBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "org/testAccSonarqubeGitlabBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "GitLab"),
				),
			},
		},
	})
}
