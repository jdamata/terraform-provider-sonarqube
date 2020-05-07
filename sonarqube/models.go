package sonarqube

/*
 * Make sure the fields are public (First letter Uppercase),
 * otherwise the JSON serialization will fail.
 */

// QualityGateResponse for unmarshalling response body
type QualityGateResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// QualityGateConditionResponse for unmarshalling response body
type QualityGateConditionResponse struct {
	ID      int64  `json:"id"`
	Metric  string `json:"metric"`
	OP      string `json:"op"`
	Error   string `json:"error"`
	Warning string `json:"warning"`
}
