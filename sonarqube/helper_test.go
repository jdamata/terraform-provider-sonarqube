package sonarqube

import (
	"fmt"
	"strings"

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
