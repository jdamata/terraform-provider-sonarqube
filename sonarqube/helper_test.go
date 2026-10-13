package sonarqube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func generateRandomResourceName() string {
	return acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
}

func generateHCLList(s []string) string {
	semiformat := fmt.Sprintf("%+q", s)      // Turn the slice into a string that looks like ["one" "two" "three"]
	tokens := strings.Split(semiformat, " ") // Split this string by spaces
	return strings.Join(tokens, ", ")        // Join the Slice together (that was split by spaces) with commas
}

// Turns the map into a string that looks like {"one": "two", "three": "four"}
func generateHCLMap(m map[string]string) string {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "{")
	for key, value := range m {
		fmt.Fprintf(b, "\"%s\":\"%s\", ", key, value)
	}
	b.Truncate(b.Len() - 2)
	fmt.Fprintf(b, "}")
	return b.String()
}

// Asks the server instead of reading testAccProvider.Meta(), which is nil until a test has
// configured the provider, so the answer does not depend on test order.
func testAccIsCommunityEdition(t *testing.T) bool {
	var global struct {
		Edition string `json:"edition"`
	}
	if err := json.Unmarshal(testAccSonarqubeDirectRequest(t, "GET", "/api/navigation/global", url.Values{}, http.StatusOK), &global); err != nil {
		t.Fatal(err)
	}
	return strings.EqualFold(global.Edition, "community")
}
