resource "sonarqube_user_token" "token" {
  name            = "my-token"
  type            = "GLOBAL_ANALYSIS_TOKEN"
  expiration_date = "2099-01-01"
}

output "global_analysis_token" {
  value = sonarqube_user_token.token.token
}
