package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_bitbucket_binding", &resource.Sweeper{
		Name: "sonarqube_bitbucket_binding",
		F:    testSweepSonarqubeBitbucketBinding,
	})
}

// TODO: implement sweeper to clean up projects: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeBitbucketBinding(r string) error {
	return nil
}

func testAccPreCheckBitbucketBindingSupport(t *testing.T) {
	if err := checkBitbucketBindingSupport(testAccProvider.Meta().(*ProviderConfiguration)); err != nil {
		t.Skipf("Skipping test of unsupported feature (Bitbucket Binding)")
	}
}

func testAccSonarqubeBitbucketBindingName(rnd string, name string, repository string, slug string) string {
	return fmt.Sprintf(`
		resource "sonarqube_project" "%[1]s" {
			name       = "%[2]s"
			project    = "%[2]s"
			visibility = "public"
		}

        resource "sonarqube_alm_bitbucket" "%[1]s" {
            personal_access_token = "test-token"
            key                   = "%[2]s"
            url                   = "https://bitbucket.example.com"
        }

        resource "sonarqube_bitbucket_binding" "%[1]s" {
            alm_setting = sonarqube_alm_bitbucket.%[1]s.key
            monorepo    = "false"
            project     = sonarqube_project.%[1]s.project
            repository  = "%[3]s"
            slug        = "%[4]s"
        }`, rnd, name, repository, slug)
}

func TestAccSonarqubeBitbucketBindingName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_bitbucket_binding." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckBitbucketBindingSupport(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeBitbucketBindingName(rnd, "testAccSonarqubeBitbucketBindingName", "repo-key", "repo-slug"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeBitbucketBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "repo-key"),
					resource.TestCheckResourceAttr(name, "slug", "repo-slug"),
					resource.TestCheckResourceAttr(name, "monorepo", "false"),
				),
			},
			{
				Config: testAccSonarqubeBitbucketBindingName(rnd, "testAccSonarqubeBitbucketBindingName", "updated-repo-key", "updated-repo-slug"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeBitbucketBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "updated-repo-key"),
					resource.TestCheckResourceAttr(name, "slug", "updated-repo-slug"),
					resource.TestCheckResourceAttr(name, "monorepo", "false"),
				),
			},
		},
	})
}
