resource "sonarqube_project" "main" {
  name       = "SonarQube"
  project    = "my_project"
  visibility = "public"
}
resource "sonarqube_project_main_branch" "mybranch" {
  name    = "release"
  project = "my_project"
}
