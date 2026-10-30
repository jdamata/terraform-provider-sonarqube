package sonarqube

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testingTerraform "github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("sonarqube_user_token", &resource.Sweeper{
		Name: "sonarqube_user_token",
		F:    testSweepSonarqubeUserTokenSweeper,
	})
}

func testSweepSonarqubeUserTokenSweeper(r string) error {
	return nil
}

func testAccSonarqubeUserTokenBasicConfig(rnd string, name string) string {
	return fmt.Sprintf(`
        resource "sonarqube_user" "%[1]s" {
            login_name = "%[2]s"
            name       = "%[2]s"
            password   = "secret-sauce37!"
        }
        resource "sonarqube_user_token" "%[1]s" {
            login_name = sonarqube_user.%[1]s.login_name
            name       = "%[2]s"
        }`, rnd, name)
}

func TestAccSonarqubeUserTokenBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_user_token." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeUserTokenBasicConfig(rnd, "testAccSonarqubeUserToken"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeUserToken"),
					resource.TestCheckResourceAttr(name, "type", string(UserToken)),
				),
			},
		},
	})
}

func testAccSonarqubeUserTokenExpirationDateConfig(rnd string, name string, expiration_date string) string {
	return fmt.Sprintf(`
        resource "sonarqube_user" "%[1]s" {
            login_name = "%[2]s"
            name       = "%[2]s"
            password   = "secret-sauce37!"
        }
        resource "sonarqube_user_token" "%[1]s" {
            login_name      = sonarqube_user.%[1]s.login_name
            name            = "%[2]s"
            expiration_date = "%s"
        }`, rnd, name, expiration_date)
}

func TestAccSonarqubeUserTokenWithExpirationDate(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_user_token." + rnd
	expiration_date := time.Now().AddDate(0, 0, 5).Format("2006-01-02")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeUserTokenExpirationDateConfig(rnd, "testAccSonarqubeUserTokenWithExpirationDate", expiration_date),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeUserTokenWithExpirationDate"),
					resource.TestCheckResourceAttr(name, "expiration_date", expiration_date),
					resource.TestCheckResourceAttr(name, "type", string(UserToken)),
				),
			},
		},
	})
}

func testAccSonarqubeUserTokenNoLoginConfig(rnd string, name string) string {
	return fmt.Sprintf(`
        resource "sonarqube_user" "%[1]s" {
            login_name = "%[2]s"
            name       = "%[2]s"
            password   = "secret-sauce37!"
        }
        resource "sonarqube_user_token" "%[1]s" {
            name       = "%[2]s"
        }`, rnd, name)
}

func TestAccSonarqubeUserTokenNoLogin(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_user_token." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeUserTokenNoLoginConfig(rnd, "testAccSonarqubeUserTokenNoLogin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeUserTokenNoLogin"),
					resource.TestCheckResourceAttr(name, "type", string(UserToken)),
				),
			},
		},
	})
}

func testAccSonarqubeUserTokenGlobalAnalysisTokenConfig(rnd string, name string) string {
	return fmt.Sprintf(`
        resource "sonarqube_user" "%[1]s" {
            login_name = "%[2]s"
            name       = "%[2]s"
            password   = "secret-sauce37!"
        }
        resource "sonarqube_user_token" "%[1]s" {
            name       = "%[2]s"
            type       = "GLOBAL_ANALYSIS_TOKEN"
        }`, rnd, name)
}

func TestAccSonarqubeUserTokenGlobalAnalysisToken(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_user_token." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeUserTokenGlobalAnalysisTokenConfig(rnd, "testAccSonarqubeUserTokenGlobalAnalysisToken"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeUserTokenGlobalAnalysisToken"),
					resource.TestCheckResourceAttr(name, "type", string(GlobalAnalysisToken)),
				),
			},
		},
	})
}

func testAccSonarqubeUserTokenProjectAnalysisTokenConfig(rnd string, name string) string {
	return fmt.Sprintf(`
        resource "sonarqube_user" "%[1]s" {
            login_name = "%[2]s"
            name       = "%[2]s"
            password   = "secret-sauce37!"
        }
		resource "sonarqube_project" "%[1]s" {
			name = "my-project"
			project = "my-project"
			visibility = "public"
		}
        resource "sonarqube_user_token" "%[1]s" {
            name        = "%[2]s"
            type        = "PROJECT_ANALYSIS_TOKEN"
            project_key = sonarqube_project.%[1]s.project
        }`, rnd, name)
}

