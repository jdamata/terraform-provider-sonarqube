resource "sonarqube_project" "reference" {
  name = "my-project"
}

resource "sonarqube_new_code_periods" "reference" {
  project = sonarqube_project.reference.project
  type = "REFERENCE_BRANCH"
  value = "main"
