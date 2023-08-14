package sonarqube

import (
	"bytes"
	"fmt"
	"strings"

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