func TestAccSonarqubeUserTokenProjectAnalysisToken(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_user_token." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeUserTokenProjectAnalysisTokenConfig(rnd, "testAccSonarqubeUserTokenProjectAnalysisToken"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", "testAccSonarqubeUserTokenProjectAnalysisToken"),
					resource.TestCheckResourceAttr(name, "type", string(ProjectAnalysisToken)),
				),
			},
		},
	})
}

func TestAccSonarqubeUserTokenRevokedOutOfBand(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_user_token." + rnd
	tokenName := "testAccSonarqubeUserTokenRevokedOutOfBand"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeUserTokenBasicConfig(rnd, tokenName),
				Check:  resource.TestCheckResourceAttr(name, "name", tokenName),
			},
			{
				// Revoke behind Terraform's back, then refresh. Before the fix this step
				// failed with "Failed to find user token"; now the resource drops out of
				// state and the plan proposes recreating it.
				PreConfig:          func() { testAccRevokeUserToken(t, tokenName, tokenName) },
				Config:             testAccSonarqubeUserTokenBasicConfig(rnd, tokenName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRevokeUserToken(t *testing.T, login string, tokenName string) {
	conf := testAccProvider.Meta().(*ProviderConfiguration)
	sonarQubeURL := conf.sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/user_tokens/revoke"
	sonarQubeURL.RawQuery = url.Values{
		"login": []string{login},
		"name":  []string{tokenName},
	}.Encode()
	resp, err := httpRequestHelper(conf.httpClient, "POST", sonarQubeURL.String(), http.StatusNoContent, "testAccRevokeUserToken")
	if err != nil {
		t.Fatalf("revoking token out of band: %v", err)
	}
	defer resp.Body.Close()
}

func TestAccSonarqubeUserTokenMetadataOnly(t *testing.T) {
	name := "metadata-only-" + generateRandomResourceName()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{PreConfig: func() { testAccIssueExternalToken(t, name) },
				Config: fmt.Sprintf(`resource "sonarqube_user_token" "test" {
     name = %q
     metadata_only = true
    }`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sonarqube_user_token.test", "name", name),
					resource.TestCheckResourceAttr("sonarqube_user_token.test", "token", ""),
				)},
			{ResourceName: "sonarqube_user_token.test", ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"metadata_only", "login_name", "token"}},
		},
	})
}

func TestAccSonarqubeUserTokenMetadataOnlyMatchingMetadata(t *testing.T) {
	expiration := time.Now().UTC().AddDate(0, 0, 30).Format("2006-01-02")
	for _, tc := range []struct {
		name       string
		tokenType  TokenType
		expiration string
	}{
		{name: "project", tokenType: ProjectAnalysisToken},
		{name: "project_with_expiration", tokenType: ProjectAnalysisToken, expiration: expiration},
		{name: "user_with_expiration", tokenType: UserToken, expiration: expiration},
	} {
		t.Run(tc.name, func(t *testing.T) {
			name := "metadata-match-" + generateRandomResourceName()
			projectConfig := ""
			projectKeyConfig := ""
			options := url.Values{"type": {string(tc.tokenType)}}
			checks := []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("sonarqube_user_token.test", "name", name),
				resource.TestCheckResourceAttr("sonarqube_user_token.test", "type", string(tc.tokenType)),
				resource.TestCheckResourceAttr("sonarqube_user_token.test", "metadata_only", "true"),
				resource.TestCheckResourceAttr("sonarqube_user_token.test", "token", ""),
			}
			if tc.tokenType == ProjectAnalysisToken {
				projectConfig = fmt.Sprintf(`
resource "sonarqube_project" "test" {
  name                 = %[1]q
  project              = %[1]q
}
`, name)
				projectKeyConfig = "project_key = sonarqube_project.test.project"
				options.Set("projectKey", name)
				checks = append(checks, resource.TestCheckResourceAttr("sonarqube_user_token.test", "project_key", name))
			}
			expirationConfig := ""
			if tc.expiration != "" {
				options.Set("expirationDate", tc.expiration)
				expirationConfig = fmt.Sprintf("expiration_date = %q", tc.expiration)
				checks = append(checks, resource.TestCheckResourceAttr("sonarqube_user_token.test", "expiration_date", tc.expiration))
			}
			config := projectConfig + fmt.Sprintf(`
resource "sonarqube_user_token" "test" {
  name          = %q
  type          = %q
  metadata_only = true
  %s
  %s
}
`, name, tc.tokenType, projectKeyConfig, expirationConfig)
			steps := []resource.TestStep{}
			if projectConfig != "" {
				steps = append(steps, resource.TestStep{Config: projectConfig})
			}
			steps = append(steps,
				resource.TestStep{
					PreConfig: func() { testAccIssueExternalToken(t, name, options) },
					Config:    config,
					Check:     resource.ComposeTestCheckFunc(checks...),
				},
				resource.TestStep{Config: config, PlanOnly: true},
			)
			resource.Test(t, resource.TestCase{
				PreCheck: func() { testAccPreCheck(t) }, Providers: testAccProviders,
				Steps: steps,
			})
		})
	}
}

