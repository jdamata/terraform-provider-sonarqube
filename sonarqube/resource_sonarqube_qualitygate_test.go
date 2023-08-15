package sonarqube

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("sonarqube_qualitygate", &resource.Sweeper{
		Name: "sonarqube_qualitygate",
		F:    testSweepSonarqubeQualitygateSweeper,
	})
}

func testSweepSonarqubeQualitygateSweeper(r string) error {
	return nil
}

func testAccSonarqubeQualitygateBasicConfig(rnd string, name string, is_default string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualitygate" "%[1]s" {
			name = "%[2]s"
			is_default = "%[3]s"

			condition {
				metric    = "new_coverage"
				op        = "LT"
				threshold = "50"
			}
		}`, rnd, name, is_default)
}

// Add a basic quality gate
func TestAccSonarqubeQualitygateBasic(t *testing.T) {

	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate." + rnd
	expectedConditions := 1

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygateBasicConfig(rnd, "testAccSonarqubeQualitygate", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeQualitygate"),
					resource.TestCheckResourceAttr(name, "is_default", "true"),
					resource.TestCheckResourceAttr(name, "condition.#", strconv.Itoa(expectedConditions)),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_default"},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeQualitygate"),
					resource.TestCheckResourceAttr(name, "is_default", "true"),
				),
			},
		},
	})
}

// Add a gate with multiple conditions and check tha the values are as expected.
func testAccSonarqubeQualitygateWithConditionsConfig(rnd string, name string, is_default string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualitygate" "%[1]s" {
			name = "%[2]s"
			is_default = "%[3]s"

			condition {
				metric    = "new_coverage"
				op        = "LT"
				threshold = "50"
			}

			condition {
				metric    = "reliability_rating"
				op        = "GT"
				threshold = "2"
			}

		}`, rnd, name, is_default)
}

func TestAccSonarqubeQualitygateConditions(t *testing.T) {

	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate." + rnd
	expectedConditions := 2

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygateWithConditionsConfig(rnd, "TestAccSonarqubeQualitygateConditions", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "TestAccSonarqubeQualitygateConditions"),
					resource.TestCheckResourceAttr(name, "condition.#", strconv.Itoa(expectedConditions)),
					resource.TestCheckResourceAttr(name, "condition.0.metric", "new_coverage"),
					resource.TestCheckResourceAttr(name, "condition.0.op", "LT"),
					resource.TestCheckResourceAttr(name, "condition.0.threshold", "50"),
					resource.TestCheckResourceAttr(name, "condition.1.metric", "reliability_rating"),
					resource.TestCheckResourceAttr(name, "condition.1.op", "GT"),
					resource.TestCheckResourceAttr(name, "condition.1.threshold", "2"),
				),
			},
		},
	})
}

func testAccSonarqubeQualitygateChangeDefaultConfig(rnd string, name string, firstIsDefault bool, threshold2 string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualitygate" "%[1]s-1" {
			name = "%[2]s-1"
			is_default = "%[3]s"

			condition {
				metric    = "new_coverage"
				op        = "LT"
				threshold = "50"
			}
		}

		resource "sonarqube_qualitygate" "%[1]s-2" {
			name = "%[2]s-2"
			is_default = "%[4]s"

			condition {
				metric    = "new_coverage"
				op        = "LT"
				threshold = "%[5]s"
			}
		}`, rnd, name, strconv.FormatBool(firstIsDefault), strconv.FormatBool(!firstIsDefault), threshold2)
}

// Add two quality gates and then change which one is the default
func TestAccSonarqubeQualitygateChangeDefault(t *testing.T) {

	rnd := generateRandomResourceName()
	baseName := "sonarqube_qualitygate." + rnd
	gate1 := baseName + "-1"
	gate2 := baseName + "-2"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygateChangeDefaultConfig(rnd, "TestAccSonarqubeQualitygateChangeDefaultName", true, "10"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(gate1, "name", "TestAccSonarqubeQualitygateChangeDefaultName-1"),
					resource.TestCheckResourceAttr(gate1, "is_default", "true"),
					resource.TestCheckResourceAttr(gate2, "name", "TestAccSonarqubeQualitygateChangeDefaultName-2"),
					resource.TestCheckResourceAttr(gate2, "is_default", "false"),
				),
			},
			// Now Make the second one the default
			{
				Config: testAccSonarqubeQualitygateChangeDefaultConfig(rnd, "TestAccSonarqubeQualitygateChangeDefaultName", false, "10"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(gate1, "is_default", "false"),
					resource.TestCheckResourceAttr(gate2, "is_default", "true"),
				),
			},
			// Change a property of the second one and check it is still the default (and the threshold changed)
			{
				Config: testAccSonarqubeQualitygateChangeDefaultConfig(rnd, "TestAccSonarqubeQualitygateChangeDefaultName", false, "20"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(gate1, "is_default", "false"),
					resource.TestCheckResourceAttr(gate2, "is_default", "true"),
					resource.TestCheckResourceAttr(gate2, "condition.0.threshold", "20"),
				),
			},
		},
	})
}

func testAccSonarqubeQualitygateDeleteDefaultConfig(rnd string, name string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualitygate" "%[1]s" {
			name = "%[2]s"
			is_default = "true"

			condition {
				metric    = "new_coverage"
				op        = "LT"
				threshold = "50"
			}
		}`, rnd, name)
}

