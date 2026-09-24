package sonarqube

import (
	"fmt"
	"net/http"
	"net/url"
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
