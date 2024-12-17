---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarqube_portfolio Resource - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Provides a Sonarqube Portfolio resource. This can be used to create and manage Sonarqube Portfolio. Note that the SonarQube API for Portfolios is called views
---

# sonarqube_portfolio (Resource)

Provides a Sonarqube Portfolio resource. This can be used to create and manage Sonarqube Portfolio. Note that the SonarQube API for Portfolios is called ``views``

## Example Usage

```terraform
resource "sonarqube_portfolio" "main" {
  key         = "portfolio-key"
  name        = "portfolio-name"
  description = "portfolio-description"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `description` (String) A description of the Portfolio to create
- `key` (String) The key of the Portfolio to create
- `name` (String) The name of the Portfolio to create

### Optional

- `branch` (String) Which branch to analyze. If nothing, or '' is specified, the main branch is used.
- `regexp` (String) A regular expression that is used to match Projects with a matching name OR key. If they match, they are added to the Portfolio
- `selected_projects` (Block Set) A set of projects to add to the portfolio. (see [below for nested schema](#nestedblock--selected_projects))
- `selection_mode` (String) How to populate the Portfolio to create. Possible values are `NONE`, `MANUAL`, `TAGS`, `REGEXP` or `REST`. [See docs](https://docs.sonarqube.org/9.8/project-administration/managing-portfolios/#populating-portfolios) for how Portfolio population works
- `tags` (List of String) List of Project tags to populate the Portfolio from. Only active when `selection_mode` is `TAGS`
- `visibility` (String) Whether the created portfolio should be visible to everyone, or only specific user/groups. If no visibility is specified, the default portfolio visibility will be `public`.

### Read-Only

- `id` (String) The ID of this resource.
- `qualifier` (String)

<a id="nestedblock--selected_projects"></a>
### Nested Schema for `selected_projects`

Required:

- `project_key` (String) The project key of the project to add to the portfolio

Optional:

- `selected_branches` (Set of String) A set of branches for the project to add to the portfolio