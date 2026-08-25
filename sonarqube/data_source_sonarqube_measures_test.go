package sonarqube

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccSonarqubeMeasuresDataSourceConfig(rnd, projectKey string, metricKeys []string) string {
	formattedKeys := generateHCLList(metricKeys)
	return fmt.Sprintf(`
resource "sonarqube_project" "%[1]s" {
  name       = "%[2]s"
  project    = "%[2]s"
  visibility = "public"
}

data "sonarqube_measures" "%[1]s" {
  component   = sonarqube_project.%[1]s.project
  metric_keys = %[3]s
}
`, rnd, projectKey, formattedKeys)
}

func TestAccSonarqubeMeasuresDataSource(t *testing.T) {
	rnd := generateRandomResourceName()
	dataName := "data.sonarqube_measures." + rnd
	projectKey := "testAccMeasures" + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// A freshly created project has no scan data, so measures will be empty.
				// This test verifies the data source calls the API successfully and sets its ID.
				Config: testAccSonarqubeMeasuresDataSourceConfig(rnd, projectKey, []string{"coverage", "reliability_rating", "ncloc"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataName, "id"),
					resource.TestCheckResourceAttr(dataName, "component", projectKey),
				),
			},
		},
	})
}
