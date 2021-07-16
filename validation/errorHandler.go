package validation

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"strings"
)

type Errors struct {
	Error string `json:"error"`
	ErrorValues map[string]interface{} `json:"error_values"`
	Location string `json:"location"`
	LocationType string `json:"location_type"`
	Type string `json:"type"`
}

type errorsArray []Errors

func GetErrors () errorsArray {
	return errorsArray{}
}

func FormatError(err error) []byte {
	errorsArr := GetErrors()
	for _, me := range err.(openapi3.MultiError) {
		re := me.(*openapi3filter.RequestError)
		for _, se := range re.Err.(openapi3.MultiError) {
			schemaError := se.(*openapi3.SchemaError)
			reason := strings.Replace(schemaError.Reason, "\"", "'", -1)
			jsonPath := strings.Join(schemaError.JSONPointer(), ".")
			fieldName := schemaError.JSONPointer()[len(schemaError.JSONPointer())-1]
			fieldValue := ""
			switch schemaError.SchemaField {
			case "required":
				fieldValue = ""
			default:
				fieldValue = fmt.Sprintf("%v", schemaError.Value)
			}
			error := Errors{
				Error:        reason,
				ErrorValues:  map[string]interface{}{fieldName: fieldValue},
				Location:     jsonPath,
				LocationType: "json-path",
				Type:         "ch:validation",
			}
			errorsArr = append(errorsArr, error)
		}
	}
	mr, _ := json.Marshal(errorsArr)
	return mr
}