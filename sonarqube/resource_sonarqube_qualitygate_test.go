package sonarqube

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		}`, rnd, name, is_default)
}

func TestAccSonarqubeQualitygateBasic(t *testing.T) {

	numDefaultConditions := 6 // 9.9 and above automatically add conditions to the quality gate

	sonarQubeVersion, err := getSonarQubeVersion()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	minimumVersion, _ := version.NewVersion("9.9")
	if sonarQubeVersion.LessThan(minimumVersion) {
		numDefaultConditions = 0
	}
	t.Logf("Sonar version: %s - %d default quality gate conditions expected", sonarQubeVersion.Core().String(), numDefaultConditions)

	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygateBasicConfig(rnd, "testAccSonarqubeQualitygate", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeQualitygate"),
					resource.TestCheckResourceAttr(name, "is_default", "true"),
					resource.TestCheckResourceAttr(name, "conditions.#", strconv.Itoa(numDefaultConditions)),
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

func testAccSonarqubeQualitygateCopyConfig(rnd string, baseName string, conditionName string, threshold string, op string, copyName string) string {
	return fmt.Sprintf(`
	resource "sonarqube_qualitygate" "%[2]s" {
		name = "%[2]s"
	}

	resource "sonarqube_qualitygate_condition" "qualitygate_condition" {
		gatename  = sonarqube_qualitygate.%[2]s.id
		metric    = "%[3]s"
		threshold = "%[4]s"
		op        = "%[5]s"
	}

	resource "sonarqube_qualitygate" "%[1]s" {
		depends_on = [sonarqube_qualitygate.%[2]s, sonarqube_qualitygate_condition.qualitygate_condition]
		name = "%[6]s"
		copy_from = "%[2]s"
	}`, rnd, baseName, conditionName, threshold, op, copyName)
}

func TestAccSonarqubeQualitygateCopy(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualitygate." + rnd

	// Copy test variables
	baseGateName := "baseGate"
	baseGateResourceName := "sonarqube_qualitygate." + baseGateName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualitygateCopyConfig(rnd, baseGateName, "comment_lines_density", "68", "LT", baseGateName+"Copy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(baseGateResourceName, "conditions", name, "conditions"),
				),
			},
		},
	})
}
