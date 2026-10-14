package sonarqube

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccPreCheckApplicationSupport(t *testing.T) {
	if testAccIsCommunityEdition(t) {
		t.Skip("applications are not available in the Community edition")
	}
}

func testAccSonarqubeApplicationConfig(rnd, name, description, visibility, project, tags string) string {
	return fmt.Sprintf(`
		resource "sonarqube_project" "%[1]s_a" {
			name       = "%[1]s-a"
			project    = "%[1]s-a"
			visibility = "public"
		}
		resource "sonarqube_project" "%[1]s_b" {
			name       = "%[1]s-b"
			project    = "%[1]s-b"
			visibility = "public"
		}
		resource "sonarqube_application" "%[1]s" {
			key         = "%[1]s-app"
			name        = "%[2]s"
			description = "%[3]s"
			visibility  = "%[4]s"
			projects    = [sonarqube_project.%[1]s_%[5]s.project]
			tags        = [%[6]s]
		}`, rnd, name, description, visibility, project, tags)
}

func TestAccSonarqubeApplicationBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_application." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckApplicationSupport(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeApplicationConfig(rnd, "First name", "first", "private", "a", `"team-a"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", rnd+"-app"),
					resource.TestCheckResourceAttr(name, "name", "First name"),
					resource.TestCheckResourceAttr(name, "description", "first"),
					resource.TestCheckResourceAttr(name, "visibility", "private"),
					resource.TestCheckResourceAttr(name, "projects.#", "1"),
					resource.TestCheckTypeSetElemAttr(name, "projects.*", rnd+"-a"),
					resource.TestCheckResourceAttr(name, "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr(name, "tags.*", "team-a"),
				),
			},
			{
				Config: testAccSonarqubeApplicationConfig(rnd, "Second name", "second", "public", "b", `"team-b", "prod"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "Second name"),
					resource.TestCheckResourceAttr(name, "description", "second"),
					resource.TestCheckResourceAttr(name, "visibility", "public"),
					resource.TestCheckResourceAttr(name, "projects.#", "1"),
					resource.TestCheckTypeSetElemAttr(name, "projects.*", rnd+"-b"),
					resource.TestCheckResourceAttr(name, "tags.#", "2"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSonarqubeApplicationDeletedOutOfBand(t *testing.T) {
	rnd := generateRandomResourceName()
	testAccDeletedOutOfBand(t, func() { testAccPreCheckApplicationSupport(t) },
		testAccSonarqubeApplicationConfig(rnd, "Out of band", "oob", "public", "a", `"oob"`),
		testAccPost(t, "/api/applications/delete", url.Values{"application": {rnd + "-app"}}))
}
