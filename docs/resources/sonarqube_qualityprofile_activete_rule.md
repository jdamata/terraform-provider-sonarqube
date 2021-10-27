# sonarqube_qualityprofile_activate_rule

Provides a Sonarqube Rules resource. This can be used to manage Sonarqube rules.

## Example: activate a rule

```terraform

resource "sonarqube_rule" "allowed_maven_dependencies" {
  custom_key = "Only_use_allowed_Maven_dependencies"
  markdown_description = "Description"  
  name = "Only use allowed Maven dependencies"
  params = "FilePattern=**/pom.xml"
  severity = "BLOCKER"
  status = "READY"
  template_key = "xml:XPathCheck"
  type = "VULNERABILITY"
}

resource "sonarqube_qualityprofile" "xml" {
    name     = "test way - xml"
    language = "xml"
    is_default = "false"
    parent = "Sonar way"
}

resource "sonarqube_qualityprofile_activate_rule" "xml_rule" {
  key = sonarqube_qualityprofile.xml.key # .id also works
  rule = sonarqube_rule.allowed_maven_dependencies.id
  severity = "BLOCKER"
}

```

## Argument Reference

The following arguments are supported

- key - (Required) Quality Profile key. Can be obtained through api/qualityprofiles/search 
- params - (Optional) Parameters as semi-colon list of key=value. Ignored if parameter reset is true.
- reset - (Optional) Reset severity and parameters of activated rule. Set the values defined on parent profile or from rule default values.
  - Possible values true false yes no (Default false)
- rule - Rule key
- severity - (Optional) Severity. Ignored if parameter reset is true.
  - Possible values - INFO, MINOR, MAJOR, CRITICAL, BLOCKER

## Attribute Reference

The following attributes are exported:

- key - the Key of the Rule.
