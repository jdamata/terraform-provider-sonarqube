resource "sonarqube_alm_gitlab" "gitlab-alm" {
  key                   = "myalm"
  personal_access_token = "my_personal_access_token"
  url                   = "https://gitlab.com/api/v4"
}

resource "sonarqube_project" "main" {
  name       = "SonarQube"
  project    = "my_project"
  visibility = "public"
}

resource "sonarqube_gitlab_binding" "gitlab-binding" {
  alm_setting = sonarqube_alm_gitlab.gitlab-alm.key
  project     = "my_project"
  repository  = "123"
}
