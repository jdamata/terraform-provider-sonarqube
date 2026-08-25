resource "sonarqube_project_links" "homepage" {
  project_key = "my_project"
  name        = "Homepage"
  url         = "https://example.com"
}
