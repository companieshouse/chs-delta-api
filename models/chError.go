package models

import "fmt"

// CHError is a struct representation of the CH Error object.
type CHError struct {
	Error        string                 `json:"error"`
	ErrorValues  map[string]interface{} `json:"error_values"`
	Location     string                 `json:"location"`
	LocationType string                 `json:"location_type"`
	Type         string                 `json:"type"`
}

// String provides a formatted string of a CHError.
func (che CHError) String() string {
	return fmt.Sprintf("{Error:%s, ErrorValues:%s, Location:%s, LocationType:%s, Type:%s}",
		che.Error, che.ErrorValues, che.Location, che.LocationType, che.Type)
}