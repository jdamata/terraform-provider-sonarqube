package sonarqube

/*
 * Make sure the fields are public (First letter Uppercase),
 * otherwise the JSON serialization will fail.
 */

// Version for unmarshalling
type Version struct {
	PageIndex int64 `json:"pageIndex"`
}

// GetQualityGate for unmarshalling response body of quality gate get
type GetQualityGate struct {
	ID         int64                                `json:"id"`
	Name       string                               `json:"name"`
	Conditions []CreateQualityGateConditionResponse `json:"conditions"`
	IsBuiltIn  bool                                 `json:"isBuiltIn"`
	Actions    QualityGateActions                   `json:"actions"`
}

// QualityGateActions used in GetQualityGate
type QualityGateActions struct {
	Rename            bool `json:"rename"`
	SetAsDefault      bool `json:"setAsDefault"`
	Copy              bool `json:"copy"`
	AssociateProjects bool `json:"associateProjects"`
	Delete            bool `json:"delete"`
	ManageConditions  bool `json:"manageConditions"`
}

// CreateQualityGateResponse for unmarshalling response body of quality gate creation
type CreateQualityGateResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// CreateQualityGateConditionResponse for unmarshalling response body of condition creation
type CreateQualityGateConditionResponse struct {
	ID      int64  `json:"id"`
	Metric  string `json:"metric"`
	OP      string `json:"op"`
	Error   string `json:"error"`
	Warning string `json:"warning"`
}

// CreateProjectResponse for unmarshalling response body of project creation
type CreateProjectResponse struct {
	Project Project `json:"project"`
}

// Project used in CreateProjectResponse
type Project struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	Qualifier string `json:"qualifier"`
}

// GetProject for unmarshalling response body from getting project details
type GetProject struct {
	Paging     Paging              `json:"paging"`
	Components []ProjectComponents `json:"components"`
}

// CreateGroupResponse for unmarshalling response body of group creation
type CreateGroupResponse struct {
	Group Group `json:"group"`
}

// Group used in CreateGroupResponse
type Group struct {
	ID           int    `json:"id"`
	Organization string `json:"organization"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	MembersCount int    `json:"membersCount"`
	IsDefault    bool   `json:"default"`
}

// GetGroup for unmarshalling response body from getting group details
type GetGroup struct {
	Paging Paging  `json:"paging"`
	Groups []Group `json:"groups"`
}

// Paging used in /search API endpoints
type Paging struct {
	PageIndex int64 `json:"pageIndex"`
	PageSize  int64 `json:"pageSize"`
	Total     int64 `json:"total"`
}

// ProjectComponents used in GetProject
type ProjectComponents struct {
	Organization     string `json:"organization"`
	Key              string `json:"key"`
	Name             string `json:"name"`
	Qualifier        string `json:"qualifier"`
	Visibility       string `json:"visibility"`
	LastAnalysisDate string `json:"lastAnalysisDate"`
	Revision         string `json:"revision"`
}

// GetQualityGateAssociation for unmarshalling response body from getting quality gate association
type GetQualityGateAssociation struct {
	Paging  Paging                              `json:"paging"`
	Results []GetQualityGateAssociationProjects `json:"results"`
}

// GetQualityGateAssociationProjects used in GetQualityGateAssociation
type GetQualityGateAssociationProjects struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	Selected bool   `json:"selected"`
}

// GetInstalledPlugins for unmarshalling response body from geting installed plugins
type GetInstalledPlugins struct {
	Plugins []Plugin `json:"plugins"`
}

// Plugin used in GetInstalledPlugins
type Plugin struct {
	Key                string `json:"key"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Version            string `json:"version"`
	License            string `json:"license"`
	OrganizationName   string `json:"organizationName"`
	OrganizationURL    string `json:"organizationUrl"`
	EditionBundled     bool   `json:"editionBundled"`
	HomepageURL        string `json:"homepageUrl"`
	IssueTrackerURL    string `json:"issueTrackerUrl"`
	Filename           string `json:"filename"`
	Hash               string `json:"hash"`
	SonarLintSupported bool   `json:"sonarLintSupported"`
	DocumentationPath  bool   `json:"documentationPath"`
	UpdatedAt          int    `json:"updatedAt"`
}
