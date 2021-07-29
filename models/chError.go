package models

// CHError is a struct representation of the CH Error object.
type CHError struct {
	Error        string                 `json:"error"`
	ErrorValues  map[string]interface{} `json:"error_values"`
	Location     string                 `json:"location"`
	LocationType string                 `json:"location_type"`
	Type         string                 `json:"type"`
}
