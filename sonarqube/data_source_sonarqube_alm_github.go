package sonarqube

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubeAlmGithub() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a GitHub ALM definition.",
		Read:        dataSourceSonarqubeAlmGithubRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique key of the GitHub instance setting.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GitHub API URL.",
			},
			"app_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GitHub App ID.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GitHub App Client ID.",
			},
		},
	}
}

func dataSourceSonarqubeAlmGithubRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/list_definitions"
	sonarQubeURL.RawQuery = url.Values{}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"dataSourceSonarqubeAlmGithubRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	almGithubReadResponse := GetAlmGithub{}
	if err := json.NewDecoder(resp.Body).Decode(&almGithubReadResponse); err != nil {
		return fmt.Errorf("dataSourceSonarqubeAlmGithubRead: Failed to decode json into struct: %+v", err)
	}

	for _, value := range almGithubReadResponse.Github {
		if d.Get("key").(string) == value.Key {
			d.SetId(value.Key)
			errs := []error{}
			errs = append(errs, d.Set("url", value.URL))
			errs = append(errs, d.Set("app_id", value.AppID))
			errs = append(errs, d.Set("client_id", value.ClientID))
			return errors.Join(errs...)
		}
	}

	return fmt.Errorf("dataSourceSonarqubeAlmGithubRead: Failed to find github alm definition: %+v", d.Get("key").(string))
}
