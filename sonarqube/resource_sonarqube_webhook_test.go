package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSonarqubeWebhookBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	resourceName := "sonarqube_webhook." + rnd

	name := acctest.RandString(10)
	url := fmt.Sprintf("https://%s.com", acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	secret := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeWebhookBasicConfig(rnd, name, url, secret),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "url", url),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Version 10.1 of sonarqube does not return the secret in the api response anymore. Field 'secret' replaced by flag 'hasSecret' in response
				ImportStateVerifyIgnore: []string{"secret"},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "url", url),
				),
			},
		},
	})
}

func TestAccSonarqubeWebhookUpdate(t *testing.T) {
	rnd := generateRandomResourceName()
	resourceName := "sonarqube_webhook." + rnd

	firstName := acctest.RandString(10)
	firstUrl := fmt.Sprintf("https://%s.com", acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	firstSecret := acctest.RandString(10)

	secondName := acctest.RandString(10)
	secondUrl := fmt.Sprintf("https://%s.com", acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	secondSecret := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeWebhookBasicConfig(rnd, firstName, firstUrl, firstSecret),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", firstName),
					resource.TestCheckResourceAttr(resourceName, "url", firstUrl),
					resource.TestCheckResourceAttr(resourceName, "secret", firstSecret),
				),
			},
			{
				Config: testAccSonarqubeWebhookBasicConfig(rnd, secondName, secondUrl, secondSecret),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", secondName),
					resource.TestCheckResourceAttr(resourceName, "url", secondUrl),
					resource.TestCheckResourceAttr(resourceName, "secret", secondSecret),
				),
			},
		},
	})
}

func testAccSonarqubeWebhookBasicConfig(rnd, name, url, secret string) string {
	return fmt.Sprintf(`
resource "sonarqube_webhook" "%s" {
	name   = "%s"
	url    = "%s"
	secret = "%s"
}
`, rnd, name, url, secret)
}

func TestAccSonarqubeWebhookProjectBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	resourceName := "sonarqube_webhook." + rnd

	name := acctest.RandString(10)
	url := fmt.Sprintf("https://%s.com", acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	project := "testAccSonarqubeWebhookProject"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeWebhookProjectBasicConfig(rnd, name, url, project),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(resourceName, "project", project),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateIdFunc:       testAccSonarqubeWebhookProjectImportID(resourceName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret"},
			},
		},
	})
}

func testAccSonarqubeWebhookProjectImportID(resourceNode string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceNode]
		if !ok {
			return "", fmt.Errorf("Resource node not found: %s", resourceNode)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["id"], rs.Primary.Attributes["project"]), nil
	}
}

func testAccSonarqubeWebhookProjectBasicConfig(rnd string, name string, url string, project string) string {
	return fmt.Sprintf(`
		resource "sonarqube_project" "%[1]s" {
			name       = "%[4]s"
			project    = "%[4]s"
			visibility = "public" 
		}

		resource "sonarqube_webhook" "%[1]s" {
			name    = "%[2]s"
			url     = "%[3]s"
			project = sonarqube_project.%[1]s.project
		}`, rnd, name, url, project)
}
