package sonarqube

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"sonarqube": testAccProvider,
	}
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	testSonarHost(t)
	if v := os.Getenv("SONAR_TOKEN"); v == "" {
		testSonarUser(t)
		testSonarPass(t)
	}
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
