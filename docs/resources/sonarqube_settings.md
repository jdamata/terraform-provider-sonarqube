# sonarqube_settings

Provides a Sonarqube Settings resource. This can be used to manage Sonarqube settings.

## Example: create a setting

```terraform
resource "sonarqube_setting" "setting" {
  key   = "sonar.demo"
  value = "sonarqube@example.org"
}

```

## Argument Reference

The following arguments are supported

- key - (Required) Setting key
- value - (Required) Setting value

## Attribute Reference

The following attributes are exported:

- key - The Key of the Setting.
