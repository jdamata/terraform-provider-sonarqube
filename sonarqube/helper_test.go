package sonarqube

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
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

// Returns the version of SonarQube running on the instance the tests are running against
func getSonarQubeVersion() (*version.Version, error) {
	sonarQubeURL := strings.TrimSuffix(os.Getenv("SONAR_HOST"), "/") + "/api/server/version"

	resp, err := httpRequestHelper(
		retryablehttp.NewClient(),
		"GET",
		sonarQubeURL,
		http.StatusOK,
		"api/server",
	)
	if err != nil {
		return &version.Version{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return &version.Version{}, err
	}
	version, _ := version.NewVersion(string(bodyBytes))
	return version, nil
}
