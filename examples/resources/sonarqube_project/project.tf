resource "sonarqube_project" "main" {
  name       = "SonarQube"
  project    = "my_project"
  visibility = "public"
}
