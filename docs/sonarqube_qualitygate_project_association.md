# sonarqube_qualitygate_project_association
Provides a Sonarqube Quality Gate Project association resource. This can be used to associate a Quality Gate to a Project

## Example: create a quality gate
```terraform
resource "sonarqube_qualitygate" "main" {
    name = "my_qualitygate"
}

resource "sonarqube_project" "main" {
    name = "my_project"
}

resource "sonarqube_qualitygate_project_association" "main" {
    gateid = sonarqube_qualitygate.main.id
    projectkey = ""
}
```

## Argument Reference
The following arguments are supported:

- gateid - (Required) The id of the Quality Gate
- metric - (Required) Condition metric. Only metric of the following types are allowed: INT, MILLISEC, RATING, WORK_DUR, FLOAT, PERCENT and LEVEL. Following metrics are forbidden: alert_status, security_hotspots and new_security_hotspots
- error - (Required) Condition error threshold
- op - (Required) Condition operator. Possible values are: LT and GT

## Attributes Reference
The following attributes are exported:

- id - ID of the Sonarqube Quality Gate
- metric - Condition metric
- error - Condition error threshold
- warning - Condition warning threshold
- op - Condition operator

