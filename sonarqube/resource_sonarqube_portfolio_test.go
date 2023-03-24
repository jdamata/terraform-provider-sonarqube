package sonarqube

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_portfolio", &resource.Sweeper{
		Name: "sonarqube_portfolio",
		F:    testSweepSonarqubePortfolioSweeper,
	})
}

// TODO: implement sweeper to clean up portfolio: https://www.terraform.io/docs/extend/testing/acceptance-tests/sweepers.html
func testSweepSonarqubePortfolioSweeper(r string) error {
	return nil
}

func testAccSonarqubePortfolioBasicConfig(rnd string, key string, name string, description string, visibility string) string {
	return fmt.Sprintf(`
		resource "sonarqube_portfolio" "%[1]s" {
		  key       = "%[2]s"
		  name    = "%[3]s"
		  description = "%[4]s"
		  visibility = "%[5]s"
		}
		`, rnd, key, name, description, visibility)
}

func testAccSonarqubePortfolioConfigSelectionMode(rnd string, key string, name string, description string, visibility string, selectionMode string) string {
	return fmt.Sprintf(`
		resource "sonarqube_portfolio" "%[1]s" {
		  key       = "%[2]s"
		  name    = "%[3]s"
		  description = "%[4]s"
		  visibility = "%[5]s"
		  selection_mode = "%[6]s"
		}
		`, rnd, key, name, description, visibility, selectionMode)
}


func TestAccSonarqubePortfolioBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_portfolio." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePortfolioBasicConfig(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "testAccSonarqubePortfolioDescription", "public"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubePortfolioName"),
					resource.TestCheckResourceAttr(name, "qualifier", "VW"), // Qualifier for Portfolios seems to always be "VW" (views)
					resource.TestCheckResourceAttr(name, "description", "testAccSonarqubePortfolioDescription"),
					resource.TestCheckResourceAttr(name, "visibility", "public"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubePortfolioName"),
					resource.TestCheckResourceAttr(name, "description", "testAccSonarqubePortfolioDescription"),
					resource.TestCheckResourceAttr(name, "visibility", "public"),
				),
			},
		},
	})
}

func TestAccSonarqubePortfolioNameUpdate(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_portfolio." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePortfolioBasicConfig(rnd, "testAccSonarqubePortfolioKey", "oldName", "testAccSonarqubePortfolioDescription", "public"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "name", "oldName"),
				),
			},
			{
				Config: testAccSonarqubePortfolioBasicConfig(rnd, "testAccSonarqubePortfolioKey", "newName", "testAccSonarqubePortfolioDescription", "public"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "name", "newName"),
				),
			},
		},
	})
}

func TestAccSonarqubePortfolioDescriptionUpdate(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_portfolio." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePortfolioBasicConfig(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "oldDescription", "public"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "description", "oldDescription"),
				),
			},
			{
				Config: testAccSonarqubePortfolioBasicConfig(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "newDescription", "public"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "description", "newDescription"),
				),
			},
		},
	})
}

func TestAccSonarqubePortfolioVisibilityUpdate(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_portfolio." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePortfolioBasicConfig(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "testAccSonarqubePortfolioDescription", "public"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "visibility", "public"),
				),
			},
			{
				Config: testAccSonarqubePortfolioBasicConfig(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "testAccSonarqubePortfolioDescription", "private"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "visibility", "private"),
				),
			},
		},
	})
}


func TestAccSonarqubePortfolioVisibilityError(t *testing.T) {
	rnd := generateRandomResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePortfolioBasicConfig(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "testAccSonarqubePortfolioDescription", "badValue"),
				ExpectError: regexp.MustCompile("Accepted values are .* for key .* got:"),
			},
		},
	})
}

func TestAccSonarqubePortfolioSelectionModeError(t *testing.T) {
	rnd := generateRandomResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePortfolioConfigSelectionMode(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "testAccSonarqubePortfolioDescription", "public", "badValue"),
				ExpectError: regexp.MustCompile("Accepted values are .* for key .* got:"),
			},
		},
	})
}

func TestAccSonarqubePortfolioSelectionModeManual(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_portfolio." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePortfolioConfigSelectionMode(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "testAccSonarqubePortfolioDescription", "public", "MANUAL"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "selection_mode", "MANUAL"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "selection_mode", "MANUAL"),
				),
			},
		},
	})
}
