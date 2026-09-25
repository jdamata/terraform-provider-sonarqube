package sonarqube

import (
	"encoding/json"
	"fmt"
	version "github.com/hashicorp/go-version"
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

func testAccSonarqubeQualityprofileDeactivateRuleInheritedConfig(rnd string, withDeactivate bool) string {
	deactivate := ""
	if withDeactivate {
		deactivate = fmt.Sprintf(`
		resource "sonarqube_qualityprofile_deactivate_rule" "%[1]s" {
			key  = sonarqube_qualityprofile.%[1]s_child.key
			rule = sonarqube_rule.%[1]s.id
		}`, rnd)
	}
	return fmt.Sprintf(`
		resource "sonarqube_qualityprofile" "%[1]s_parent" {
			name     = "%[1]s_parent"
			language = "xml"
		}

		resource "sonarqube_rule" "%[1]s" {
			custom_key           = "%[1]s"
			markdown_description = "My rule"
			name                 = "%[1]s"
			severity             = "MAJOR"
			template_key         = "xml:XPathCheck"
			type                 = "VULNERABILITY"
		}

		resource "sonarqube_qualityprofile_activate_rule" "%[1]s" {
			key      = sonarqube_qualityprofile.%[1]s_parent.key
			rule     = sonarqube_rule.%[1]s.id
			severity = "MAJOR"
		}

		resource "sonarqube_qualityprofile" "%[1]s_child" {
			name     = "%[1]s_child"
			language = "xml"
			parent   = sonarqube_qualityprofile.%[1]s_parent.name

			depends_on = [sonarqube_qualityprofile_activate_rule.%[1]s]
		}
		%[2]s`, rnd, deactivate)
}

// Deactivating an inherited rule and then removing the deactivation must leave the rule
// active on the child again. The second step exercises Delete.
func TestAccSonarqubeQualityprofileDeactivateRuleInheritedReactivate(t *testing.T) {
	rnd := generateRandomResourceName()
	deactivateName := "sonarqube_qualityprofile_deactivate_rule." + rnd
	childName := "sonarqube_qualityprofile." + rnd + "_child"
	ruleName := "sonarqube_rule." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckInheritedRuleDeactivation(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualityprofileDeactivateRuleInheritedConfig(rnd, true),
				Check:  testAccCheckRuleNotActive(deactivateName),
			},
			{
				Config: testAccSonarqubeQualityprofileDeactivateRuleInheritedConfig(rnd, false),
				Check:  testAccCheckRuleActiveOnProfile(childName, ruleName),
			},
		},
	})
}

func testAccCheckRuleActiveOnProfile(profileResource string, ruleResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		profile, ok := s.RootModule().Resources[profileResource]
		if !ok {
			return fmt.Errorf("not found: %s", profileResource)
		}
		rule, ok := s.RootModule().Resources[ruleResource]
		if !ok {
			return fmt.Errorf("not found: %s", ruleResource)
		}
		profileKey := profile.Primary.Attributes["key"]
		ruleKey := rule.Primary.ID

		providerConfig := testAccProvider.Meta().(*ProviderConfiguration)
		sonarQubeURL := providerConfig.sonarQubeURL
		sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/rules/show"
		sonarQubeURL.RawQuery = url.Values{
			"key":     []string{ruleKey},
			"actives": []string{"true"},
		}.Encode()

		resp, err := httpRequestHelper(providerConfig.httpClient, "GET", sonarQubeURL.String(), http.StatusOK, "testAccCheckRuleActiveOnProfile")
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
				return nil
			}
		}
		return fmt.Errorf("rule %s not active on profile %s after the deactivation was removed", ruleKey, profileKey)
	}
}

// SonarQube 9.9 LTS rejects deactivating a rule inherited from the parent profile
// ("Cannot deactivate inherited rule"); the capability exists from 10.x onwards.
func testAccPreCheckInheritedRuleDeactivation(t *testing.T) {
	sonarQubeVersion := testAccProvider.Meta().(*ProviderConfiguration).sonarQubeVersion
	minimumVersion, _ := version.NewVersion("10.0")
	if sonarQubeVersion.LessThan(minimumVersion) {
		t.Skipf("inherited rules cannot be deactivated on SonarQube %s (requires 10.0+)", sonarQubeVersion)
	}
}
