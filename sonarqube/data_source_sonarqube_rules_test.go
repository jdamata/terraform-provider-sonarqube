package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccSonarqubeRulesTotal(t *testing.T, filters url.Values) int {
	filters.Set("ps", "1")
	var search struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal(testAccSonarqubeDirectRequest(t, "GET", "/api/rules/search", filters, http.StatusOK), &search); err != nil {
		t.Fatal(err)
	}
	return search.Total
}

func TestAccSonarqubeRulesDataSourceFilters(t *testing.T) {
	rnd := generateRandomResourceName()
	want := testAccSonarqubeRulesTotal(t, url.Values{
		"languages": {"java"}, "types": {"CODE_SMELL"}, "severities": {"CRITICAL"}, "statuses": {"READY"},
	})
	if want == 0 {
		t.Skip("no matching Java rules on this server")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "sonarqube_rules" "%s" {
						languages  = ["java"]
						types      = ["CODE_SMELL"]
						severities = ["CRITICAL"]
						statuses   = ["READY"]
					}`, rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sonarqube_rules."+rnd, "rules.#", strconv.Itoa(want)),
					resource.TestCheckResourceAttr("data.sonarqube_rules."+rnd, "rules.0.language", "java"),
					resource.TestCheckResourceAttr("data.sonarqube_rules."+rnd, "rules.0.severity", "CRITICAL"),
					resource.TestCheckResourceAttr("data.sonarqube_rules."+rnd, "rules.0.type", "CODE_SMELL"),
				),
			},
		},
	})
}

// More than one page of api/rules/search (500 rules), so a missing page shows up as a short count.
func TestAccSonarqubeRulesDataSourceAllPages(t *testing.T) {
	rnd := generateRandomResourceName()
	want := testAccSonarqubeRulesTotal(t, url.Values{"repositories": {"java,javascript,typescript,php,py,xml"}})
	if want <= 500 {
		t.Skipf("only %d rules: not enough to span pages on this server", want)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "sonarqube_rules" "%s" {
						repositories = ["java", "javascript", "typescript", "php", "py", "xml"]
					}`, rnd),
				Check: resource.TestCheckResourceAttr("data.sonarqube_rules."+rnd, "rules.#", strconv.Itoa(want)),
			},
		},
	})
}
