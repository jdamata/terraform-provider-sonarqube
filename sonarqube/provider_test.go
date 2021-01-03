package sonarqube

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"sonarqube": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testSonarPreCheck(t *testing.T) {
	testSonarHost(t)
	testSonarUser(t)
	testSonarPass(t)
	testSonarScheme(t)
}

func testSonarUser(t *testing.T) {
	if v := os.Getenv("SONAR_USER"); v == "" {
		t.Fatal("SONAR_USER must be set for this acceptance test")
	}
}

func testSonarPass(t *testing.T) {
	if v := os.Getenv("SONAR_PASS"); v == "" {
		t.Fatal("SONAR_PASS must be set for this acceptance test")
	}
}

func testSonarHost(t *testing.T) {
	if v := os.Getenv("SONAR_HOST"); v == "" {
		t.Fatal("SONAR_HOST must be set for this acceptance test")
	}
}

func testSonarScheme(t *testing.T) {
	if v := os.Getenv("SONAR_SCHEME"); v == "" {
		t.Fatal("SONAR_SCHEME must be set for this acceptance test")
	}
}

func generateRandomResourceName() string {
	return acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
}
