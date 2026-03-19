package sonarqube

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSonarqubeNewCodePeriodsImportUnit(t *testing.T) {
	rnd := generateRandomResourceName()
	
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Test global import - use Community-compatible PREVIOUS_VERSION
			{
				Config: testAccSonarqubeNewCodePeriodsGlobalPreviousVersion(rnd + "global"),
			},
			{
				ResourceName:      "sonarqube_new_code_periods." + rnd + "global",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "newCodePeriod",
			},
			// Test project import  
			{
				Config: testAccSonarqubeNewCodePeriodsProjectPreviousVersion(rnd + "project"),
			},
			{
				ResourceName:      "sonarqube_new_code_periods." + rnd + "project",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("newCodePeriod/%s", rnd + "project"),
			},
			// Test another project import for variety
			{
				Config: testAccSonarqubeNewCodePeriodsProjectPreviousVersion(rnd + "project2"),
			},
			{
				ResourceName:      "sonarqube_new_code_periods." + rnd + "project2",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("newCodePeriod/%s", rnd + "project2"),
			},
			// Test invalid import format - use Community-compatible PREVIOUS_VERSION
			{
				Config: testAccSonarqubeNewCodePeriodsGlobalPreviousVersion(rnd + "invalid"),
			},
			{
				ResourceName:  "sonarqube_new_code_periods." + rnd + "invalid",
				ImportState:   true,
				ImportStateId: "invalidFormat",
				ExpectError:   regexp.MustCompile("invalid import ID format"),
			},
		},
	})
}
