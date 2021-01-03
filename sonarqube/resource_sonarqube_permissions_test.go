package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_permissions", &resource.Sweeper{
		Name: "sonarqube_permissions",
		F:    testSweepPermissionSweeper,
	})
}

// TODO: implement sweeper to clean up permission_template: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepPermissionSweeper(r string) error {
	return nil
}

func TestAccSonarqubePermission_GroupName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_permissions." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePermissionGroupNameConfig(rnd, "permission-acc-test-group", []string{"admin"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "group_name", "permission-acc-test-group"),
					//resource.TestCheckResourceAttr(name, "permissions", generateHCLList([]string{"admin"})),
				),
			},
		},
	})
}

func TestAccSonarqubePermission_LoginName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_permissions." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePermissionLoginNameConfig(rnd, "permission-acc-test-user", []string{"gateadmin", "profileadmin"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "login_name", "permission-acc-test-user"),
					//resource.TestCheckResourceAttr(name, "permissions", generateHCLList([]string{"gateadmin", "profileadmin"})),
				),
			},
		},
	})
}

func testAccSonarqubePermissionGroupNameConfig(id string, groupName string, permissions []string) string {
	formattedPermissions := generateHCLList(permissions)
	return fmt.Sprintf(`
		resource "sonarqube_group" "%[1]s" {
	  	  	name        = "%[2]s"
	  	  	description = "example description"
		}

	  	resource "sonarqube_permissions" "%[1]s" {
			group_name  = "%[2]s"
			permissions = %[3]s
		}`, id, groupName, formattedPermissions)
}

func testAccSonarqubePermissionLoginNameConfig(id string, loginName string, permissions []string) string {
	formattedPermissions := generateHCLList(permissions)
	return fmt.Sprintf(`
	    resource "sonarqube_user" "%[1]s" {
			login_name = "%[2]s"
		  	name       = "%[2]s"
		  	password   = "secret-sauce37!"
		}
	  
		resource "sonarqube_permissions" "%[1]s" {
		  	login_name  = "%[2]s"
		  	permissions = %[3]s
		}
		`, id, loginName, formattedPermissions)
}
