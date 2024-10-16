---
page_title: "sonarqube_project Data Source - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Use this data source to get a Sonarqube project resource
---

# sonarqube_project (Data Source)

Use this data source to get a Sonarqube project resource

## Example usage

```terraform
data "sonarqube_project" "project" {
  project = "projet-key-id"
}
```

## Argument Reference

The following arguments are supported:

- project - (Required) The project key of the project

## Attributes Reference

The following attributes are exported:

- name - Name of the project
- project - Key of the project
- visibility - Project visibility
