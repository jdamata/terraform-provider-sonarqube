package sonarqube

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
func testAccPreCheckPortfolioSupport(t *testing.T) {
	if err := checkPortfolioSupport(testAccProvider.Meta().(*ProviderConfiguration)); err != nil {
		t.Skipf("Skipping test of unsupported feature (Portfolio)")
	}
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

func testAccSonarqubePortfolioConfigSelectionModeTags(rnd string, key string, name string, description string, visibility string, selectionMode string, tags []string) string {
	formattedTags := generateHCLList(tags)
	return fmt.Sprintf(`
		resource "sonarqube_portfolio" "%[1]s" {
		  key       = "%[2]s"
		  name    = "%[3]s"
		  description = "%[4]s"
		  visibility = "%[5]s"
		  selection_mode = "%[6]s"
		  tags = %[7]s // Note that the "" should be missing since this is a list
		}
		`, rnd, key, name, description, visibility, selectionMode, formattedTags)
}

func testAccSonarqubePortfolioConfigSelectionModeRegex(rnd string, key string, name string, description string, visibility string, selectionMode string, regexp string) string {
	return fmt.Sprintf(`
		resource "sonarqube_portfolio" "%[1]s" {
		  key       = "%[2]s"
		  name    = "%[3]s"
		  description = "%[4]s"
		  visibility = "%[5]s"
		  selection_mode = "%[6]s"
		  regexp = "%[7]s"

		}
		`, rnd, key, name, description, visibility, selectionMode, regexp)
}

func TestAccSonarqubePortfolioBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_portfolio." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckPortfolioSupport(t) },
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
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckPortfolioSupport(t) },
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
		PreCheck: func() { testAccPreCheck(t); testAccPreCheckPortfolioSupport(t) },

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
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckPortfolioSupport(t) },
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
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckPortfolioSupport(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccSonarqubePortfolioBasicConfig(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "testAccSonarqubePortfolioDescription", "badValue"),
				ExpectError: regexp.MustCompile("expected .* to be one of"),
			},
		},
	})
}

func TestAccSonarqubePortfolioSelectionModeError(t *testing.T) {
	rnd := generateRandomResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckPortfolioSupport(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccSonarqubePortfolioConfigSelectionMode(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "testAccSonarqubePortfolioDescription", "public", "badValue"),
				ExpectError: regexp.MustCompile("expected .* to be one of"),
			},
		},
	})
}

func TestAccSonarqubePortfolioSelectionModeNone(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_portfolio." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckPortfolioSupport(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePortfolioConfigSelectionMode(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "testAccSonarqubePortfolioDescription", "public", "NONE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "selection_mode", "NONE"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "selection_mode", "NONE"),
				),
			},
		},
	})
}

func TestAccSonarqubePortfolioSelectionModeManual(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_portfolio." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckPortfolioSupport(t) },
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

func TestAccSonarqubePortfolioSelectionModeTags(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_portfolio." + rnd
	tags := []string{"tag1", "tag2"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckPortfolioSupport(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePortfolioConfigSelectionModeTags(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "testAccSonarqubePortfolioDescription", "public", "TAGS", tags),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "selection_mode", "TAGS"),
					resource.TestCheckResourceAttr(name, "tags.0", tags[0]),
					resource.TestCheckResourceAttr(name, "tags.1", tags[1]),
					resource.TestCheckNoResourceAttr(name, "branch"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "selection_mode", "TAGS"),
					resource.TestCheckResourceAttr(name, "tags.0", tags[0]),
					resource.TestCheckResourceAttr(name, "tags.1", tags[1]),
					resource.TestCheckNoResourceAttr(name, "branch"),
				),
			},
		},
	})
}

func TestAccSonarqubePortfolioSelectionModeRegexp(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_portfolio." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckPortfolioSupport(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePortfolioConfigSelectionModeRegex(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "testAccSonarqubePortfolioDescription", "public", "REGEXP", "regexp1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "selection_mode", "REGEXP"),
					resource.TestCheckResourceAttr(name, "regexp", "regexp1"),
					resource.TestCheckNoResourceAttr(name, "branch"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "selection_mode", "REGEXP"),
					resource.TestCheckResourceAttr(name, "regexp", "regexp1"),
					resource.TestCheckNoResourceAttr(name, "branch"),
				),
			},
		},
	})
}

func TestAccSonarqubePortfolioSelectionModeUpdate(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_portfolio." + rnd
	tags := []string{"tag1", "tag2"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckPortfolioSupport(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubePortfolioBasicConfig(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "testAccSonarqubePortfolioDescription", "public"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
				),
			},
			{
				Config: testAccSonarqubePortfolioConfigSelectionModeTags(rnd, "testAccSonarqubePortfolioKey", "testAccSonarqubePortfolioName", "testAccSonarqubePortfolioDescription", "public", "TAGS", tags),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "selection_mode", "TAGS"),
					resource.TestCheckResourceAttr(name, "tags.0", tags[0]),
					resource.TestCheckResourceAttr(name, "tags.1", tags[1]),
					resource.TestCheckNoResourceAttr(name, "branch"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "testAccSonarqubePortfolioKey"),
					resource.TestCheckResourceAttr(name, "selection_mode", "TAGS"),
					resource.TestCheckResourceAttr(name, "tags.0", tags[0]),
					resource.TestCheckResourceAttr(name, "tags.1", tags[1]),
					resource.TestCheckNoResourceAttr(name, "branch"),
				),
			},
		},
	})
}
