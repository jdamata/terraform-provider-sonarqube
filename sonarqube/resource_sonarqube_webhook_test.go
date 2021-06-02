package sonarqube

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"testing"

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

func testAccSonarqubeWebhookBasicConfig(rnd, name, url, secret string) string {
	return fmt.Sprintf(`
resource "sonarqube_webhook" "%s" {
	name   = "%s"
	url    = "%s"
	secret = "%s"
}
`, rnd, name, url, secret)
}
