# sonarqube_qualitygate
Provides a Sonarqube Quality Gate resource. This can be used to create and manage Sonarqube Quality Gates.

## Example: create a quality gate
```terraform
provider "sonarqube" {
    user = "admin"
    pass = "admin" 
    url = "http://127.0.0.1:9000"
}

resource "sonarqube_qualitygate" "test" {
    name = "test"
}

output "guality-gate-name" {
    value = sonarqube_qualitygate.test.name
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
terraform import sonarqube_qualitygate.test 11
```

