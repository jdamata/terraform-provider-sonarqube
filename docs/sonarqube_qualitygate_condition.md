# sonarqube_qualitygate_condition
Provides a Sonarqube Quality Gate Condition resource. This can be used to create and manage Sonarqube Quality Gate conditions.

## Example: create a quality gate
```terraform
provider "sonarqube" {
    user = "admin"
    pass = "admin" 
    url = "http://127.0.0.1:9000"
}

resource "sonarqube_qualitygate" "main" {
    name = "example"
}

resource "sonarqube_qualitygate_condition" "main" {
    error = 10
    gateId = sonarqube_qualitygate.main.id
    metric = "vulnerabilities"
    op = "GT"
}
```

## Argument Reference
The following arguments are supported:

- name - (Required) The name of the Quality Gate to create. Maximum length 100

## Attributes Reference
The following attributes are exported:

- name - Name of the Sonarqube Quality Gate
- id - ID of the Sonarqube Quality Gate

## Import 
Quality Gates can be imported using their numeric value

```terraform
terraform import sonarqube_qualitygate.main 11
```