func TestAccSonarqubeUserTokenMetadataOnlyMissing(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) }, Providers: testAccProviders,
		Steps: []resource.TestStep{{
			Config: fmt.Sprintf(`resource "sonarqube_user_token" "test" {
    name = %q
    metadata_only = true
   }`, "missing-"+generateRandomResourceName()),
			ExpectError: regexp.MustCompile("must already exist when metadata_only is enabled"),
		}},
	})
}

func TestAccSonarqubeUserTokenMetadataOnlyRejectsType(t *testing.T) {
	name := "metadata-type-" + generateRandomResourceName()
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) }, Providers: testAccProviders,
		Steps: []resource.TestStep{{
			PreConfig: func() { testAccIssueExternalToken(t, name) },
			Config: fmt.Sprintf(`resource "sonarqube_user_token" "test" {
    name = %q
    type = "GLOBAL_ANALYSIS_TOKEN"
    metadata_only = true
   }`, name),
			ExpectError: regexp.MustCompile("has type.*expected"),
		}},
	})
}

func TestAccSonarqubeUserTokenMetadataOnlyMigration(t *testing.T) {
	name := "metadata-migration-" + generateRandomResourceName()
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) }, Providers: testAccProviders,
		Steps: []resource.TestStep{
			{Config: fmt.Sprintf(`resource "sonarqube_user_token" "test" { name = %q }`, name),
				Check: resource.TestCheckResourceAttrSet("sonarqube_user_token.test", "token")},
			{Config: fmt.Sprintf(`resource "sonarqube_user_token" "test" {
    name = %q
    metadata_only = true
   }`, name),
				Check: resource.TestCheckResourceAttr("sonarqube_user_token.test", "token", "")},
			{Config: fmt.Sprintf(`resource "sonarqube_user_token" "test" { name = %q }`, name),
				Check: resource.TestCheckResourceAttr("sonarqube_user_token.test", "token", "")},
		},
	})
}

// Issue a fixture token outside Terraform and discard its value. Only the token
// name is needed to test metadata adoption; no credential is logged or retained.
func testAccIssueExternalToken(t *testing.T, name string, options ...url.Values) {
	t.Helper()
	query := url.Values{"name": {name}}
	for _, option := range options {
		for key, values := range option {
			query[key] = values
		}
	}
	endpoint := strings.TrimSuffix(os.Getenv("SONAR_HOST"), "/") + "/api/user_tokens/generate?" + query.Encode()
	request, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}
	if token := os.Getenv("SONAR_TOKEN"); token != "" {
		request.SetBasicAuth(token, "")
	} else {
		request.SetBasicAuth(os.Getenv("SONAR_USER"), os.Getenv("SONAR_PASS"))
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("fixture token creation returned HTTP %d", response.StatusCode)
	}
	t.Cleanup(func() {
		cleanup := strings.TrimSuffix(os.Getenv("SONAR_HOST"), "/") + "/api/user_tokens/revoke?" + url.Values{"name": {name}}.Encode()
		request, err := http.NewRequest("POST", cleanup, nil)
		if err != nil {
			t.Error(err)
			return
		}
		if token := os.Getenv("SONAR_TOKEN"); token != "" {
			request.SetBasicAuth(token, "")
		} else {
			request.SetBasicAuth(os.Getenv("SONAR_USER"), os.Getenv("SONAR_PASS"))
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Error(err)
			return
		}
		response.Body.Close()
	})
}

