package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProjectLink struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
	Type string `json:"type"`
}

type CreateProjectLinkResponse struct {
	Link ProjectLink `json:"link"`
}

type SearchProjectLinksResponse struct {
	Links []ProjectLink `json:"links"`
}

func resourceSonarqubeProjectLinks() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Sonarqube Project Links resource. This can be used to manage Sonarqube project links.",
		Create:      resourceSonarqubeProjectLinksCreate,
		Read:        resourceSonarqubeProjectLinksRead,
		Delete:      resourceSonarqubeProjectLinksDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarqubeProjectLinksImport,
		},

		Schema: map[string]*schema.Schema{
			"project_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The key of the project to associate the link with.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the link.",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The URL of the link.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the link, assigned by SonarQube.",
			},
		},
	}
}

func resourceSonarqubeProjectLinksCreate(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/project_links/create"

	sonarQubeURL.RawQuery = url.Values{
		"projectKey": []string{d.Get("project_key").(string)},
		"name":       []string{d.Get("name").(string)},
		"url":        []string{d.Get("url").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeProjectLinksCreate",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeProjectLinksCreate: Failed to call %s: %+v", sonarQubeURL.Path, err)
	}
	defer resp.Body.Close()

	linkResponse := CreateProjectLinkResponse{}
	err = json.NewDecoder(resp.Body).Decode(&linkResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeProjectLinksCreate: Failed to decode json into struct: %+v", err)
	}

	d.SetId(linkResponse.Link.Id)

	return resourceSonarqubeProjectLinksRead(d, m)
}

func resourceSonarqubeProjectLinksRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/project_links/search"

	sonarQubeURL.RawQuery = url.Values{
		"projectKey": []string{d.Get("project_key").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarQubeURL.String(),
		http.StatusOK,
		"resourceSonarqubeProjectLinksRead",
	)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("resourceSonarqubeProjectLinksRead: Failed to call %s: %+v", sonarQubeURL.Path, err)
	}
	defer resp.Body.Close()

	linksResponse := SearchProjectLinksResponse{}
	err = json.NewDecoder(resp.Body).Decode(&linksResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeProjectLinksRead: Failed to decode json into struct: %+v", err)
	}

	for _, link := range linksResponse.Links {
		if link.Id == d.Id() {
			if err := d.Set("name", link.Name); err != nil {
				return fmt.Errorf("resourceSonarqubeProjectLinksRead: Failed to set name: %+v", err)
			}
			if err := d.Set("url", link.Url); err != nil {
				return fmt.Errorf("resourceSonarqubeProjectLinksRead: Failed to set url: %+v", err)
			}
			if err := d.Set("type", link.Type); err != nil {
				return fmt.Errorf("resourceSonarqubeProjectLinksRead: Failed to set type: %+v", err)
			}
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceSonarqubeProjectLinksDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/project_links/delete"

	sonarQubeURL.RawQuery = url.Values{
		"id": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarQubeURL.String(),
		http.StatusNoContent,
		"resourceSonarqubeProjectLinksDelete",
	)
	if err != nil {
		return fmt.Errorf("resourceSonarqubeProjectLinksDelete: Failed to delete project link: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarqubeProjectLinksImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// import id in format {id}/{projectKey}
	importIdComponents := strings.SplitN(d.Id(), "/", 2)

	if len(importIdComponents) != 2 {
		return nil, fmt.Errorf("resourceSonarqubeProjectLinksImport: Import id '%s' is not in format {id}/{projectKey}", d.Id())
	}

	d.SetId(importIdComponents[0])
	if err := d.Set("project_key", importIdComponents[1]); err != nil {
		return nil, fmt.Errorf("resourceSonarqubeProjectLinksImport: Failed to set project_key: %+v", err)
	}

	if err := resourceSonarqubeProjectLinksRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
