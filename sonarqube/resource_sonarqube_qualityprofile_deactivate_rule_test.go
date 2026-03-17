package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("sonarqube_qualityprofile_deactivate_rule", &resource.Sweeper{
		Name: "sonarqube_qualityprofile_deactivate_rule",
		F:    testSweeepSonarqibeQualityprofileDeactivateRuleSweeper,
	})
}

func testSweeepSonarqibeQualityprofileDeactivateRuleSweeper(r string) error {
	return nil
}

func testAccSonarqubeQualityprofileDeactivateRuleBasicConfig(rnd string, name string, key string, severity string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualityprofile" "%[1]s" {
			name     = "%[2]s"
			language = "xml"
		}

		resource "sonarqube_rule" "%[1]s" {
			custom_key = "%[3]s"
			markdown_description = "My rule"
			name = "%[3]s"
			severity = "%[4]s"
			template_key = "xml:XPathCheck"
			type = "VULNERABILITY"
		}

		resource "sonarqube_qualityprofile_activate_rule" "%[1]s" {
			key = sonarqube_qualityprofile.%[1]s.key
			rule = sonarqube_rule.%[1]s.id
			severity = "%[4]s"
		}

		resource "sonarqube_qualityprofile_deactivate_rule" "%[1]s" {
			key = sonarqube_qualityprofile.%[1]s.key
			rule = sonarqube_rule.%[1]s.id
			depends_on = [sonarqube_qualityprofile_activate_rule.%[1]s]
		}`, rnd, name, key, severity)
}

func TestAccSonarqubeQualityprofileDeactivateRuleBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualityprofile_deactivate_rule." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualityprofileDeactivateRuleBasicConfig(rnd, "testProfileDeactivate", "deactivateRule", "BLOCKER"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "key"),
					resource.TestCheckResourceAttrSet(name, "rule"),
					testAccCheckRuleNotActive(name),
				),
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key", "rule"},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(name, "key"),
					resource.TestCheckResourceAttrSet(name, "rule"),
					testAccCheckRuleNotActive(name),
				),
			},
		},
	})
}

func testAccCheckRuleNotActive(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		profileKey := rs.Primary.Attributes["key"]
		ruleKey := rs.Primary.Attributes["rule"]

		providerConfig := testAccProvider.Meta().(*ProviderConfiguration)
		sonarQubeURL := providerConfig.sonarQubeURL
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/rules/show"
		sonarQubeURL.RawQuery = url.Values{
			"key":     []string{ruleKey},
			"actives": []string{"true"},
		}.Encode()

		resp, err := httpRequestHelper(
			providerConfig.httpClient,
			"GET",
			sonarQubeURL.String(),
			http.StatusOK,
			"testAccCheckRuleNotActive",
		)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		activeRuleReadResponse := GetActiveRules{}
		if err := json.NewDecoder(resp.Body).Decode(&activeRuleReadResponse); err != nil {
			return err
		}

		for _, active := range activeRuleReadResponse.Actives {
			if active.QProfile == profileKey {
				return fmt.Errorf("rule %s still active on profile %s", ruleKey, profileKey)
			}
		}

		return nil
	}
}
