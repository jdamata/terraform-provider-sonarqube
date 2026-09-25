package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_alm_bitbucket", &resource.Sweeper{
		Name: "sonarqube_alm_bitbucket",
		F:    testSweepSonarqubeAlmBitbucket,
	})
}

// TODO: implement sweeper to clean up projects: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeAlmBitbucket(r string) error {
	return nil
}

func testAccSonarqubeAlmBitbucketName(rnd string, name string, personalAccessToken string) string {
	return fmt.Sprintf(`

        resource "sonarqube_alm_bitbucket" "%[1]s" {
            personal_access_token       = "%[3]s"
            key    = "%[2]s"
            url    = "https://bitbucket.%[3]s.com"
        }`, rnd, name, personalAccessToken)
}

func TestAccSonarqubeAlmBitbucketName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_alm_bitbucket." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeAlmBitbucketName(rnd, "testAccSonarqubeAlmBitbucketName", "123456"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubeAlmBitbucketName"),
					resource.TestCheckResourceAttr(name, "personal_access_token", "123456"),
					resource.TestCheckResourceAttr(name, "url", "https://bitbucket.123456.com"),
				),
			},
			{
				Config: testAccSonarqubeAlmBitbucketName(rnd, "testAccSonarqubeAlmBitbucketNameUpdate", "654321"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubeAlmBitbucketNameUpdate"),
					resource.TestCheckResourceAttr(name, "personal_access_token", "654321"),
					resource.TestCheckResourceAttr(name, "url", "https://bitbucket.654321.com"),
				),
			},
		},
	})
}
