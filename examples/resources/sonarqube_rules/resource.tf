resource "sonarqube_rule" "rule" {
  key                  = "rule_test"
  markdown_description = "Description of my rule"
  name                 = "rule test"
  params               = "FilePattern=**/pom.xml"
  prevent_reactivation = "false"
  severity             = "CRITICAL"
  status               = "READY"
  template_key         = "xml:XPathCheck"
  type                 = "VULNERABILITY"
}
