# Data Source: sonarqube_rule

Use this data source to get a Sonarqube rule resource

## Example usage

```terraform
data "sonarqube_rule" "rule" {
  key = "squid:forbidSonar"
}
```

## Argument Reference

The following arguments are supported:

- key - (Required) The key of the sonarqube rule. Should be <repo>:<name>. https://next.sonarqube.com/sonarqube/web_api/api/rules?query=api%2Frules%2Fcreate

## Attributes Reference

The following attributes are exported:

- markdown_description - Rule description
- name - Rule name
- severity - Rule severity
- status - Rule status
- template_key - Key of the template rule
- type - Rule type