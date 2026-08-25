package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccSonarqubePermissionTemplateDefaultConfig(rnd, templateName, qualifier string) string {
	return fmt.Sprintf(`
resource "sonarqube_permission_template" "%[1]s" {
  name        = "%[2]s"
  description = "Test template for default assignment"
}

resource "sonarqube_permission_template_default" "%[1]s" {
  template_id = sonarqube_permission_template.%[1]s.id
  qualifier   = "%[3]s"
}
`, rnd, templateName, qualifier)
}

func testAccSonarqubePermissionTemplateDefaultUpdateConfig(rnd, templateName, templateNameUpdated, qualifier string) string {
	return fmt.Sprintf(`
resource "sonarqube_permission_template" "%[1]s" {
  name        = "%[2]s"
  description = "Test template for default assignment"
}

resource "sonarqube_permission_template" "%[1]s_b" {
  name        = "%[3]s"
  description = "Updated test template for default assignment"
}

resource "sonarqube_permission_template_default" "%[1]s" {
  template_id = sonarqube_permission_template.%[1]s_b.id
  qualifier   = "%[4]s"
}
`, rnd, templateName, templateNameUpdated, qualifier)
}

func TestAccSonarqubePermissionTemplateDefaultBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	resourceName := "sonarqube_permission_template_default." + rnd
	templateName := "testAccPermTemplateDefault" + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePermissionTemplateDefaultConfig(rnd, templateName, "TRK"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "qualifier", "TRK"),
					resource.TestCheckResourceAttrSet(resourceName, "template_id"),
				),
			},
			// Update: switch to a different template (in-place, no ForceNew on template_id)
			{
				Config: testAccSonarqubePermissionTemplateDefaultUpdateConfig(rnd, templateName, templateName+"B", "TRK"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "qualifier", "TRK"),
					resource.TestCheckResourceAttrSet(resourceName, "template_id"),
				),
			},
			// Import by qualifier
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "TRK",
				ImportStateVerify: true,
			},
		},
	})
}
