package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_new_code_periods", &resource.Sweeper{
		Name: "sonarqube_new_code_periods",
		F:    testSweepSonarqubeNewCodePeriods,
	})
}

// TODO: implement sweeper to clean up projects: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeNewCodePeriods(r string) error {
	return nil
}

func testAccSonarqubeNewCodePeriodsGlobalPreviousVersion(rnd string) string {
	return fmt.Sprintf(`
        resource "sonarqube_new_code_periods" "%[1]s" {
			type = "PREVIOUS_VERSION"
        }`, rnd)
}

func TestAccSonarqubeNewCodePeriodsGlobalPreviousVersion(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_new_code_periods." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeNewCodePeriodsGlobalPreviousVersion(rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "type", "PREVIOUS_VERSION"),
				),
			},
		},
	})
}

func testAccSonarqubeNewCodePeriodsGlobalNumberOfDays(rnd string) string {
	return fmt.Sprintf(`
        resource "sonarqube_new_code_periods" "%[1]s" {
			type = "NUMBER_OF_DAYS"
			value = "5"
        }`, rnd)
}

func TestAccSonarqubeNewCodePeriodsGlobalNumberOfDays(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_new_code_periods." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeNewCodePeriodsGlobalNumberOfDays(rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "type", "NUMBER_OF_DAYS"),
					resource.TestCheckResourceAttr(name, "value", "5"),
				),
			},
		},
	})
}

func testAccSonarqubeNewCodePeriodsBranchPreviousVersion(rnd string) string {
	return fmt.Sprintf(`
	    resource "sonarqube_project" "%[1]s" {
			name = "%[1]s"
			project = "%[1]s"
			visibility = "public"
		}

        resource "sonarqube_new_code_periods" "%[1]s" {
			branch = "main"
			project = sonarqube_project.%[1]s.project
			type = "PREVIOUS_VERSION"
        }`, rnd)
}

func TestAccSonarqubeNewCodePeriodsBranchPreviousVersion(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_new_code_periods." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeNewCodePeriodsBranchPreviousVersion(rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "type", "PREVIOUS_VERSION"),
					resource.TestCheckResourceAttr(name, "branch", "main"),
					resource.TestCheckResourceAttr(name, "project", rnd),
				),
			},
		},
	})
}

func testAccSonarqubeNewCodePeriodsBranchNumberOfDays(rnd string) string {
	return fmt.Sprintf(`
	    resource "sonarqube_project" "%[1]s" {
			name = "%[1]s"
			project = "%[1]s"
			visibility = "public"
		}

        resource "sonarqube_new_code_periods" "%[1]s" {
			branch = "main"
			project = sonarqube_project.%[1]s.project
			type = "NUMBER_OF_DAYS"
			value = "5"
        }`, rnd)
}

func TestAccSonarqubeNewCodePeriodsBranchNumberOfDays(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_new_code_periods." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeNewCodePeriodsBranchNumberOfDays(rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "type", "NUMBER_OF_DAYS"),
					resource.TestCheckResourceAttr(name, "branch", "main"),
					resource.TestCheckResourceAttr(name, "project", rnd),
					resource.TestCheckResourceAttr(name, "value", "5"),
				),
			},
		},
	})
}

// Note sure how to get a specific analysis UUID on the fly in a test environment
//
// func testAccSonarqubeNewCodePeriodsBranchSpecificAnalysis(rnd string) string {
// 	return fmt.Sprintf(`
// 	    resource "sonarqube_project" "%[1]s" {
// 			name = "%[1]s"
// 			project = "%[1]s"
// 			visibility = "public"
// 		}
//
//         resource "sonarqube_new_code_periods" "%[1]s" {
// 			branch = "main"
// 			project = sonarqube_project.%[1]s.project
// 			type = "SPECIFIC_ANALYSIS"
//         }`, rnd)
// }
//
// func TestAccSonarqubeNewCodePeriodsBranchSpecificAnalysis(t *testing.T) {
// 	rnd := generateRandomResourceName()
// 	name := "sonarqube_new_code_periods." + rnd
//
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:  func() { testAccPreCheck(t) },
// 		Providers: testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccSonarqubeNewCodePeriodsBranchSpecificAnalysis(rnd),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(name, "type", "SPECIFIC_ANALYSIS"),
// 					resource.TestCheckResourceAttr(name, "branch", "main"),
// 					resource.TestCheckResourceAttr(name, "project", rnd),
// 				),
// 			},
// 		},
// 	})
// }

