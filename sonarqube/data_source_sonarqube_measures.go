package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MeasuresComponent struct {
	Key      string    `json:"key"`
	Measures []Measure `json:"measures"`
}

type Measure struct {
	Metric    string `json:"metric"`
	Value     string `json:"value"`
	BestValue bool   `json:"bestValue"`
}

type GetMeasuresResponse struct {
	Component MeasuresComponent `json:"component"`
}

func dataSourceSonarqubeMeasures() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Sonarqube Measures data source. Use this to read metric measures for a component.",
		Read:        dataSourceSonarqubeMeasuresRead,

		Schema: map[string]*schema.Schema{
			"component": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The key of the component (project) to retrieve measures for.",
			},
			"metric_keys": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of metric keys to retrieve.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"measures": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of measures returned by SonarQube.",
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func dataSourceSonarqubeMeasuresRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/measures/component"

	component := d.Get("component").(string)

	metricKeysRaw := d.Get("metric_keys").([]interface{})
	metricKeysList := make([]string, 0, len(metricKeysRaw))
	for _, k := range metricKeysRaw {
		metricKeysList = append(metricKeysList, k.(string))
	}

	sonarQubeURL.RawQuery = url.Values{
		"component":  []string{component},
		"metricKeys": []string{strings.Join(metricKeysList, ",")},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"dataSourceSonarqubeMeasuresRead",
	)
	if err != nil {
		return fmt.Errorf("dataSourceSonarqubeMeasuresRead: Failed to call %s: %+v", sonarQubeURL.Path, err)
	}
	defer resp.Body.Close()

	measuresResponse := GetMeasuresResponse{}
	err = json.NewDecoder(resp.Body).Decode(&measuresResponse)
	if err != nil {
		return fmt.Errorf("dataSourceSonarqubeMeasuresRead: Failed to decode json into struct: %+v", err)
	}

	d.SetId(component)

	measures := make([]interface{}, 0, len(measuresResponse.Component.Measures))
	for _, measure := range measuresResponse.Component.Measures {
		measures = append(measures, map[string]interface{}{
			"metric": measure.Metric,
			"value":  measure.Value,
		})
	}

	if err := d.Set("measures", measures); err != nil {
		return fmt.Errorf("dataSourceSonarqubeMeasuresRead: Failed to set measures: %+v", err)
	}

	return nil
}
