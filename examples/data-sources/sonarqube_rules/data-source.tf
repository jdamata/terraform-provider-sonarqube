data "sonarqube_rules" "critical_code_smells" {
  languages  = ["java"]
  types      = ["CODE_SMELL"]
  severities = ["CRITICAL"]
  statuses   = ["READY"]
}

resource "sonarqube_qualityprofile" "java" {
  name     = "Java strict"
  language = "java"
}

resource "sonarqube_qualityprofile_activate_rule" "critical_code_smells" {
  for_each = { for rule in data.sonarqube_rules.critical_code_smells.rules : rule.key => rule }
  key      = sonarqube_qualityprofile.java.key
  rule     = each.key
  severity = each.value.severity
}
