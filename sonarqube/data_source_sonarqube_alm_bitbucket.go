package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubeAlmBitbucket() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get a Bitbucket Data Center ALM definition.",
		Read:        dataSourceSonarqubeAlmBitbucketRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique key of the Bitbucket instance setting.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Bitbucket server API URL.",
			},
		},
	}
}

func dataSourceSonarqubeAlmBitbucketRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/list_definitions"
	sonarQubeURL.RawQuery = url.Values{}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"dataSourceSonarqubeAlmBitbucketRead",
	)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	almBitbucketReadResponse := GetAlmBitbucket{}
	if err := json.NewDecoder(resp.Body).Decode(&almBitbucketReadResponse); err != nil {
		return fmt.Errorf("dataSourceSonarqubeAlmBitbucketRead: Failed to decode json into struct: %+v", err)
	}

	for _, value := range almBitbucketReadResponse.Bitbucket {
		if d.Get("key").(string) == value.Key {
			d.SetId(value.Key)
			return d.Set("url", value.URL)
		}
	}

	return fmt.Errorf("dataSourceSonarqubeAlmBitbucketRead: Failed to find bitbucket alm definition: %+v", d.Get("key").(string))
}
