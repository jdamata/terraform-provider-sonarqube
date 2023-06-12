# Data Source: sonarqube_portfolio

Use this data source to get a Sonarqube portfolio resource

## Example usage

```terraform
data "sonarqube_portfolio" "portfolio" {
  key = "portfolio-key"
}
```

## Argument Reference

The following arguments are supported:

- key - (Required) The key of the portfolio

## Attributes Reference

The following attributes are exported:

- name - Name of the portfolio
- description - Description of the portfolio
- qualifier - `VW` (portfolios always have this qualifier)
- visibility - Portfolio visibility
- selection_mode - How the Portfolio is populated. Possible values are ``NONE``, ``MANUAL``, ``TAGS``, ``REGEXP`` or ``REST``. [See docs](https://docs.sonarqube.org/9.8/project-administration/managing-portfolios/#populating-portfolios) for how Portfolio population works
- branch - Which branch is analyzed
- tags - The list of tags used to populate the Portfolio. Only active when `selection_mode` is `TAGS`
- regexp - The regular expression used to populate the portfolio. Only active when `selection_mode` is `REGEXP`