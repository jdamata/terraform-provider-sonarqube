# Sonarqube_rules

Provides a Sonarqube Rules resource. this can be used to manage Sonarqube rules.

## Example: create a rule

```terraform
resource "sonarqube_rule" "rule" {
  key = "rule_test"
  markdown_description = "Description of my rule"
  name = "rule test"
  params = "FilePattern=**/pom.xml"
  prevent_reactivation = "false"
  severity = "CRITICAL"
  status = "READY"
  template_key = "xml:XPathCheck"
  type = "VULNERABILITY"
}

```

## Argument Reference

the following arguments are supported

- key - (Required) key of the custom rule should only contain : a-z, 0-9, _ 
- markdown_description - (Required) Rule description
- name - (Required) Rule name
- params - (Optional) Parameters as semi-colon list of =, for example 'params=key1=v1;key2=v2' (Only for custom rule)
- prevent_reactivation - (Optional) If set to true and if the rule has been deactivated (status 'REMOVED'), a status 409 will be returned
  - Possible values - true, false, yes, no
- severity - (Optional) Rule severity
  - Possible values - INFO, MINOR, MAJOR, CRITICAL, BLOCKER
- status - (Optional) Rule status
  - Possible values - BETA, DEPRECATED, READY, REMOVED
  - Default value READY
- template_key - Key of the template rule in order to create a custom rule (mandatory for custom rule)
  - [Example values](https://docs.sonarqube.org/latest/user-guide/rules/#header-4)
- type - (Optional) Rule type
  - Possible values - CODE_SMELL, BUG, VULNERABILITY, SECURITY_HOTSPOT

## Attribute Reference

The following attributes are exported:

- key - the Key of the Rule.
