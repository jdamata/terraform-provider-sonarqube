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
	permissions := []string{"admin"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePermissionGroupNameConfig(rnd, "sonar-users", permissions),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "group_name", "sonar-users"),
					//resource.TestCheckResourceAttr(name, "permissions", generateHCLList(permissions)),
				),
			},
		},
	})
}

func TestAccSonarqubePermission_LoginName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_permissions." + rnd
	permissions := []string{"gateadmin", "profileadmin"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePermissionLoginNameConfig(rnd, "admin", permissions),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "login_name", "admin"),
					//resource.TestCheckResourceAttr(name, "permissions", generateHCLList(permissions)),
				),
			},
		},
	})
}

func testAccSonarqubePermissionGroupNameConfig(id string, groupName string, permissions []string) string {
	formattedPermissions := generateHCLList(permissions)
	return fmt.Sprintf(`
	  	resource "sonarqube_permissions" "%[1]s" {
			group_name  = "%[2]s"
			permissions = %[3]s
		}`, id, groupName, formattedPermissions)
}

func testAccSonarqubePermissionLoginNameConfig(id string, loginName string, permissions []string) string {
	formattedPermissions := generateHCLList(permissions)
	return fmt.Sprintf(`	  
		resource "sonarqube_permissions" "%[1]s" {
		  	login_name  = "%[2]s"
		  	permissions = %[3]s
		}
		`, id, loginName, formattedPermissions)
}
