resource "sonarqube_alm_bitbucket" "bitbucket-alm" {
  key                   = "mybitbucket"
  personal_access_token = "my_personal_access_token"
  url                   = "https://bitbucket.example.com"
}
