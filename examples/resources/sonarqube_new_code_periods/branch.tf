resource "sonarqube_project" "main" {
  name    = "my-project"
  project = "my-project"
}

resource "sonarqube_new_code_periods" "main_branch" {
  project = sonarqube_project.main.project
  branch  = "main"
  type    = "PREVIOUS_VERSION"
}
