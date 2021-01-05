# sonarqube_qualitygate_condition
Provides a Sonarqube Quality Gate Condition resource. This can be used to create and manage Sonarqube Quality Gate conditions.

## Example: create a quality gate condition
```terraform
resource "sonarqube_qualitygate" "main" {
    name = "my_qualitygate"
}

resource "sonarqube_qualitygate_condition" "main" {
    gatename  = sonarqube_qualitygate.main.id
    metric    = "vulnerabilities"
    threshold = "10"
    op        = "GT"
}
```

## Argument Reference
The following arguments are supported:

- gatename - (Required) The name of the Quality Gate
- metric - (Required) Condition metric. Only metric of the following types are allowed: INT, MILLISEC, RATING, WORK_DUR, FLOAT, PERCENT and LEVEL. Following metrics are forbidden: alert_status, security_hotspots and new_security_hotspots
- threshold - (Required) Condition error threshold
- op - (Required) Condition operator. Possible values are: LT and GT

## Attributes Reference
The following attributes are exported:

- id - ID of the Sonarqube Quality Gate
- metric - Condition metric
- threshold - Condition error threshold
- warning - Condition warning threshold
- op - Condition operator

