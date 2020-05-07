package main

/*
 * Make sure the fields are public (First letter Uppercase),
 * otherwise the JSON serialization will fail.
 */

// QualityGate struct for enconding the request body
type QualityGate struct {
	Name         string `json:"name"`
	Organization string `json:"organization"`
}

// QualityGateResponse for unmarshalling response body
type QualityGateResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
