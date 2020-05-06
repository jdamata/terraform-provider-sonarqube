package main

/*
 * Make sure the fields are public (First letter Uppercase),
 * otherwise the JSON serialization will fail.
 */

type QualityGate struct {
	Name         string `json:"name"`
	Organization string `json:"organization"`
}
