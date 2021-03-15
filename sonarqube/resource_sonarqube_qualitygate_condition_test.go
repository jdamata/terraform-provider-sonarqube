package sonarqube

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_qualitygate_condition", &resource.Sweeper{
		Name: "sonarqube_qualitygate_condition",
		F:    testSweepSonarqubeQualitygateConditionSweeper,
	})
}

// TODO: implement sweeper to clean up projects: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubeQualitygateConditionSweeper(r string) error {
	return nil
}

func testAccSonarqubeQualitygateConditionGateName(id string, name string, metric string, threshold string, op string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualitygate" "%[1]s" {
			name = "%[2]s"
		}

		resource "sonarqube_qualitygate_condition" "%[1]s" {
			gatename  = sonarqube_qualitygate.%[1]s.name
			metric    = "%[3]s"
			threshold = "%[4]s"
			op        = "%[5]s"
		}
		`, id, name, metric, threshold, op)
}

func TestAccSonarqubeQualitygateConditionGateName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate_condition." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygateConditionGateName(rnd, "testAccSonarqubeQualitygateCondition", "vulnerabilities", "10", "GT"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "metric", "vulnerabilities"),
					resource.TestCheckResourceAttr(name, "threshold", "10"),
					resource.TestCheckResourceAttr(name, "op", "GT"),
				),
			},
			{
				Config: testAccSonarqubeQualitygateConditionGateName(rnd, "testAccSonarqubeQualitygateCondition", "vulnerabilities", "11", "GT"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "metric", "vulnerabilities"),
					resource.TestCheckResourceAttr(name, "threshold", "11"),
					resource.TestCheckResourceAttr(name, "op", "GT"),
				),
			},
		},
	})
}

func testAccSonarqubeQualitygateConditionGateID(id string, name string, metric string, threshold string, op string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualitygate" "%[1]s" {
			name = "%[2]s"
		}

		resource "sonarqube_qualitygate_condition" "%[1]s" {
			gateid    = sonarqube_qualitygate.%[1]s.id
			metric    = "%[3]s"
			threshold = "%[4]s"
			op        = "%[5]s"
		}
		`, id, name, metric, threshold, op)
}

func TestAccSonarqubeQualitygateConditionGateID(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate_condition." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// If sonarqube version is <8.0, this test case will fail because of a missing gatenmame
				ExpectError: regexp.MustCompile("Error: API returned an error: No quality gate has been found for name"),
				Config:      testAccSonarqubeQualitygateConditionGateID(rnd, "testAccSonarqubeQualitygateCondition", "vulnerabilities", "10", "GT"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "metric", "vulnerabilities"),
					resource.TestCheckResourceAttr(name, "threshold", "10"),
					resource.TestCheckResourceAttr(name, "op", "GT"),
				),
			},
			{
				// If sonarqube version is <8.0, this test case will fail because of a missing gatenmame
				ExpectError: regexp.MustCompile("Error: API returned an error: No quality gate has been found for name"),
				Config:      testAccSonarqubeQualitygateConditionGateID(rnd, "testAccSonarqubeQualitygateCondition", "vulnerabilities", "11", "GT"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "metric", "vulnerabilities"),
					resource.TestCheckResourceAttr(name, "threshold", "11"),
					resource.TestCheckResourceAttr(name, "op", "GT"),
				),
			},
		},
	})
}
