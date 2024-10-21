resource "sonarqube_alm_github" "github-alm" {
  app_id         = "12345"
  client_id      = "56789"
  client_secret  = "secret"
  key            = "myalm"
  private_key    = "myprivate_key"
  url            = "https://api.github.com"
  webhook_secret = "mysecret"
}
