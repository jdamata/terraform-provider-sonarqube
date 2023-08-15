package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_setting", &resource.Sweeper{
		Name: "sonarqube_setting",
		F:    testSweepSonarqubeSettingSweeper,
	})
}

func testSweepSonarqubeSettingSweeper(r string) error {
	return nil
}

func testAccSonarqubeSettingBasicConfig(rnd string, key string, value string) string {
	return fmt.Sprintf(`
		resource "sonarqube_setting" "%[1]s" {
			key = "%[2]s"
			value = "%[3]s"
		}`, rnd, key, value)
}

func testAccSonarqubeSettingConfigMultipleValues(rnd string, key string, values []string) string {
	formattedValues := generateHCLList(values)
	return fmt.Sprintf(`
		resource "sonarqube_setting" "%[1]s" {
			key = "%[2]s"
			values = %[3]s
		}`, rnd, key, formattedValues)
}
func testAccSonarqubeSettingConfigMultipleFields(rnd string, key string, fields map[string]string) string {
	formattedFields := generateHCLMap(fields)
	return fmt.Sprintf(`
		resource "sonarqube_setting" "%[1]s" {
			key = "%[2]s"
			field_values = [%[3]s]
		}`, rnd, key, formattedFields)
}
func TestAccSonarqubeSettingBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_setting." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeSettingBasicConfig(rnd, "sonar.demo", "sonarqube@example.org"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "sonar.demo"),
					resource.TestCheckResourceAttr(name, "value", "sonarqube@example.org"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", "sonar.demo"),
					resource.TestCheckResourceAttr(name, "value", "sonarqube@example.org"),
				),
			},
		},
	})
}
func TestAccSonarqubeSettingMultipleValues(t *testing.T) {
	key := "sonar.global.exclusions" // Needs to be a setting that accepts multiple values
	rnd := generateRandomResourceName()
	name := "sonarqube_setting." + rnd
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeSettingConfigMultipleValues(rnd, key, []string{"foo", "bar"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", key),
					resource.TestCheckTypeSetElemAttr(name, "values.*", "foo"),
					resource.TestCheckTypeSetElemAttr(name, "values.*", "bar"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", key),
					resource.TestCheckTypeSetElemAttr(name, "values.*", "foo"),
					resource.TestCheckTypeSetElemAttr(name, "values.*", "bar"),
				),
			},
		},
	})
}
func TestAccSonarqubeSettingMultipleFields(t *testing.T) {
	key := "sonar.issue.ignore.multicriteria" // Needs to be a setting that accepts multiple fields
	rnd := generateRandomResourceName()
	name := "sonarqube_setting." + rnd
	expected := map[string]string{"ruleKey": "foo", "resourceKey": "bar"}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeSettingConfigMultipleFields(rnd, key, expected),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", key),
					resource.TestCheckTypeSetElemNestedAttrs(name, "field_values.*", expected),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "key", key),
					resource.TestCheckTypeSetElemNestedAttrs(name, "field_values.*", expected),
				),
			},
		},
	})
}
