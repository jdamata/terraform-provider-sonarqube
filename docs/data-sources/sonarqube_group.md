# Data Source: sonarqube_group

Use this data source to get a Sonarqube Group resource

## Example usage

```terraform
data "sonarqube_group" "group" {
  name = "terraform-test"
}
```


## Argument Reference

The following arguments are supported:

- name - (Required) The name of the Group

## Attributes Reference

The following attributes are exported:

- id - The ID of the Group.
- name - The name of the group.
- description - The group description.
