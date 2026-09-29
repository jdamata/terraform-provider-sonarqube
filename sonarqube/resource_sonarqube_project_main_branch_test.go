package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_project_main_branch", &resource.Sweeper{
		Name: "sonarqube_project_main_branch",
		F:    testSweepSonarqubeProjectMainBranchSweeper,
	})
}

// TODO: implement sweeper to clean up projects: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeProjectMainBranchSweeper(r string) error {
	return nil
}

func testAccSonarqubeProjectMainBranchName(rnd string, projName string, branchName string) string {
	return fmt.Sprintf(`
		resource "sonarqube_project" "%[1]s" {
			name       = "%[2]s"
			project    = "%[2]s"
			visibility = "public"
		}

		resource "sonarqube_project_main_branch" "%[1]s" {
			name   = "%[3]s"
			project = sonarqube_project.%[1]s.project
		}`, rnd, projName, branchName)
}

func TestAccSonarqubeProjectMainBranchName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_project_main_branch." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeProjectMainBranchName(rnd, "testAccSonarqubeProjectMainBranchName", "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProjectMainBranchName"),
					resource.TestCheckResourceAttr(name, "name", "test"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProjectMainBranchName"),
					resource.TestCheckResourceAttr(name, "name", "test"),
				),
			},
			{
				Config: testAccSonarqubeProjectMainBranchName(rnd, "testAccSonarqubeProjectMainBranchName", "main"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProjectMainBranchName"),
					resource.TestCheckResourceAttr(name, "name", "main"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProjectMainBranchName"),
					resource.TestCheckResourceAttr(name, "name", "main"),
				),
			},
			{
				Config: testAccSonarqubeProjectMainBranchName(rnd, "testAccSonarqubeProjectMainBranchName", "slash/branch"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProjectMainBranchName"),
					resource.TestCheckResourceAttr(name, "name", "slash/branch"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProjectMainBranchName"),
					resource.TestCheckResourceAttr(name, "name", "slash/branch"),
				),
			},
		},
	})
}

func TestAccSonarqubeProjectMainBranchSetMain(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_project_main_branch." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeProjectMainBranchName(rnd, "testAccSonarqubeProjectMainBranchSetMain", "main"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProjectMainBranchSetMain"),
					resource.TestCheckResourceAttr(name, "name", "main"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProjectMainBranchSetMain"),
					resource.TestCheckResourceAttr(name, "name", "main"),
				),
			},
		},
	})
}

func TestAccSonarqubeProjectMainBranchRenameNonExistent(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_project_main_branch." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeProjectMainBranchName(rnd, "testAccSonarqubeProjectMainBranchRenameNonExistent", "non-existent-branch"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProjectMainBranchRenameNonExistent"),
					resource.TestCheckResourceAttr(name, "name", "non-existent-branch"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProjectMainBranchRenameNonExistent"),
					resource.TestCheckResourceAttr(name, "name", "non-existent-branch"),
				),
			},
		},
	})
}
