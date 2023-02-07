# sonarqube_qualitygate_usergroup_association
Provides a Sonarqube Quality Gate Usergroup association resource. This can be used to associate a Quality Gate to an User or to a Group

## Example: create a quality gate group association
```terraform
resource "sonarqube_qualitygate" "main" {
    name = "my_qualitygate"
}

resource "sonarqube_group" "qa_team" {
    name        = "QA-Team"
    description = "Quality Assurence Team"
}

resource "sonarqube_qualitygate_usergroup_association" "main" {
    gatename   = sonarqube_qualitygate.main.id
    group_name = sonarqube_group.qa_team.name
}
```

## Example: create a quality gate user association
```terraform
resource "sonarqube_qualitygate" "main" {
    name = "my_qualitygate"
}

resource "sonarqube_user" "qa_user" {
  login_name = "qa-user"
  name       = "qa-user"
  password   = "secret-sauce37!"
}

resource "sonarqube_qualitygate_usergroup_association" "main" {
    gatename   = sonarqube_qualitygate.main.id
    group_name = sonarqube_user.qa_user.id
}
```

## Argument Reference
The following arguments are supported:

- gatename - (Required) The name of the Quality Gate
- group_name - (Optional) The name of the Group to associate. 
- login_name - (Required) The name of the User to associate.

## Attributes Reference

The following attributes are exported:

- id - A randomly generated UUID for the permission entry.

## Import

Importing is not supported for the `sonarqube_qualitygate_usergroup_association` resource.
