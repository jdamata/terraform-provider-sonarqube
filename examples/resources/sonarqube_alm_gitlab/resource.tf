resource "sonarqube_alm_gitlab" "gitlab-alm" {
  key                   = "myalm"
  personal_access_token = "my_personal_access_token"
  url                   = "https://gitlab.com/api/v4"
}
