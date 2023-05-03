# sonarqube_qualitygate

Provides a Sonarqube Quality Gate resource. This can be used to create and manage Sonarqube Quality Gates and their Conditions.

## Example: create a quality gate

```terraform
resource "sonarqube_qualitygate" "main" {
    name = "example"
    is_default = true

    condition {
        metric    = "new_coverage"
        op        = "LT"
        threshold = "50"
    }

    condition {
        metric    = "vulnerabilities"
        threshold = "10"
        op        = "GT"
    }
}
```

## Argument Reference

The following arguments are supported:

- name - (Required) The name of the Quality Gate to create. Maximum length 100
- copy_from - (Optional) Name of an existing Quality Gate to copy from
- is_default - (Optional) When set to true this will make the added Quality Gate default
- condition - (Required) The definition of a Condition to be used by this Quality Gate as documented in the `condition` block below.  At least one condition is required for a Quality Gate.

A `condition` block supports:

- metric - (Required) Condition metric.

  Only metrics of the following types are allowed:
  - INT
  - MILLISEC
  - RATING
  - WORK_DUR
  - FLOAT
  - PERCENT
  - LEVEL.

  The following metrics are forbidden:
  - alert_status
  - security_hotspots
  - new_security_hotspots
-
- threshold - (Required) Condition error threshold (For ratings: A=1, B=2, C=3, D=4)
- op - (Required) Condition operator. Possible values are: LT and GT

## Attributes Reference

The following attributes are exported:

- name - Name of the Sonarqube Quality Gate
- id - ID of the Sonarqube Quality Gate (Deprecated in SonarQube 8.4 and removed in 10.0 so recommended you do not rely on this)

## Import

Quality Gates can be imported using their numeric value

```terraform
terraform import sonarqube_qualitygate.main 11
```
