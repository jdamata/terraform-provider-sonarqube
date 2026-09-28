package sonarqube

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Each test creates the resource, removes it through the API behind Terraform's back, then
// refreshes. Read must drop the resource from state so the plan proposes recreating it,
// instead of failing the refresh and blocking every later plan.
func testAccDeletedOutOfBand(t *testing.T, preCheck func(), config string, deleteOutOfBand func()) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preCheck()
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig:          deleteOutOfBand,
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccSonarqubeAPIPost(t *testing.T, path string, params url.Values) {
	conf := testAccProvider.Meta().(*ProviderConfiguration)
	sonarQubeURL := conf.sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + path
	sonarQubeURL.RawQuery = params.Encode()
	resp, err := httpRequestHelper(conf.httpClient, "POST", sonarQubeURL.String(), http.StatusNoContent, "testAccSonarqubeAPIPost")
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	defer resp.Body.Close()
}

func testAccDeleteAlmSetting(t *testing.T, key string) func() {
	return func() { testAccSonarqubeAPIPost(t, "/api/alm_settings/delete", url.Values{"key": {key}}) }
}

func testAccDeleteProjectBinding(t *testing.T, project string) func() {
	return func() {
		testAccSonarqubeAPIPost(t, "/api/alm_settings/delete_binding", url.Values{"project": {project}})
	}
}

func noPreCheck() {}

func TestAccSonarqubeAlmAzureDeletedOutOfBand(t *testing.T) {
	key := "testAccAlmAzureOutOfBand"
	testAccDeletedOutOfBand(t, noPreCheck,
		testAccSonarqubeAlmAzureName(generateRandomResourceName(), key, "https://dev.azure.com/my-org"),
		testAccDeleteAlmSetting(t, key))
}

func TestAccSonarqubeAlmBitbucketDeletedOutOfBand(t *testing.T) {
	key := "testAccAlmBitbucketOutOfBand"
	testAccDeletedOutOfBand(t, noPreCheck,
		testAccSonarqubeAlmBitbucketName(generateRandomResourceName(), key, "123456"),
		testAccDeleteAlmSetting(t, key))
}

func TestAccSonarqubeAlmGithubDeletedOutOfBand(t *testing.T) {
	key := "testAccAlmGithubOutOfBand"
	testAccDeletedOutOfBand(t, noPreCheck,
		testAccSonarqubeAlmGithubName(generateRandomResourceName(), key, "123456", "234567"),
		testAccDeleteAlmSetting(t, key))
}

func TestAccSonarqubeAlmGitlabDeletedOutOfBand(t *testing.T) {
	key := "testAccAlmGitlabOutOfBand"
	testAccDeletedOutOfBand(t, noPreCheck,
		testAccSonarqubeAlmGitlabName(generateRandomResourceName(), key, "123456"),
		testAccDeleteAlmSetting(t, key))
}

func TestAccSonarqubeAzureBindingDeletedOutOfBand(t *testing.T) {
	project := "testAccAzureBindingOutOfBand"
	testAccDeletedOutOfBand(t, func() { testAccPreCheckAzureBindingSupport(t) },
		testAccSonarqubeAzureBindingName(generateRandomResourceName(), project, "azureOutOfBand", "testAzProjName", "testAzRepoName"),
		testAccDeleteProjectBinding(t, project))
}

func TestAccSonarqubeBitbucketBindingDeletedOutOfBand(t *testing.T) {
	project := "testAccBitbucketBindingOutOfBand"
	testAccDeletedOutOfBand(t, func() { testAccPreCheckBitbucketBindingSupport(t) },
		testAccSonarqubeBitbucketBindingName(generateRandomResourceName(), project, "repo-key", "repo-slug"),
		testAccDeleteProjectBinding(t, project))
}

func TestAccSonarqubeGithubBindingDeletedOutOfBand(t *testing.T) {
	project := "testAccGithubBindingOutOfBand"
	testAccDeletedOutOfBand(t, func() { testAccPreCheckGithubBindingSupport(t) },
		testAccSonarqubeGithubBindingName(generateRandomResourceName(), project, "githubOutOfBand", project),
		testAccDeleteProjectBinding(t, project))
}

func TestAccSonarqubeGitlabBindingDeletedOutOfBand(t *testing.T) {
	project := "testAccGitlabBindingOutOfBand"
	testAccDeletedOutOfBand(t, func() { testAccPreCheckGitlabBindingSupport(t) },
		testAccSonarqubeGitlabBindingName(generateRandomResourceName(), project, "gitlabOutOfBand", project),
		testAccDeleteProjectBinding(t, project))
}
