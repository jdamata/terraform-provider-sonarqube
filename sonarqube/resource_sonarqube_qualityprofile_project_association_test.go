package sonarqube

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("sonarqube_qualityprofile_project_association", &resource.Sweeper{
		Name: "sonarqube_qualityprofile_project_association",
		F:    testSweepSonarqubeQualityProfileProjectAssociationSweeper,
	})
}

func testSweepSonarqubeQualityProfileProjectAssociationSweeper(r string) error {
	return nil
}

func testAccSonarqubeQualityProfileProjectAssociationBasicConfig(rnd string, name string, language string) string {
	return fmt.Sprintf(`
		resource "sonarqube_qualityprofile" "%[1]s" {
			name     = "%[2]s"
			language = "%[3]s"
		}

		resource "sonarqube_project" "%[1]s" {
			name       = "%[2]s"
			project    = "%[2]s"
			visibility = "public" 
		}

		resource "sonarqube_qualityprofile_project_association" "%[1]s" {
			quality_profile = sonarqube_qualityprofile.%[1]s.name
			project         = sonarqube_project.%[1]s.name
			language        = "%[3]s"
		}`, rnd, name, language)
}

func TestAccSonarqubeQualityProfileProjectAssociationBasic(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualityprofile_project_association." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualityProfileProjectAssociationBasicConfig(rnd, "testAccSonarqubeProfileProjectAssociation", "js"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "quality_profile", "testAccSonarqubeProfileProjectAssociation"),
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProfileProjectAssociation"),
					resource.TestCheckResourceAttr(name, "language", "js"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "quality_profile", "testAccSonarqubeProfileProjectAssociation"),
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProfileProjectAssociation"),
					resource.TestCheckResourceAttr(name, "language", "js"),
				),
			},
		},
	})
}

func testAccSonarqubeQualityProfileProjectAssociationSonarWay(rnd string, name string, language string, qualityProfile string) string {
	return fmt.Sprintf(`
		resource "sonarqube_project" "%[1]s" {
			name       = "%[2]s"
			project    = "%[2]s"
			visibility = "public" 
		}

		resource "sonarqube_qualityprofile_project_association" "%[1]s" {
			quality_profile = "%[4]s"
			project         = sonarqube_project.%[1]s.name
			language        = "%[3]s"
		}`, rnd, name, language, qualityProfile)
}

func TestAccSonarqubeQualityProfileProjectAssociationSonarWay(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "sonarqube_qualityprofile_project_association." + rnd
	builtIn := testAccBuiltInQualityProfile(t, "js")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSonarqubeQualityProfileProjectAssociationSonarWay(rnd, "testAccSonarqubeProfileProjectAssociation", "js", builtIn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "quality_profile", builtIn),
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProfileProjectAssociation"),
					resource.TestCheckResourceAttr(name, "language", "js"),
				),
			},
			{
				ResourceName:      name,
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "quality_profile", builtIn),
					resource.TestCheckResourceAttr(name, "project", "testAccSonarqubeProfileProjectAssociation"),
					resource.TestCheckResourceAttr(name, "language", "js"),
				),
			},
		},
	})
}

// testAccSonarqubeDirectRequest calls the SonarQube API with the test credentials. It does not
// go through the provider, so it works before the provider has been configured.
func testAccSonarqubeDirectRequest(t *testing.T, method string, path string, params url.Values, expectedStatus int) []byte {
	// Callers run this before resource.Test, which is what normally skips acceptance tests
	// when TF_ACC is unset, so without this check plain go test ./... fails here.
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance test: set TF_ACC to run")
	}
	testAccPreCheck(t)
	target := strings.TrimSuffix(os.Getenv("SONAR_HOST"), "/") + path + "?" + params.Encode()
	req, err := http.NewRequest(method, target, http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	if token := os.Getenv("SONAR_TOKEN"); token != "" {
		req.SetBasicAuth(token, "")
	} else {
		req.SetBasicAuth(os.Getenv("SONAR_USER"), os.Getenv("SONAR_PASS"))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != expectedStatus {
		t.Fatalf("%s %s: got HTTP %d, want %d: %s", method, path, resp.StatusCode, expectedStatus, body)
	}
	return body
}

// testAccBuiltInQualityProfile returns the name of the default built-in profile for a language.
// Built-in names differ across SonarQube versions and editions ("Sonar way" versus
// "Sonar way core"), so tests must not hardcode them.
func testAccBuiltInQualityProfile(t *testing.T, language string) string {
	body := testAccSonarqubeDirectRequest(t, "GET", "/api/qualityprofiles/search", url.Values{"language": {language}}, http.StatusOK)
	var profiles struct {
		Profiles []struct {
			Name      string `json:"name"`
			IsBuiltIn bool   `json:"isBuiltIn"`
			IsDefault bool   `json:"isDefault"`
		} `json:"profiles"`
	}
	if err := json.Unmarshal(body, &profiles); err != nil {
		t.Fatal(err)
	}
	fallback := ""
	for _, p := range profiles.Profiles {
		if !p.IsBuiltIn {
			continue
		}
		if p.IsDefault {
			return p.Name
		}
		if fallback == "" {
			fallback = p.Name
		}
	}
	if fallback == "" {
		t.Fatalf("no built-in quality profile for language %q", language)
	}
	return fallback
}

// A profile that Terraform does not manage is deleted between steps. Read must drop the
// association rather than send an empty profile key and fail the refresh.
func TestAccSonarqubeQualityProfileProjectAssociationProfileDeletedOutOfBand(t *testing.T) {
	rnd := generateRandomResourceName()
	profile := "testAccProfileAssociationProfileGone"
	testAccSonarqubeDirectRequest(t, "POST", "/api/qualityprofiles/create", url.Values{"name": {profile}, "language": {"py"}}, http.StatusOK)
	config := testAccSonarqubeQualityProfileProjectAssociationSonarWay(rnd, "testAccProfileAssociationProfileGone", "py", profile)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() {
					testAccSonarqubeDirectRequest(t, "POST", "/api/qualityprofiles/delete", url.Values{"qualityProfile": {profile}, "language": {"py"}}, http.StatusNoContent)
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// The association is removed in SonarQube between steps. Read must drop it from state so the
// plan proposes recreating it.
func TestAccSonarqubeQualityProfileProjectAssociationRemovedOutOfBand(t *testing.T) {
	rnd := generateRandomResourceName()
	name := "testAccProfileAssociationRemoved"
	config := testAccSonarqubeQualityProfileProjectAssociationBasicConfig(rnd, name, "py")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() {
					testAccSonarqubeDirectRequest(t, "POST", "/api/qualityprofiles/remove_project", url.Values{"qualityProfile": {name}, "language": {"py"}, "project": {name}}, http.StatusNoContent)
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
