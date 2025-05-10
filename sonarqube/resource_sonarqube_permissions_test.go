package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

func testAccSonarqubePermissionGroupNameConfig(id string, groupName string, permissions []string) string {
	formattedPermissions := generateHCLList(permissions)
	return fmt.Sprintf(`
		resource "sonarqube_group" "%[1]s" {
			name        = "%[2]s"
			description = "%[2]s"
		}

		resource "sonarqube_permissions" "%[1]s" {
			group_name  = sonarqube_group.%[1]s.name
			permissions = %[3]s
		}`, id, groupName, formattedPermissions)
}

func TestAccSonarqubePermissionGroupName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_permissions." + rnd
	permissions := []string{"admin"}
	updatedPermissions := []string{"admin", "profileadmin"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePermissionGroupNameConfig(rnd, "testAccSonarqubePermissions", permissions),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "group_name", "testAccSonarqubePermissions"),
				),
			},
			{
				Config: testAccSonarqubePermissionGroupNameConfig(rnd, "testAccSonarqubePermissions", updatedPermissions),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "group_name", "testAccSonarqubePermissions"),
				),
			},
		},
	})
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
		  	login_name  = sonarqube_user.%[1]s.name
		  	permissions = %[3]s
		}
		`, id, loginName, formattedPermissions)
}

func TestAccSonarqubePermissionLoginName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_permissions." + rnd
	permissions := []string{"gateadmin", "profileadmin"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePermissionLoginNameConfig(rnd, "testAccSonarqubePermissions", permissions),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "login_name", "testAccSonarqubePermissions"),
					//resource.TestCheckResourceAttr(name, "permissions", generateHCLList(permissions)),
				),
			},
		},
	})
}

func testAccSonarqubePermissionLoginNameTemplateNameConfig(id string, loginName string, permissions []string) string {
	formattedPermissions := generateHCLList(permissions)
	return fmt.Sprintf(`
		resource "sonarqube_user" "%[1]s" {
			login_name = "%[2]s"
			name       = "%[2]s"
			password   = "secret-sauce37!"
		}

		resource "sonarqube_permission_template" "this" {
			name = "foo"
		}
		resource "sonarqube_permissions" "%[1]s" {
		  	login_name    = sonarqube_user.%[1]s.name
			template_name = sonarqube_permission_template.this.name
		  	permissions   = %[3]s
		}
		`, id, loginName, formattedPermissions)
}

func TestAccSonarqubePermissionLoginNameTemplateName(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_permissions." + rnd
	permissions := []string{"codeviewer", "scan"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePermissionLoginNameTemplateNameConfig(rnd, "testAccSonarqubePermissions", permissions),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "login_name", "testAccSonarqubePermissions"),
				),
			},
		},
	})
}

func TestAccSonarqubePermissionUpdate(t *testing.T) {
	rnd := generateRandomResourceName()
	resourceName := "sonarqube_permissions." + rnd
	username := "update-test-user"
	initialPermissions := []string{"admin", "scan"}
	updatedPermissions := []string{"admin", "provisioning", "applicationcreator"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePermissionLoginNameConfig(rnd, username, initialPermissions),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "login_name", username),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", fmt.Sprintf("%d", len(initialPermissions))),
					resource.TestCheckTypeSetElemAttr(resourceName, "permissions.*", initialPermissions[0]),
					resource.TestCheckTypeSetElemAttr(resourceName, "permissions.*", initialPermissions[1]),
				),
			},
			{
				Config: testAccSonarqubePermissionLoginNameConfig(rnd, username, updatedPermissions),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "login_name", username),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", fmt.Sprintf("%d", len(updatedPermissions))),
					resource.TestCheckTypeSetElemAttr(resourceName, "permissions.*", updatedPermissions[0]),
					resource.TestCheckTypeSetElemAttr(resourceName, "permissions.*", updatedPermissions[1]),
					resource.TestCheckTypeSetElemAttr(resourceName, "permissions.*", updatedPermissions[2]),
				),
			},
		},
	})
}

func TestAccSonarqubePermissionGroupUpdate(t *testing.T) {
	rnd := generateRandomResourceName()
	resourceName := "sonarqube_permissions." + rnd
	groupName := "update-test-group"
	initialPermissions := []string{"profileadmin", "gateadmin"}
	updatedPermissions := []string{"profileadmin"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create with initial permissions
			{
				Config: testAccSonarqubePermissionGroupNameConfig(rnd, groupName, initialPermissions),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", groupName),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", fmt.Sprintf("%d", len(initialPermissions))),
					resource.TestCheckTypeSetElemAttr(resourceName, "permissions.*", initialPermissions[0]),
					resource.TestCheckTypeSetElemAttr(resourceName, "permissions.*", initialPermissions[1]),
				),
			},
			// Update by removing a permission
			{
				Config: testAccSonarqubePermissionGroupNameConfig(rnd, groupName, updatedPermissions),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", groupName),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", fmt.Sprintf("%d", len(updatedPermissions))),
					resource.TestCheckTypeSetElemAttr(resourceName, "permissions.*", updatedPermissions[0]),
				),
			},
		},
	})
}
