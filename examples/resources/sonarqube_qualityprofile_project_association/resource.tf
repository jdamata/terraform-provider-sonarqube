resource "sonarqube_qualityprofile" "main" {
  name     = "my_quality_profile"
  language = "js"
}

resource "sonarqube_project" "main" {
  name       = "SonarQube"
  project    = "my_project"
  visibility = "public"
}

resource "sonarqube_qualityprofile_project_association" "main" {
  quality_profile = sonarqube_qualityprofile.main.name
  project         = sonarqube_project.main.project
  language        = "js"
}
