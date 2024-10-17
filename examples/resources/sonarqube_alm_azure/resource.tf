resource "sonarqube_alm_azure" "az1" {
  key                   = "az1"
  personal_access_token = "my_pat"
  url                   = "https://dev.azure.com/my-org"
}
