package sonarqube

/*
 * Make sure the fields are public (First letter Uppercase),
 * otherwise the JSON serialization will fail.
 */

// QualityGate struct for enconding the request body
type QualityGate struct {
	Name string `json:"name"`
}

// QualityGateResponse for unmarshalling response body
type QualityGateResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
