package sonarqube

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	testingTerraform "github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Exercise real Terraform planning and apply with an isolated SonarQube API fixture.
func TestProjectBadgeTokenPropagation(t *testing.T) {
	t.Setenv("SONAR_USER", "")
	t.Setenv("SONAR_PASS", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/projects/create":
			fmt.Fprint(w, `{"project":{"key":"badge-propagation","name":"badge-propagation"}}`)
		case "/api/components/show":
			fmt.Fprint(w, `{"component":{"key":"badge-propagation","name":"badge-propagation","visibility":"private"}}`)
		case "/api/project_badges/token":
			fmt.Fprint(w, `{"token":"fixture-badge-token"}`)
		case "/api/settings/values":
			fmt.Fprint(w, `{"settings":[]}`)
		case "/api/project_tags/set", "/api/projects/delete":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected API request: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: projectBadgePropagationSteps("badge-propagation", fmt.Sprintf(`
provider "sonarqube" {
 host = %q
 token = "fixture"
 installed_version = "26.10.0"
 installed_edition = "Community"
}
`, server.URL)),
	})
}

func TestAccSonarqubeProjectBadgeTokenOutputs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) }, Providers: testAccProviders,
		Steps: projectBadgePropagationSteps("badge-output-"+generateRandomResourceName(), ""),
	})
}

func projectBadgePropagationSteps(name, providerConfig string) []resource.TestStep {
	var steps []resource.TestStep
	// Change each flag independently and together, in both directions.
	for _, mode := range [][2]bool{{true, true}, {false, true}, {false, false}, {true, false}, {true, true}, {false, false}, {true, true}} {
		resourceEnabled, dataEnabled := mode[0], mode[1]
		config := providerConfig + fmt.Sprintf(`
resource "sonarqube_project" "test" {
 name = %[1]q
 project = %[1]q
 visibility = "private"
 retrieve_badge_token = %[2]t
}
data "sonarqube_project" "test" {
 project = sonarqube_project.test.project
 retrieve_badge_token = %[3]t
}
resource "terraform_data" "resource_consumer" {
 input = sonarqube_project.test.badge_token
}
resource "terraform_data" "data_consumer" {
 input = data.sonarqube_project.test.badge_token
}
output "resource_badge_url" {
 value = "https://sonarqube.example/badge?token=${sonarqube_project.test.badge_token}"
}
output "data_badge_url" {
 value = "https://sonarqube.example/badge?token=${data.sonarqube_project.test.badge_token}"
}
`, name, resourceEnabled, dataEnabled)
		steps = append(steps, resource.TestStep{Config: config, Check: func(state *testingTerraform.State) error {
			for _, item := range []struct {
				address, consumer, output string
				enabled                   bool
			}{
				{"sonarqube_project.test", "terraform_data.resource_consumer", "resource_badge_url", resourceEnabled},
				{"data.sonarqube_project.test", "terraform_data.data_consumer", "data_badge_url", dataEnabled},
			} {
				token := state.RootModule().Resources[item.address].Primary.Attributes["badge_token"]
				if (token != "") != item.enabled {
					return fmt.Errorf("%s: unexpected badge token presence", item.address)
				}
				output, ok := state.RootModule().Outputs[item.output]
				if !ok || output.Value != "https://sonarqube.example/badge?token="+token {
					return fmt.Errorf("%s retained a value from the previous mode", item.output)
				}
				consumer := state.RootModule().Resources[item.consumer]
				if consumer == nil || consumer.Primary.Attributes["input"] != token || consumer.Primary.Attributes["output"] != token {
					return fmt.Errorf("%s did not receive this apply's badge value", item.consumer)
				}
			}
			return nil
		}}, resource.TestStep{Config: config, PlanOnly: true})
	}
	return steps
}
