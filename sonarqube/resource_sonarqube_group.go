package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Returns the resource represented by this file.
func resourceSonarqubeGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarqubeGroupCreate,
		Read:   resourceSonarqubeGroupRead,
		Update: resourceSonarqubeGroupUpdate,
		Delete: resourceSonarqubeGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeGroupImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSonarqubeGroupCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/user_groups/create"
	sonarQubeURL.RawQuery = url.Values{
		"name":        []string{d.Get("name").(string)},
		"description": []string{d.Get("description").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeGroupCreate",
	)
	if err != nil {
		return fmt.Errorf("Error creating Sonarqube group: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	groupResponse := CreateGroupResponse{}
	err = json.NewDecoder(resp.Body).Decode(&groupResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeGroupRead: Failed to decode json into struct: %+v", err)
	}

	d.SetId(strconv.Itoa(groupResponse.Group.ID))
	return resourceSonarqubeGroupRead(d, m)
}

func resourceSonarqubeGroupRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/user_groups/search"
	sonarQubeURL.RawQuery = url.Values{
		"q": []string{d.Get("name").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeGroupRead",
	)
	if err != nil {
		return fmt.Errorf("Error reading Sonarqube group: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	groupReadResponse := GetGroup{}
	err = json.NewDecoder(resp.Body).Decode(&groupReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeGroupRead: Failed to decode json into struct: %+v", err)
	}

	// Loop over all groups to see if the group we need exists.
	for _, value := range groupReadResponse.Groups {
		if d.Id() == strconv.Itoa(value.ID) {
			// If it does, set the values of that group
			d.SetId(strconv.Itoa(value.ID))
			d.Set("name", value.Name)
			d.Set("description", value.Description)
		} else {
			d.SetId("")
		}
	}

	return nil
}

func resourceSonarqubeGroupUpdate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/user_groups/update"

	rawQuery := url.Values{
		"id": []string{d.Id()},
	}

	if _, ok := d.GetOk("description"); ok {
		rawQuery.Add("description", d.Get("description").(string))
	} else {
		rawQuery.Add("description", "")
	}

	sonarQubeURL.RawQuery = rawQuery.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeGroupUpdate",
	)
	if err != nil {
		return fmt.Errorf("Error updating Sonarqube group: %+v", err)
	}
	defer resp.Body.Close()

	return resourceSonarqubeGroupRead(d, m)
}

func resourceSonarqubeGroupDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = "api/user_groups/delete"
	sonarQubeURL.RawQuery = url.Values{
		"id": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeGroupDelete",
	)
	if err != nil {
		return fmt.Errorf("Error deleting Sonarqube group: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarqubeGroupRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
