package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_alm_azure", &resource.Sweeper{
		Name: "sonarqube_alm_azure",
		F:    testSweepSonarqubeAlmAzure,
	})
}

// TODO: implement sweeper to clean up projects: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeAlmAzure(r string) error {
	return nil
}

func testAccSonarqubeAlmAzureName(rnd string, name string, url string) string {
	return fmt.Sprintf(`
		
		resource "sonarqube_alm_azure" "%[1]s" {
			key    = "%[2]s"
			personal_access_token    = "my_pat"
			url    = "%[3]s"
		}`, rnd, name, url)
}

func TestAccSonarqubeAlmAzureName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_alm_azure." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeAlmAzureName(rnd, "testAccSonarqubeAlmAzureName", "https://dev.azure.com/my-org"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubeAlmAzureName"),
					resource.TestCheckResourceAttr(name, "url", "https://dev.azure.com/my-org"),
				),
			},
			{
				Config: testAccSonarqubeAlmAzureName(rnd, "testAccSonarqubeAlmAzureNameUpdate", "https://dev.azure.com/my-other-org"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubeAlmAzureNameUpdate"),
					resource.TestCheckResourceAttr(name, "url", "https://dev.azure.com/my-other-org"),
				),
			},
		},
	})
}
