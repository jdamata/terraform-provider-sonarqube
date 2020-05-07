package sonarqube

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Returns the resource represented by this file.
func qualityGateCondition() *schema.Resource {
	return &schema.Resource{
		Create: qualityGateConditionCreate,
		Read:   qualityGateConditionRead,
		Update: qualityGateConditionUpdate,
		Delete: qualityGateConditionDelete,
		Importer: &schema.ResourceImporter{
			State: qualityGateConditionImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"error": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"gateid": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
			},
			"op": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func qualityGateConditionCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func qualityGateConditionRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func qualityGateConditionUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func qualityGateConditionDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func qualityGateConditionImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := qualityGateConditionRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func getQualityGateConditionResponse(resp *http.Response) (QualityGateConditionResponse, error) {
	qualityGateConditionResponse := QualityGateConditionResponse{}
	err := json.NewDecoder(resp.Body).Decode(&qualityGateConditionResponse)
	if err != nil {
		log.WithError(err).Error("getQualityGateConditionResponse")
		return qualityGateConditionResponse, err
	}
	return qualityGateConditionResponse, nil
}
