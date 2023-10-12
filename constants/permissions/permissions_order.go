package permissions

var PERMISSIONS_ORDER_FOR_PROJECT = []string{
	Browse,
	CodeViewer,
	IssueAdmin,
	SecurityHotspotAdmin,
	Admin,
	Scan,
}

var PERMISSIONS_ORDER_GLOBAL = []string{
	Admin,
	QualityGateAdmin, QualityProfileAdmin,
	Scan,
	ProjectCreation,
	ApplicationCreation,
	PortfolioCreation,
}
