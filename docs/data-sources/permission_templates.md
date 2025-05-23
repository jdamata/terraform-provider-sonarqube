---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarqube_permission_templates Data Source - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Use this data source to get Sonarqube permission templates resources
---

# sonarqube_permission_templates (Data Source)

Use this data source to get Sonarqube permission templates resources

## Example Usage

```terraform
data "sonarqube_permission_templates" "permission_templates" {

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `search` (String) Search permission templates by name.

### Read-Only

- `id` (String) The ID of this resource.
- `permission_templates` (List of Object) The list of permission templates. (see [below for nested schema](#nestedatt--permission_templates))

<a id="nestedatt--permission_templates"></a>
### Nested Schema for `permission_templates`

Read-Only:

- `description` (String)
- `id` (String)
- `name` (String)
- `project_key_pattern` (String)
