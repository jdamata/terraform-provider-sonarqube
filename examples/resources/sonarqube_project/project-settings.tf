resource "sonarqube_project" "main" {
  name       = "SonarQube"
  project    = "my_project"
  visibility = "public"

  setting {
    key   = "sonar.demo"
    value = "sonarqube@example.org"
  }
}
