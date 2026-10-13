package sonarqube

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type GetApplication struct {
	Application struct {
		Key         string   `json:"key"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Visibility  string   `json:"visibility"`
		Tags        []string `json:"tags"`
		Projects    []struct {
			Key string `json:"key"`
		} `json:"projects"`
	} `json:"application"`
}

func resourceSonarqubeApplication() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Sonarqube Application resource. An application groups several projects into one view. Applications are only available in the Developer edition and above.",
		Create:      resourceSonarqubeApplicationCreate,
		Read:        resourceSonarqubeApplicationRead,
		Update:      resourceSonarqubeApplicationUpdate,
		Delete:      resourceSonarqubeApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key of the application. Changing it creates a new application.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the application.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the application.",
			},
			"visibility": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "Whether the application is `public` or `private`. Defaults to SonarQube's default visibility.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"public", "private"}, false)),
			},
			"projects": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Keys of the projects in the application.",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Tags of the application.",
			},
		},
	}
}

func applicationPost(m interface{}, path string, params url.Values, expectedStatus int, resource string) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + path
	sonarQubeURL.RawQuery = params.Encode()
	resp, err := httpRequestHelper(m.(*ProviderConfiguration).httpClient, "POST", sonarQubeURL.String(), expectedStatus, resource)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func stringSet(v interface{}) []string {
	values := []string{}
	for _, value := range v.(*schema.Set).List() {
		values = append(values, value.(string))
	}
	return values
}

func resourceSonarqubeApplicationCreate(d *schema.ResourceData, m interface{}) error {
	key := d.Get("key").(string)
	params := url.Values{
		"key":         []string{key},
		"name":        []string{d.Get("name").(string)},
		"description": []string{d.Get("description").(string)},
	}
	if visibility, ok := d.GetOk("visibility"); ok {
		params.Set("visibility", visibility.(string))
	}
	if err := applicationPost(m, "/api/applications/create", params, http.StatusOK, "resourceSonarqubeApplicationCreate"); err != nil {
		conf := m.(*ProviderConfiguration)
		if strings.ToLower(conf.sonarQubeEdition) == "community" {
			return fmt.Errorf("applications need the Developer edition or above (running SonarQube %s %s): %w", conf.sonarQubeEdition, conf.sonarQubeVersion, err)
		}
		return err
	}
	d.SetId(key)

	for _, project := range stringSet(d.Get("projects")) {
		if err := applicationPost(m, "/api/applications/add_project", url.Values{"application": {key}, "project": {project}}, http.StatusNoContent, "resourceSonarqubeApplicationCreate"); err != nil {
			return err
		}
	}
	if tags := stringSet(d.Get("tags")); len(tags) > 0 {
		if err := applicationPost(m, "/api/applications/set_tags", url.Values{"application": {key}, "tags": {strings.Join(tags, ",")}}, http.StatusNoContent, "resourceSonarqubeApplicationCreate"); err != nil {
			return err
		}
	}
	return resourceSonarqubeApplicationRead(d, m)
}

func resourceSonarqubeApplicationRead(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/applications/show"
	sonarQubeURL.RawQuery = url.Values{"application": []string{d.Id()}}.Encode()

	resp, err := httpRequestHelper(m.(*ProviderConfiguration).httpClient, "GET", sonarQubeURL.String(), http.StatusOK, "resourceSonarqubeApplicationRead")
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			d.SetId("")
			return nil
		}
		return err
	}
	defer resp.Body.Close()

	app := GetApplication{}
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return fmt.Errorf("resourceSonarqubeApplicationRead: Failed to decode json into struct: %+v", err)
	}
	projects := []string{}
	for _, p := range app.Application.Projects {
		projects = append(projects, p.Key)
	}
	return errors.Join(
		d.Set("key", app.Application.Key),
		d.Set("name", app.Application.Name),
		d.Set("description", app.Application.Description),
		d.Set("visibility", app.Application.Visibility),
		d.Set("projects", projects),
		d.Set("tags", app.Application.Tags),
	)
}

func resourceSonarqubeApplicationUpdate(d *schema.ResourceData, m interface{}) error {
	key := d.Id()
	if d.HasChanges("name", "description") {
		// update clears any field it is not sent, so name and description always go together.
		params := url.Values{"application": {key}, "name": {d.Get("name").(string)}, "description": {d.Get("description").(string)}}
		if err := applicationPost(m, "/api/applications/update", params, http.StatusNoContent, "resourceSonarqubeApplicationUpdate"); err != nil {
			return err
		}
	}
	if d.HasChange("visibility") {
		params := url.Values{"project": {key}, "visibility": {d.Get("visibility").(string)}}
		if err := applicationPost(m, "/api/projects/update_visibility", params, http.StatusNoContent, "resourceSonarqubeApplicationUpdate"); err != nil {
			return err
		}
	}
	if d.HasChange("projects") {
		// add_project fails for a project that is already in the application, so only send the difference.
		oldSet, newSet := d.GetChange("projects")
		for _, project := range stringSet(oldSet.(*schema.Set).Difference(newSet.(*schema.Set))) {
			if err := applicationPost(m, "/api/applications/remove_project", url.Values{"application": {key}, "project": {project}}, http.StatusNoContent, "resourceSonarqubeApplicationUpdate"); err != nil {
				return err
			}
		}
		for _, project := range stringSet(newSet.(*schema.Set).Difference(oldSet.(*schema.Set))) {
			if err := applicationPost(m, "/api/applications/add_project", url.Values{"application": {key}, "project": {project}}, http.StatusNoContent, "resourceSonarqubeApplicationUpdate"); err != nil {
				return err
			}
		}
	}
	if d.HasChange("tags") {
		params := url.Values{"application": {key}, "tags": {strings.Join(stringSet(d.Get("tags")), ",")}}
		if err := applicationPost(m, "/api/applications/set_tags", params, http.StatusNoContent, "resourceSonarqubeApplicationUpdate"); err != nil {
			return err
		}
	}
	return resourceSonarqubeApplicationRead(d, m)
}

func resourceSonarqubeApplicationDelete(d *schema.ResourceData, m interface{}) error {
	sonarQubeURL := m.(*ProviderConfiguration).sonarQubeURL
	sonarQubeURL.Path = strings.TrimSuffix(sonarQubeURL.Path, "/") + "/api/applications/delete"
	sonarQubeURL.RawQuery = url.Values{"application": []string{d.Id()}}.Encode()

	resp, err := httpRequestHelper(m.(*ProviderConfiguration).httpClient, "POST", sonarQubeURL.String(), http.StatusNoContent, "resourceSonarqubeApplicationDelete")
	if err != nil {
		// Already gone, e.g. deleted in the UI: nothing to do.
		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			return nil
		}
		return err
	}
	resp.Body.Close()
	return nil
}
