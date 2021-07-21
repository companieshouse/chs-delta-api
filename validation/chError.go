package validation

import (
	"encoding/json"
	"fmt"
	"github.com/companieshouse/chs.go/log"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"strings"
)

type CHError struct {
	Error        string                 `json:"error"`
	ErrorValues  map[string]interface{} `json:"error_values"`
	Location     string                 `json:"location"`
	LocationType string                 `json:"location_type"`
	Type         string                 `json:"type"`
}

func FormatError(err error) []byte {
	var errorsArr []CHError

	// Range over every MultiError to pull all RequestErrors.
	for _, me := range err.(openapi3.MultiError) {

		// Retrieve RequestErrors and range over them to grab their inner MultiErrors, as these contain the SchemaErrors.
		re := me.(*openapi3filter.RequestError)
		for _, me := range re.Err.(openapi3.MultiError) {

			// Cast to SchemaError so that we can pull out all of the necessary data to build our CH Errors response.
			schemaError := me.(*openapi3.SchemaError)
			reason := strings.Replace(schemaError.Reason, "\"", "'", -1)
			jsonPath := strings.Join(schemaError.JSONPointer(), ".")
			fieldName := schemaError.JSONPointer()[len(schemaError.JSONPointer())-1]

			// Switch over validation error for fieldValue to replace required with an empty string. Without this the
			// error simply returns nothing when a required error is found, as it returns what the user gave (nothing).
			fieldValue := ""
			switch schemaError.SchemaField {
			case "required":
				fieldValue = ""
			default:
				fieldValue = fmt.Sprintf("%v", schemaError.Value)
			}

			// Construct a CHError and append it to the previously created CHError slice.
			err := CHError{
				Error:        reason,
				ErrorValues:  map[string]interface{}{fieldName: fieldValue},
				Location:     jsonPath,
				LocationType: "json-path",
				Type:         "ch:validation",
			}
			errorsArr = append(errorsArr, err)
		}
	}

	// Format the array into a JSON object for better viewing.
	mr, err := json.Marshal(errorsArr)
	if err != nil {
		log.Error(fmt.Errorf("error occured while formatting CHError array into JSON object: %s", err))
		return nil
	}

	return mr
}
