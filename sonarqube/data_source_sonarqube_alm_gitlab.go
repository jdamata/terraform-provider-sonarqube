package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubeAlmGitlab() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a GitLab ALM definition.",
		Read:        dataSourceSonarqubeAlmGitlabRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique key of the GitLab instance setting.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GitLab API URL.",
			},
		},
	}
}

func dataSourceSonarqubeAlmGitlabRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/list_definitions"
	sonarQubeURL.RawQuery = url.Values{}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"dataSourceSonarqubeAlmGitlabRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	almGitlabReadResponse := GetAlmGitlab{}
	if err := json.NewDecoder(resp.Body).Decode(&almGitlabReadResponse); err != nil {
		return fmt.Errorf("dataSourceSonarqubeAlmGitlabRead: Failed to decode json into struct: %+v", err)
	}

	for _, value := range almGitlabReadResponse.Gitlab {
		if d.Get("key").(string) == value.Key {
			d.SetId(value.Key)
			return d.Set("url", value.URL)
		}
	}

	return fmt.Errorf("dataSourceSonarqubeAlmGitlabRead: Failed to find gitlab alm definition: %+v", d.Get("key").(string))
}
