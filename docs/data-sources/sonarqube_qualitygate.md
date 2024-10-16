---
page_title: "sonarqube_qualitygate Data Source - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Use this data source to get a Sonarqube qualitygate resource
---

# sonarqube_qualitygate (Data Source)

Use this data source to get a Sonarqube qualitygate resource

## Example usage

```terraform
data "sonarqube_qualitygate" "main" {
  name     = "example"
}
```

## Argument Reference

The following arguments are supported:

- name - (Required) The name of the quality gate

## Attributes Reference

The following attributes are exported:

- name - The name of the Quality Gate.
- copy_from - Origin of Quality Gate.
- is_default - Quality Gate default.
- condition - List of Quality Gate conditions.
