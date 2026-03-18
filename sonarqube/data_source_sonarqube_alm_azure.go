package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSonarqubeAlmAzure() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an Azure DevOps ALM definition.",
		Read:        dataSourceSonarqubeAlmAzureRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique key of the Azure DevOps instance setting.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Azure API URL.",
			},
		},
	}
}

func dataSourceSonarqubeAlmAzureRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/alm_settings/list_definitions"
	sonarQubeURL.RawQuery = url.Values{}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"dataSourceSonarqubeAlmAzureRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	almAzureReadResponse := GetAlmAzure{}
	if err := json.NewDecoder(resp.Body).Decode(&almAzureReadResponse); err != nil {
		return fmt.Errorf("dataSourceSonarqubeAlmAzureRead: Failed to decode json into struct: %+v", err)
	}

	for _, value := range almAzureReadResponse.Azure {
		if d.Get("key").(string) == value.Key {
			d.SetId(value.Key)
			return d.Set("url", value.URL)
		}
	}

	return fmt.Errorf("dataSourceSonarqubeAlmAzureRead: Failed to find azure alm definition: %+v", d.Get("key").(string))
}
