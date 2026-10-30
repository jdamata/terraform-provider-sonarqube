package sonarqube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubePortfolio() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a Sonarqube portfolio resource",
		Read:        dataSourceSonarqubePortfolioRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The key of the portfolio",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the portfolio",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the portfolio",
			},
			"qualifier": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "`VW` (portfolios always have this qualifier)",
			},
			"visibility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Portfolio visibility",
			},
			"selection_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "How the Portfolio is populated. Possible values are `NONE`, `MANUAL`, `TAGS`, `REGEXP` or `REST`. [See docs](https://docs.sonarqube.org/9.8/project-administration/managing-portfolios/#populating-portfolios) for how Portfolio population works",
			},
			"branch": { // Only active for TAGS, REGEXP and REST
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Which branch is analyzed",
			},
			"tags": { // Only active for TAGS
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The list of tags used to populate the Portfolio. Only active when `selection_mode` is `TAGS`",
			},
			"regexp": { // Only active for REGEXP
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The regular expression used to populate the portfolio. Only active when `selection_mode` is `REGEXP`",
			},
			"selected_projects": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project key of the project in the portfolio",
						},
						"selected_branches": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The branches of the project in the portfolio",
						},
					},
				},
				Description: "The projects in the portfolio. Only set when `selection_mode` is `MANUAL`",
			},
			"applications": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Keys of the applications in the portfolio.",
			},
		},
	}
}

func dataSourceSonarqubePortfolioRead(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("key").(string))
	return resourceSonarqubePortfolioRead(d, m)
}
