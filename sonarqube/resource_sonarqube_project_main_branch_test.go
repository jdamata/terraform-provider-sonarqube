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

func testAccSonarqubeProjectMainBranchConflictingBranch(rnd string, projName string, mainBranchName string, extraBranchName string) string {
	return fmt.Sprintf(`
		resource "sonarqube_project" "%[1]s" {
			name       = "%[2]s"
			project    = "%[2]s"
			visibility = "public"
		}

		resource "sonarqube_project_main_branch" "%[1]s" {
			name    = "%[3]s"
			project = sonarqube_project.%[1]s.project
		}

		resource "sonarqube_project_branch" "%[1]s_extra" {
			name    = "%[4]s"
			project = sonarqube_project.%[1]s.project
		}`, rnd, projName, mainBranchName, extraBranchName)
}

// TestAccSonarqubeProjectMainBranchConflictResolution tests that when updating
// the main branch to a name that already exists as a non-main branch, the provider
// correctly deletes the conflicting branch first (for SonarQube < 10.2 using rename API)
func TestAccSonarqubeProjectMainBranchConflictResolution(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_project_main_branch." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Step 1: Create a project with main branch named "main" and an extra branch "develop"
				Config: testAccSonarqubeProjectMainBranchConflictingBranch(rnd, "testAccConflictResolution", "main", "develop"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccConflictResolution"),
					resource.TestCheckResourceAttr(name, "name", "main"),
				),
			},
			{
				// Step 2: Update main branch to "develop" - should delete the existing non-main "develop" branch
				// and then rename "main" to "develop" (for SonarQube < 10.2) or set "develop" as main (>= 10.2)
				Config: testAccSonarqubeProjectMainBranchName(rnd, "testAccConflictResolution", "develop"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccConflictResolution"),
					resource.TestCheckResourceAttr(name, "name", "develop"),
				),
			},
			{
				// Step 3: Verify import still works after conflict resolution
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "project", "testAccConflictResolution"),
					resource.TestCheckResourceAttr(name, "name", "develop"),
				),
			},
		},
	})
}
