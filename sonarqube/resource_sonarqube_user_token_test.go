package sonarqube

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
   retrieve_badge_token = false
  }
  resource "sonarqube_project" "other" {
   name = "%[1]s-other"
   project = "%[1]s-other"
   retrieve_badge_token = false
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
