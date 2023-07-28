package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				),
			},
			{
				Config: testAccSonarqubeWebhookBasicConfig(rnd, secondName, secondUrl, secondSecret),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", secondName),
					resource.TestCheckResourceAttr(resourceName, "url", secondUrl),
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
