package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_alm_gitlab", &resource.Sweeper{
		Name: "sonarqube_alm_gitlab",
		F:    testSweepSonarqubeAlmGitlab,
	})
}

// TODO: implement sweeper to clean up projects: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeAlmGitlab(r string) error {
	return nil
}

func testAccSonarqubeAlmGitlabName(rnd string, name string, personalAccessToken string) string {
	return fmt.Sprintf(`

        resource "sonarqube_alm_gitlab" "%[1]s" {
            personal_access_token       = "%[3]s"
            key    = "%[2]s"
            url    = "https://%[3]s.gitlab.com/api/v4"
        }`, rnd, name, personalAccessToken)
}

func TestAccSonarqubeAlmGitlabName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_alm_gitlab." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeAlmGitlabName(rnd, "testAccSonarqubeAlmGitlabName", "123456"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubeAlmGitlabName"),
					resource.TestCheckResourceAttr(name, "personal_access_token", "123456"),
					resource.TestCheckResourceAttr(name, "url", "https://123456.gitlab.com/api/v4"),
				),
			},
			{
				Config: testAccSonarqubeAlmGitlabName(rnd, "testAccSonarqubeAlmGitlabNameUpdate", "654321"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubeAlmGitlabNameUpdate"),
					resource.TestCheckResourceAttr(name, "personal_access_token", "654321"),
					resource.TestCheckResourceAttr(name, "url", "https://654321.gitlab.com/api/v4"),
				),
			},
		},
	})
}
