resource "sonarqube_alm_azure" "az1" {
  key                   = "az1"
  personal_access_token = "my_pat"
  url                   = "https://dev.azure.com/my-org"
}

resource "sonarqube_project" "main" {
  name       = "SonarQube"
  project    = "main"
  visibility = "public"
}

resource "sonarqube_azure_binding" "main" {
  alm_setting     = sonarqube_alm_azure.az1.key
  project         = sonarqube_project.main.project
  project_name    = "my_azure_project"
  repository_name = "my_repo"
}