func testAccSonarqubeNewCodePeriodsBranchReferenceBranch(rnd string) string {
	return fmt.Sprintf(`
	    resource "sonarqube_project" "%[1]s" {
			name = "%[1]s"
			project = "%[1]s"
			visibility = "public"
		}

        resource "sonarqube_new_code_periods" "%[1]s" {
			branch = "main"
			project = sonarqube_project.%[1]s.project
			type = "REFERENCE_BRANCH"
			value = "development"
        }`, rnd)
}

func TestAccSonarqubeNewCodePeriodsBranchReferenceBranch(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_new_code_periods." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeNewCodePeriodsBranchReferenceBranch(rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "type", "REFERENCE_BRANCH"),
					resource.TestCheckResourceAttr(name, "branch", "main"),
					resource.TestCheckResourceAttr(name, "project", rnd),
					resource.TestCheckResourceAttr(name, "value", "development"),
				),
			},
		},
	})
}

func testAccSonarqubeNewCodePeriodsProjectPreviousVersion(rnd string) string {
	return fmt.Sprintf(`
	    resource "sonarqube_project" "%[1]s" {
			name = "%[1]s"
			project = "%[1]s"
			visibility = "public"
		}

        resource "sonarqube_new_code_periods" "%[1]s" {
			project = sonarqube_project.%[1]s.project
			type = "PREVIOUS_VERSION"
        }`, rnd)
}

func TestAccSonarqubeNewCodePeriodsProjectPreviousVersion(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_new_code_periods." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeNewCodePeriodsProjectPreviousVersion(rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "type", "PREVIOUS_VERSION"),
					resource.TestCheckResourceAttr(name, "project", rnd),
				),
			},
		},
	})
}

func testAccSonarqubeNewCodePeriodsProjectNumberOfDays(rnd string) string {
	return fmt.Sprintf(`
	    resource "sonarqube_project" "%[1]s" {
			name = "%[1]s"
			project = "%[1]s"
			visibility = "public"
		}

        resource "sonarqube_new_code_periods" "%[1]s" {
			project = sonarqube_project.%[1]s.project
			type = "NUMBER_OF_DAYS"
			value = "5"
        }`, rnd)
}

func TestAccSonarqubeNewCodePeriodsProjectNumberOfDays(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_new_code_periods." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeNewCodePeriodsProjectNumberOfDays(rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "type", "NUMBER_OF_DAYS"),
					resource.TestCheckResourceAttr(name, "project", rnd),
					resource.TestCheckResourceAttr(name, "value", "5"),
				),
			},
		},
	})
}

func testAccSonarqubeNewCodePeriodsProjectReferenceProject(rnd string) string {
	return fmt.Sprintf(`
	    resource "sonarqube_project" "%[1]s" {
			name = "%[1]s"
			project = "%[1]s"
			visibility = "public"
		}

        resource "sonarqube_new_code_periods" "%[1]s" {
			project = sonarqube_project.%[1]s.project
			type = "REFERENCE_BRANCH"
			value = "development"
        }`, rnd)
}

func TestAccSonarqubeNewCodePeriodsProjectReferenceProject(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_new_code_periods." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeNewCodePeriodsProjectReferenceProject(rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "type", "REFERENCE_BRANCH"),
					resource.TestCheckResourceAttr(name, "project", rnd),
					resource.TestCheckResourceAttr(name, "value", "development"),
				),
			},
		},
	})
}