// Deleting a quality gate that was set as the default should set the built in Sonar way as the default
func TestAccSonarqubeQualitygateDeleteDefault(t *testing.T) {

	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygateDeleteDefaultConfig(rnd, "TestAccSonarqubeQualitygateDeleteDefault"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "TestAccSonarqubeQualitygateDeleteDefault"),
					resource.TestCheckResourceAttr(name, "is_default", "true"),
				),
			},
			// Now remove the resource
			{
				Config:  testAccSonarqubeQualitygateDeleteDefaultConfig(rnd, "TestAccSonarqubeQualitygateDeleteDefault"),
				Destroy: true,
				Check: resource.ComposeTestCheckFunc(
					//					resource.TestCheckNoResourceAttr("sonarqube_qualitygate", rnd),
					checkSonarWayIsDefault,
				),
			},
		},
	})
}

func testAccSonarqubeQualitygateCopyConfig(rnd string, baseName string, conditionName string, threshold string, op string, copyName string) string {
	return fmt.Sprintf(`
	resource "sonarqube_qualitygate" "%[2]s" {
		name = "%[2]s"

		condition {
			metric    = "%[3]s"
			threshold = "%[4]s"
			op        = "%[5]s"
		}
	}

	resource "sonarqube_qualitygate" "%[1]s" {
		depends_on = [sonarqube_qualitygate.%[2]s]
		name = "%[6]s"
		copy_from = "%[2]s"
	}`, rnd, baseName, conditionName, threshold, op, copyName)
}

// Copy a quality gate and validate the copy exists and it does not have any explicit conditions
func TestAccSonarqubeQualitygateCopy(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate." + rnd

	// Copy test variables
	baseGateName := "baseGate"
	baseGateResourceName := "sonarqube_qualitygate." + baseGateName
	copyName := baseGateName + "-copy"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygateCopyConfig(rnd, baseGateName, "comment_lines_density", "68", "LT", copyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", copyName),
					resource.TestCheckResourceAttrPair(baseGateResourceName, "condition", name, "condition"),
					resource.TestCheckResourceAttr(name, "condition.#", "0"),
				),
			},
		},
	})
}

func checkSonarWayIsDefault(s *terraform.State) error {

	sonarQubeURL := fmt.Sprintf("%[1]s/api/qualitygates/show?name=Sonar%%20way", strings.TrimSuffix(os.Getenv("SONAR_HOST"), "/"))
	auth := os.Getenv("SONAR_USER") + ":" + os.Getenv("SONAR_PASS")
	b64 := base64.StdEncoding.EncodeToString([]byte(auth))
	authHeader := "Basic " + b64

	client := retryablehttp.NewClient()
	req, _ := retryablehttp.NewRequest("GET", sonarQubeURL, http.NoBody)
	req.Header.Set("Authorization", authHeader)

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("checkSonarWayIsDefault: failed to execute http request: %v. Request: %v", err, req)
	}

	// Check response code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("checkSonarWayIsDefault: Unexpected statusCode: %v", resp.StatusCode)
	}
	defer resp.Body.Close()

	qualityGateReadResponse := GetQualityGate{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateReadResponse)
	if err != nil {
		return fmt.Errorf("checkSonarWayIsDefault: Failed to decode json into struct: %+v", err)
	}
	if qualityGateReadResponse.Actions.SetAsDefault {
		return fmt.Errorf("checkSonarWayIsDefault: Sonar way is not the default Quality Profile: %+v", err)
	}
	return nil
}
