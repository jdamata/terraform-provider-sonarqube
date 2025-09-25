resource "sonarqube_setting" "multi_field_setting" {
  key = "sonar.issue.ignore.multicriteria"
  field_values = [
    {
      "ruleKey" : "foo",
      "resourceKey" : "bar"
    },
    {
      "ruleKey" : "foo2",
      "resourceKey" : "bar2"
    }
  ]
}
