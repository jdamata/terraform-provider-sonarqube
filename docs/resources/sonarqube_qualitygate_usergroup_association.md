# sonarqube_qualitygate_usergroup_association

Provides a Sonarqube Quality Gate Usergroup association resource. This can be used to associate a Quality Gate to an User or to a Group.
The feature is available on SonarQube 9.2 or newer.

## Example: create a quality gate group association
```terraform
resource "sonarqube_qualitygate" "main" {
    name = "my_qualitygate"

    condition {
        metric    = "new_coverage"
        op        = "LT"
        threshold = "30"
    }
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

    condition {
        metric    = "new_coverage"
        op        = "LT"
        threshold = "30"
    }
}

resource "sonarqube_user" "qa_user" {
  login_name = "qa-user"
  name       = "qa-user"
  password   = "secret-sauce37!"
}

resource "sonarqube_qualitygate_usergroup_association" "main" {
    gatename   = sonarqube_qualitygate.main.id
    login_name = sonarqube_user.qa_user.id
}
```

## Argument Reference
The following arguments are supported:

- `gatename` - (Required) The name of the Quality Gate
- `group_name` - (Optional) The name of the Group to associate. Either `group_name` or `login_name` should be provided.
- `login_name` - (Optional) The name of the User to associate. Either `group_name` or `login_name` should be provided.

## Attributes Reference

The following attributes are exported:

- `id` - A quality user/group composed ID for the association entry.

## Import

Import is not required because SonarQube overwrites existing associations.