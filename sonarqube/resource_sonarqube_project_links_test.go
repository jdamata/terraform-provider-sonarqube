package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccSonarqubeProjectLinksConfig(rnd, projectKey, linkName, linkURL string) string {
	return fmt.Sprintf(`
resource "sonarqube_project" "%[1]s" {
  name       = "%[2]s"
  project    = "%[2]s"
  visibility = "public"
}

resource "sonarqube_project_links" "%[1]s" {
  project_key = sonarqube_project.%[1]s.project
  name        = "%[3]s"
  url         = "%[4]s"
}
`, rnd, projectKey, linkName, linkURL)
}

func TestAccSonarqubeProjectLinksBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	resourceName := "sonarqube_project_links." + rnd
	projectKey := "testAccProjectLinks" + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Avoid built-in reserved names ("Homepage", "CI", "SCM") — SonarQube assigns
				// a non-deletable built-in type to those links, causing teardown to fail.
				Config: testAccSonarqubeProjectLinksConfig(rnd, projectKey, "Documentation", "https://example.com"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_key", projectKey),
					resource.TestCheckResourceAttr(resourceName, "name", "Documentation"),
					resource.TestCheckResourceAttr(resourceName, "url", "https://example.com"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccSonarqubeProjectLinksImportID(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSonarqubeProjectLinksMultiple(t *testing.T) {
	rnd := generateRandomResourceName()
	projectKey := "testAccProjectLinksMulti" + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Avoid "Homepage" and "CI" — those map to non-deletable built-in link types.
				Config: fmt.Sprintf(`
resource "sonarqube_project" "%[1]s" {
  name       = "%[2]s"
  project    = "%[2]s"
  visibility = "public"
}
resource "sonarqube_project_links" "%[1]s_docs" {
  project_key = sonarqube_project.%[1]s.project
  name        = "Documentation"
  url         = "https://example.com"
}
resource "sonarqube_project_links" "%[1]s_build" {
  project_key = sonarqube_project.%[1]s.project
  name        = "Build"
  url         = "https://build.example.com"
}
`, rnd, projectKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarqube_project_links."+rnd+"_docs", "name", "Documentation"),
					resource.TestCheckResourceAttr("sonarqube_project_links."+rnd+"_build", "name", "Build"),
				),
			},
		},
	})
}

func testAccSonarqubeProjectLinksImportID(resourceNode string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceNode]
		if !ok {
			return "", fmt.Errorf("Resource node not found: %s", resourceNode)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.ID, rs.Primary.Attributes["project_key"]), nil
	}
}
