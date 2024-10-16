---
page_title: "sonarqube_new_code_periods Resource - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Provides a Sonarqube New Code Periods resource. This can be used to manage Sonarqube New Code Periods.
---

# sonarqube_new_code_periods (Resource)

Provides a Sonarqube New Code Periods resource. This can be used to manage Sonarqube New Code Periods.

## Example: Set the global new code period to a number of days

```terraform
resource "sonarqube_new_code_periods" "code_period" {
  type = "NUMBER_OF_DAYS"
  value = "7"
}
```

## Example: create a project and set its new code period to a reference branch

```terraform
resource "sonarqube_project" "reference" {
  name = "my-project"
}

resource "sonarqube_new_code_periods" "reference" {
  project = sonarqube_project.reference.project
  type = "REFERENCE_BRANCH"
  value = "main"
```

## Argument Reference

The following arguments are supported:

- type - (Required) The kind of new code period to use. Supported values are SPECIFIC_ANALYSIS, PREVIOUS_VERSION, NUMBER_OF_DAYS, or REFERENCE_BRANCH.
- project - (Optional) The key of a project for which the new code period will be configured. Changing this will force a new resource to be created.
- branch - (Optional) The name of a branch of a project for which the new code period will be configured. Changing this will force a new resource to be created. Setting this also requires setting the 'project' argument.
- value - (Optional) The desired value of the new code period. Varies based on the 'type'. For SPECIFIC_ANALYIS, the value must be the UUID of a previous analysis. For NUMBER_OF_DAYS it must be a numeric string. For REFERENCE_BRANCH it should be the name of branch on the project. For PREVIOUS_VERSION it must **not** be set.

## Attributes Reference

The following attributes are exported:

- id - The ID of the Token.

## Import

Import is not supported for this resource.
