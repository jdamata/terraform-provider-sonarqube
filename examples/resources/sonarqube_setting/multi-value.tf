resource "sonarqube_setting" "multi_value_setting" {
  key    = "sonar.global.exclusions"
  values = ["foo", "bar/**/*.*"]
}