func TestAccSonarqubeUserTokenMetadataOnlyRejectsProject(t *testing.T) {
	name := "metadata-project-" + generateRandomResourceName()
	projects := fmt.Sprintf(`
  resource "sonarqube_project" "source" {
   name = "%[1]s-source"
   project = "%[1]s-source"
   }
  resource "sonarqube_project" "other" {
   name = "%[1]s-other"
   project = "%[1]s-other"
   }
 `, name)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) }, Providers: testAccProviders,
		Steps: []resource.TestStep{
			{Config: projects},
			{PreConfig: func() {
				testAccIssueExternalToken(t, name, url.Values{"type": {string(ProjectAnalysisToken)}, "projectKey": {name + "-source"}})
			},
				Config: projects + fmt.Sprintf(`
     resource "sonarqube_user_token" "test" {
      name = %q
      type = "PROJECT_ANALYSIS_TOKEN"
      project_key = sonarqube_project.other.project
      metadata_only = true
     }
    `, name),
				ExpectError: regexp.MustCompile("belongs to project.*expected"),
			},
		},
	})
}

func TestAccSonarqubeUserTokenMetadataOnlyRejectsExpiration(t *testing.T) {
	name := "metadata-expiration-" + generateRandomResourceName()
	actual := time.Now().UTC().AddDate(0, 0, 30).Format("2006-01-02")
	expected := time.Now().UTC().AddDate(0, 0, 60).Format("2006-01-02")
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) }, Providers: testAccProviders,
		Steps: []resource.TestStep{{
			PreConfig: func() { testAccIssueExternalToken(t, name, url.Values{"expirationDate": {actual}}) },
			Config: fmt.Sprintf(`resource "sonarqube_user_token" "test" {
    name = %q
    metadata_only = true
    expiration_date = %q
   }`, name, expected),
			ExpectError: regexp.MustCompile("does not have the expected expiration date"),
		}},
	})
}

func TestMetadataOnlyTokenReplacementDiff(t *testing.T) {
	original := map[string]string{
		"name": "issued", "type": "PROJECT_ANALYSIS_TOKEN", "project_key": "project",
		"expiration_date": "2030-01-01", "login_name": "owner", "metadata_only": "true",
	}
	for _, field := range []string{"name", "type", "project_key", "expiration_date", "login_name"} {
		for _, modes := range []struct{ old, new bool }{{true, true}, {true, false}, {false, true}, {false, false}} {
			t.Run(fmt.Sprintf("%s/%v-to-%v", field, modes.old, modes.new), func(t *testing.T) {
				attrs := map[string]string{}
				config := map[string]interface{}{}
				for key, value := range original {
					attrs[key] = value
					config[key] = value
				}
				attrs["metadata_only"] = fmt.Sprint(modes.old)
				config["metadata_only"] = modes.new
				changes := map[string]string{"name": "next", "type": "USER_TOKEN", "project_key": "other", "expiration_date": "2030-02-01", "login_name": "other"}
				config[field] = changes[field]
				_, err := resourceSonarqubeUserToken().Diff(context.Background(), &terraform.InstanceState{
					ID: "owner/issued", Attributes: attrs,
				}, terraform.NewResourceConfigRaw(config), nil)
				wantError := modes.old || modes.new
				if (err != nil) != wantError {
					t.Fatalf("Diff error = %v, want error %v", err, wantError)
				}
				if err != nil && !strings.Contains(err.Error(), "cannot change "+field) {
					t.Fatal(err)
				}
			})
		}
	}
	for _, mode := range []bool{false, true} {
		t.Run(fmt.Sprintf("mode-only/%v", mode), func(t *testing.T) {
			config := map[string]interface{}{}
			for key, value := range original {
				config[key] = value
			}
			config["metadata_only"] = mode
			if _, err := resourceSonarqubeUserToken().Diff(context.Background(), &terraform.InstanceState{
				ID: "owner/issued", Attributes: original,
			}, terraform.NewResourceConfigRaw(config), nil); err != nil {
				t.Fatal(err)
			}
		})
	}
	t.Run("initial-adoption", func(t *testing.T) {
		if _, err := resourceSonarqubeUserToken().Diff(context.Background(), nil, terraform.NewResourceConfigRaw(map[string]interface{}{
			"name": "74D93920-ED26-11E3-AC10-0800200C9A66", "metadata_only": true,
		}), nil); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("unknown-replacement-name", func(t *testing.T) {
		config := map[string]interface{}{}
		for key, value := range original {
			config[key] = value
		}
		config["metadata_only"] = true
		config["name"] = "74D93920-ED26-11E3-AC10-0800200C9A66"
		if _, err := resourceSonarqubeUserToken().Diff(context.Background(), &terraform.InstanceState{
			ID: "owner/issued", Attributes: original,
		}, terraform.NewResourceConfigRaw(config), nil); err == nil {
			t.Fatal("unknown replacement name must be refused")
		}
	})
}

func TestAccSonarqubeUserTokenMetadataOnlyReplacementRefused(t *testing.T) {
	name := "metadata-guard-" + generateRandomResourceName()
	projects := fmt.Sprintf(`
resource "sonarqube_project" "source" {
  name = %[1]q
  project = %[1]q
}
resource "sonarqube_project" "other" {
  name = "%[1]s-other"
  project = "%[1]s-other"
}
`, name)
	base := map[string]string{
		"name": fmt.Sprintf("%q", name), "type": `"PROJECT_ANALYSIS_TOKEN"`,
		"project_key": "sonarqube_project.source.project", "metadata_only": "true",
	}
	configFor := func(field, value string) string {
		values := map[string]string{}
		for key, v := range base {
			values[key] = v
		}
		if field != "" {
			values[field] = value
		}
		config := projects + `resource "sonarqube_user_token" "test" {` + "\n"
		for _, key := range []string{"name", "type", "project_key", "metadata_only", "expiration_date", "login_name"} {
			if v, ok := values[key]; ok {
				config += key + " = " + v + "\n"
			}
		}
		return config + "}"
	}
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("sonarqube_user_token.test", "token", ""),
		func(_ *testingTerraform.State) error {
			tokens, err := readSonarqubeUserTokens(testAccProvider.Meta(), "")
			if err != nil {
				return err
			}
			for _, token := range tokens.Tokens {
				if token.Name == name && token.Project.Key == name {
					return nil
				}
			}
			return fmt.Errorf("adopted token %q was revoked after a rejected plan", name)
		},
	)
	steps := []resource.TestStep{
		{Config: projects},
		{PreConfig: func() {
			testAccIssueExternalToken(t, name, url.Values{"type": {string(ProjectAnalysisToken)}, "projectKey": {name}})
		}, Config: configFor("", ""), Check: check},
	}
	for _, change := range []struct{ field, value string }{
		{"project_key", "sonarqube_project.other.project"}, {"name", fmt.Sprintf("%q", name+"-next")},
		{"type", `"USER_TOKEN"`}, {"expiration_date", fmt.Sprintf("%q", time.Now().UTC().AddDate(0, 0, 30).Format("2006-01-02"))},
		{"login_name", `"another-owner"`},
	} {
		steps = append(steps,
			resource.TestStep{Config: configFor(change.field, change.value), ExpectError: regexp.MustCompile("cannot change " + change.field)},
			resource.TestStep{Config: configFor("", ""), Check: check},
		)
	}
	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheck(t) }, Providers: testAccProviders, Steps: steps})
}

