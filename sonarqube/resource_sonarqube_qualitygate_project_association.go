package sonarqube

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Returns the resource represented by this file.
func resourceSonarqubeQualityGateProjectAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeQualityGateProjectAssociationCreate,
		Read:   resourceSonarqubeQualityGateProjectAssociationRead,
		Delete: resourceSonarqubeQualityGateProjectAssociationDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"gateid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"projectkey": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSonarqubeQualityGateProjectAssociationCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/select"
	query := url.Values{
		"gateId":     []string{d.Get("gateid").(string)},
		"projectKey": []string{d.Get("projectkey").(string)},
	}
	sonarQubeURL.RawQuery = query.Encode()

	req, err := http.NewRequest("POST", sonarQubeURL.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeQualityGateProjectAssociationCreate")
		return err
	}
	resp, err := m.(*ProviderConfiguration).httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeQualityGateProjectAssociationCreate")
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}

	id := fmt.Sprintf("%v/%v", d.Get("gateid").(string), d.Get("projectkey").(string))
	d.SetId(id)
	return nil
}

func resourceSonarqubeQualityGateProjectAssociationRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/search"
	query := url.Values{
		"gateId": []string{d.Get("gateid").(string)},
	}
	sonarQubeURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", sonarQubeURL.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeQualityGateProjectAssociationRead")
		return err
	}
	resp, err := m.(*ProviderConfiguration).httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeQualityGateProjectAssociationRead")
		return err
	}

	defer resp.Body.Close()
	log.WithField("status code", resp.StatusCode).Info("Response from server")
	if resp.StatusCode != http.StatusOK {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}

	qualityGateAssociationReadResponse := GetQualityGateAssociation{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateAssociationReadResponse)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeQualityGateProjectAssociationRead")
	}

	// ID is in format <gateid>/<projectkey>. This splits the id into gateid and projectkey
	// EG: "1/my_project" >> ["1", "my_project"]
	idSlice := strings.Split(d.Id(), "/")

	for _, value := range qualityGateAssociationReadResponse.Results {
		if idSlice[1] == value.Key {
			d.Set("gateid", idSlice[0])
			d.Set("projectkey", value.Key)
		}
	}

	return nil
}

func resourceSonarqubeQualityGateProjectAssociationDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/qualitygates/deselect"
	query := url.Values{
		"gateId":     []string{d.Get("gateid").(string)},
		"projectKey": []string{d.Get("projectkey").(string)},
	}
	sonarQubeURL.RawQuery = query.Encode()
	req, err := http.NewRequest("POST", sonarQubeURL.String(), http.NoBody)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeQualityGateProjectAssociationDelete")
		return err
	}
	resp, err := m.(*ProviderConfiguration).httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("resourceSonarqubeQualityGateProjectAssociationDelete")
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		responseBody := getResponseBodyAsString(resp)
		return errors.New(responseBody)
	}

	return nil
}
