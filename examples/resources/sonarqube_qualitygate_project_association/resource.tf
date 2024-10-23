## Example: create a quality gate project association

resource "sonarqube_qualitygate" "main" {
  name = "my_qualitygate"

  condition {
    metric    = "new_coverage"
    op        = "LT"
    threshold = "30"
  }
}

resource "sonarqube_project" "main" {
  name       = "SonarQube"
  project    = "my_project"
  visibility = "public"
}

resource "sonarqube_qualitygate_project_association" "main" {
  gatename   = sonarqube_qualitygate.main.id
  projectkey = sonarqube_project.main.project
}
