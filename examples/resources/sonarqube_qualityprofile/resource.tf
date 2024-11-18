resource "sonarqube_qualityprofile" "main" {
  name       = "example"
  language   = "js"
  is_default = false
  parent     = "sonar way"
}
