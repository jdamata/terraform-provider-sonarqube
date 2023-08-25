# sonarqube_settings

Provides a Sonarqube Settings resource. This can be used to manage Sonarqube settings.

## Example: create a setting with a single value

```terraform
resource "sonarqube_setting" "single_setting" {
  key   = "sonar.demo"
  value = "sonarqube@example.org"
}

```
## Example: create a setting associated to a component with a single value

```terraform
resource "sonarqube_setting" "single_setting" {
  key       = "sonar.demo"
  value     = "sonarqube@example.org"
  component = "sonar.component"
}

```

## Example: create a setting with multiple values
```terraform
resource "sonarqube_setting" "multi_value_setting" {
  key   = "sonar.global.exclusions"
  values = ["foo", "bar/**/*.*"]
}
```
## Example: create a setting with multiple field values
```terraform
resource "sonarqube_setting" "multi_field_setting" {
  key   = "sonar.issue.ignore.multicriteria"
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
```
## Argument Reference

The following arguments are supported

- key - (Required) Setting key
- value - (Optional) Single valued setting value
- values - (Optional) Multi-valued setting values
- field_values - (Optional) Multi-field setting values
- component - (Optional) Component key. Only keys for projects, applications, portfolios or subportfolios are accepted.

One of value, values, field_values _must_ be supplied
## Attribute Reference

The following attributes are exported:

- key - The Key of the Setting.
