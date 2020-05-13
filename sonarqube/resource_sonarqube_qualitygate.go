package sonarqube

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Returns the resource represented by this file.
func resourceSonarqubeQualityGate() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeQualityGateCreate,
		Read:   resourceSonarqubeQualityGateRead,
		Delete: resourceSonarqubeQualityGateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeQualityGateImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSonarqubeQualityGateCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/create"
	sonarQubeURL.RawQuery = url.Values{
		"name": []string{d.Get("name").(string)},
	}.Encode()

	resp := httpRequestHelper(
		*m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceQualityGateCreate",
	)
	defer resp.Body.Close()

	// Decode response into struct
	qualityGateResponse := CreateQualityGateResponse{}
	err := json.NewDecoder(resp.Body).Decode(&qualityGateResponse)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateCreate: Failed to decode json into struct")
	}

	d.SetId(strconv.FormatInt(qualityGateResponse.ID, 10))
	return nil
}

func resourceSonarqubeQualityGateRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/show"
	sonarQubeURL.RawQuery = url.Values{
		"id": []string{d.Id()},
	}.Encode()

	resp := httpRequestHelper(
		*m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceQualityGateRead",
	)
	defer resp.Body.Close()

	// Decode response into struct
	qualityGateReadResponse := GetQualityGate{}
	err := json.NewDecoder(resp.Body).Decode(&qualityGateReadResponse)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateRead: Failed to decode json into struct")
	}

	d.SetId(strconv.FormatInt(qualityGateReadResponse.ID, 10))
	d.Set("name", qualityGateReadResponse.Name)
	return nil
}

func resourceSonarqubeQualityGateDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/destroy"
	sonarQubeURL.RawQuery = url.Values{
		"id": []string{d.Id()},
	}.Encode()

	resp := httpRequestHelper(
		*m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceQualityGateDelete",
	)
	defer resp.Body.Close()

	// Decode response into struct
	qualityGateReadResponse := GetQualityGate{}
	err := json.NewDecoder(resp.Body).Decode(&qualityGateReadResponse)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateDelete: Failed to decode json into struct")
	}

	return nil
}

func resourceSonarqubeQualityGateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeQualityGateRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
