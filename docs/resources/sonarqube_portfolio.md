# sonarqube_portfolio
Provides a Sonarqube Portfolio resource. This can be used to create and manage Sonarqube Portfolio. Note that the SonarQube API for Portfolios is called ``views`` 

## Example: create a portfolio
```terraform
resource "sonarqube_portfolio" "main" {
    key         = "portfolio-key"
    name        = "portfolio-name"
    description = "portfolio-description"
}
```

## Argument Reference
The following arguments are supported:

- key - (Required) The key of the Portfolio to create
- name - (Required) The name of the Portfolio to create
- description - (Required) A description of the Portfolio to create
- visibility - (Optional) Whether the created portfolio should be visible to everyone, or only specific user/groups. If no visibility is specified, the default portfolio visibility will be `public`.
- selection_mode - (Optional) How to populate the Portfolio to create. Possible values are ``NONE``, ``MANUAL``, ``TAGS``, ``REGEXP`` or ``REST``. [See docs](https://docs.sonarqube.org/9.8/project-administration/managing-portfolios/#populating-portfolios) for how Portfolio population works
- branch - (Optional) Which branch to analyze. If nothing, or "" is specified, the main branch is used.
- tags - (Optional) List of Project tags to populate the Portfolio from. Only active when `selection_mode` is `TAGS`
- regexp - (Optional) A regular expression that is used to match Projects with a matching name OR key. If they match, they are added to the Portfolio

## Attributes Reference
The following attributes are exported in addition to the arguments above:
- qualifier - (Computed) Key of the portfolio (`VW` for views)

## Import 
Portfolios can be imported using their portfolio key

```terraform
terraform import sonarqube_portfolio.main my_portfolio
```

