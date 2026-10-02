package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func testAccPost(t *testing.T, path string, params url.Values) func() {
	return func() { testAccSonarqubeDirectRequest(t, "POST", path, params, 0) }
}

func TestAccSonarqubeUserDeletedOutOfBand(t *testing.T) {
	login := "testAccUserOutOfBand"
	testAccDeletedOutOfBand(t, noPreCheck,
		testAccSonarqubeUserNotLocalConfig(generateRandomResourceName(), login, "oob@example.org"),
		testAccPost(t, "/api/users/deactivate", url.Values{"login": {login}}))
}

func TestAccSonarqubeWebhookDeletedOutOfBand(t *testing.T) {
	name := "testAccWebhookOutOfBand"
	testAccDeletedOutOfBand(t, noPreCheck,
		testAccSonarqubeWebhookBasicConfig(generateRandomResourceName(), name, "https://example.org/oob", "secret-out-of-band"),
		func() {
			var list struct {
				Webhooks []struct{ Key, Name string } `json:"webhooks"`
			}
			if err := json.Unmarshal(testAccSonarqubeDirectRequest(t, "GET", "/api/webhooks/list", url.Values{}, http.StatusOK), &list); err != nil {
				t.Fatal(err)
			}
			for _, w := range list.Webhooks {
				if w.Name == name {
					testAccPost(t, "/api/webhooks/delete", url.Values{"webhook": {w.Key}})()
				}
			}
		})
}

func TestAccSonarqubeQualityProfileDeletedOutOfBand(t *testing.T) {
	name := "testAccQualityProfileOutOfBand"
	testAccDeletedOutOfBand(t, noPreCheck,
		testAccSonarqubeQualityProfileBasicConfig(generateRandomResourceName(), name, "js"),
		testAccPost(t, "/api/qualityprofiles/delete", url.Values{"language": {"js"}, "qualityProfile": {name}}))
}

func TestAccSonarqubePermissionTemplateDeletedOutOfBand(t *testing.T) {
	name := "testAccPermissionTemplateOutOfBand"
	testAccDeletedOutOfBand(t, noPreCheck,
		testAccSonarqubePermissionTemplateBasicConfig(generateRandomResourceName(), name, "oob", "oob.*"),
		testAccPost(t, "/api/permissions/delete_template", url.Values{"templateName": {name}}))
}

func TestAccSonarqubePermissionsDeletedOutOfBand(t *testing.T) {
	rnd := generateRandomResourceName()
	project, group := "testAccPermOOB-"+rnd, "testAccPermOOBGroup-"+rnd
	config := fmt.Sprintf(`
		resource "sonarqube_project" "%[1]s" {
			name       = "%[2]s"
			project    = "%[2]s"
			visibility = "private"
		}
		resource "sonarqube_group" "%[1]s" {
			name = "%[3]s"
		}
		resource "sonarqube_permissions" "%[1]s" {
			group_name  = sonarqube_group.%[1]s.name
			project_key = sonarqube_project.%[1]s.project
			permissions = ["issueadmin"]
		}`, rnd, project, group)
	testAccDeletedOutOfBand(t, noPreCheck, config,
		testAccPost(t, "/api/user_groups/delete", url.Values{"name": {group}}))
}

func TestAccSonarqubeSettingDeletedOutOfBand(t *testing.T) {
	key := "sonar.demo"
	testAccDeletedOutOfBand(t, noPreCheck,
		testAccSonarqubeSettingBasicConfig(generateRandomResourceName(), key, "oob@example.org"),
		testAccPost(t, "/api/settings/reset", url.Values{"keys": {key}}))
}

func TestAccSonarqubeNewCodePeriodsProjectDeletedOutOfBand(t *testing.T) {
	rnd := generateRandomResourceName()
	testAccDeletedOutOfBand(t, noPreCheck,
		testAccSonarqubeNewCodePeriodsProjectPreviousVersion(rnd),
		testAccPost(t, "/api/projects/delete", url.Values{"project": {rnd}}))
}

func TestAccSonarqubeProjectMainBranchDeletedOutOfBand(t *testing.T) {
	project := "testAccMainBranchOutOfBand"
	testAccDeletedOutOfBand(t, noPreCheck,
		testAccSonarqubeProjectMainBranchName(generateRandomResourceName(), project, "oob"),
		testAccPost(t, "/api/project_branches/rename", url.Values{"project": {project}, "name": {"renamed-oob"}}))
}

func TestAccSonarqubeRuleDeletedOutOfBand(t *testing.T) {
	key := "ruleOutOfBand"
	testAccDeletedOutOfBand(t, noPreCheck,
		testAccSonarqubeRuleBasicConfig(generateRandomResourceName(), key, "markdown_description", "name", "xml:XPathCheck", "INFO", "READY", "VULNERABILITY"),
		testAccPost(t, "/api/rules/delete", url.Values{"key": {"xml:" + key}}))
}

func TestAccSonarqubeQualityProfileRuleDeletedOutOfBand(t *testing.T) {
	rule := "activateRuleOutOfBand"
	testAccDeletedOutOfBand(t, noPreCheck,
		testAccSonarqubeQualityprofileActivateRuleBasicConfig(generateRandomResourceName(), "testAccActivateRuleOutOfBand", rule, "BLOCKER"),
		testAccPost(t, "/api/rules/delete", url.Values{"key": {"xml:" + rule}}))
}

func TestAccSonarqubeProjectDeletedOutOfBand(t *testing.T) {
	rnd := generateRandomResourceName()
	testAccDeletedOutOfBand(t, noPreCheck,
		testAccSonarqubeProjectBasicConfig(rnd, "testAccProjectOutOfBand", "testAccProjectOutOfBand", "public"),
		testAccPost(t, "/api/projects/delete", url.Values{"project": {"testAccProjectOutOfBand"}}))
}
