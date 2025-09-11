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
			// Test global import
			{
				Config: testAccSonarqubeNewCodePeriodsGlobalNumberOfDays(rnd + "global"),
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
			// Test branch import
			{
				Config: testAccSonarqubeNewCodePeriodsBranchReferenceBranch(rnd + "branch"),
			},
			{
				ResourceName:      "sonarqube_new_code_periods." + rnd + "branch",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("newCodePeriod/main/%s", rnd + "branch"),
			},
			// Test invalid import format
			{
				Config: testAccSonarqubeNewCodePeriodsGlobalNumberOfDays(rnd + "invalid"),
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
