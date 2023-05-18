package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_azure_binding", &resource.Sweeper{
		Name: "sonarqube_azure_binding",
		F:    testSweepSonarqubeAzureBinding,
	})
}

// TODO: implement sweeper to clean up projects: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeAzureBinding(r string) error {
	return nil
}
func testAccPreCheckAzureBindingSupport(t *testing.T) {
	if err := checkAzureBindingSupport(testAccProvider.Meta().(*ProviderConfiguration)); err != nil {
		t.Skipf("Skipping test of unsupported feature (Azure Binding)")
	}
}

func testAccSonarqubeAzureBindingName(rnd string, projName string, almSetting string, repoName string) string {
	return fmt.Sprintf(`
		
		resource "sonarqube_alm_azure" "%[1]s" {
			key    = "%[3]s"
			personal_access_token    = "my_pat"
			url    = "https://dev.azure.com/my-org"
		}

		resource "sonarqube_project" "%[1]s" {
			name       = "%[2]s"
			project    = "%[2]s"
			visibility = "public"
		}
		resource "sonarqube_azure_binding" "%[1]s" {
			alm_setting   = sonarqube_alm_azure.%[1]s.key
			monorepo     = "false"
			project = sonarqube_project.%[1]s.project
			repository   = "%[4]s"
		}`, rnd, projName, almSetting, repoName)
}

func TestAccSonarqubeAzureBindingName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_azure_binding." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckAzureBindingSupport(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeAzureBindingName(rnd, "testAccSonarqubeAzureBindingName", "azure", "testAccSonarqubeAzureBindingName"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeAzureBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeAzureBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "azure"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeAzureBindingName"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeAzureBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeAzureBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "azure"),
				),
			},
			{
				Config: testAccSonarqubeAzureBindingName(rnd, "testAccSonarqubeAzureBindingName", "azurea", "testAccSonarqubeAzureBindingName"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeAzureBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeAzureBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "azurea"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeAzureBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "testAccSonarqubeAzureBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "azurea"),
				),
			},
			{
				Config: testAccSonarqubeAzureBindingName(rnd, "testAccSonarqubeAzureBindingName", "azure", "org/testAccSonarqubeAzureBindingName"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeAzureBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "org/testAccSonarqubeAzureBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "azure"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeAzureBindingName"),
					resource.TestCheckResourceAttr(name, "repository", "org/testAccSonarqubeAzureBindingName"),
					resource.TestCheckResourceAttr(name, "alm_setting", "azure"),
				),
			},
		},
	})
}
