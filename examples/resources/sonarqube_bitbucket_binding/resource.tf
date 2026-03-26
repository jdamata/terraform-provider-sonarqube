resource "sonarqube_alm_bitbucket" "bitbucket-alm" {
  key                   = "mybitbucket"
  personal_access_token = "my_personal_access_token"
  url                   = "https://bitbucket.example.com"
}

resource "sonarqube_project" "main" {
  name       = "SonarQube"
  project    = "my_project"
  visibility = "public"
}

resource "sonarqube_bitbucket_binding" "bitbucket-binding" {
  alm_setting = sonarqube_alm_bitbucket.bitbucket-alm.key
  project     = sonarqube_project.main.project
  repository  = "MYREPO"
  slug        = "myrepo"
  monorepo    = "false"
}
