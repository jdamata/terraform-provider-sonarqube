resource "sonarqube_rule" "allowed_maven_dependencies" {
  custom_key           = "Only_use_allowed_Maven_dependencies"
  markdown_description = "Description"
  name                 = "Only use allowed Maven dependencies"
  params               = "FilePattern=**/pom.xml"
  severity             = "BLOCKER"
  status               = "READY"
  template_key         = "xml:XPathCheck"
  type                 = "VULNERABILITY"
}

resource "sonarqube_qualityprofile" "xml" {
  name       = "test way - xml"
  language   = "xml"
  is_default = "false"
  parent     = "Sonar way"
}

resource "sonarqube_qualityprofile_activate_rule" "xml_rule" {
  key      = sonarqube_qualityprofile.xml.key # .id also works
  rule     = sonarqube_rule.allowed_maven_dependencies.id
  severity = "BLOCKER"
}