func TestAccSonarqubeUserTokenMetadataOnlyRotationSeparateAddresses(t *testing.T) {
	name := "metadata-rotation-" + generateRandomResourceName()
	oldName, newName := name+"-old", name+"-new"
	oldConfig := fmt.Sprintf(`resource "sonarqube_user_token" "old" {
  name = %q
  metadata_only = true
}
`, oldName)
	newConfig := fmt.Sprintf(`resource "sonarqube_user_token" "next" {
  name = %q
  metadata_only = true
}
`, newName)
	check := func(wantOld bool) resource.TestCheckFunc {
		return func(_ *testingTerraform.State) error {
			tokens, err := readSonarqubeUserTokens(testAccProvider.Meta(), "")
			if err != nil {
				return err
			}
			foundOld, foundNew := false, false
			for _, token := range tokens.Tokens {
				foundOld = foundOld || token.Name == oldName
				foundNew = foundNew || token.Name == newName
			}
			if foundOld != wantOld || !foundNew {
				return fmt.Errorf("rotation token presence: old=%v (want %v), new=%v (want true)", foundOld, wantOld, foundNew)
			}
			return nil
		}
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) }, Providers: testAccProviders,
		Steps: []resource.TestStep{
			{PreConfig: func() { testAccIssueExternalToken(t, oldName) }, Config: oldConfig},
			{PreConfig: func() { testAccIssueExternalToken(t, newName) }, Config: oldConfig + newConfig,
				Check: resource.ComposeTestCheckFunc(check(true), resource.TestCheckResourceAttr("sonarqube_user_token.next", "token", ""))},
			{Config: newConfig, Check: check(false)},
			{Config: newConfig, PlanOnly: true},
		},
	})
}
