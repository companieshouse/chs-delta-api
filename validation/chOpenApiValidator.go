package validation

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/companieshouse/chs-delta-api/models"
	"github.com/companieshouse/chs.go/log"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	router "github.com/getkin/kin-openapi/routers/gorillamux"
	"net/http"
	"path/filepath"
	"strings"
)

// Used for unit testing and mocking calls to external functions/methods.
var (
	callFilepathAbs                  = filepath.Abs
	callNewRouter                    = router.NewRouter
	callOpenApiFilterValidateRequest = openapi3filter.ValidateRequest
	callFormatError                  = formatError
)

// CHValidator provides an interface to interact with the CH Validator.
type CHValidator interface {
	ValidateRequestAgainstOpenApiSpec(httpReq *http.Request, openApiSpec string) ([]byte, error)
}

// CHValidatorImpl is a concrete implementation of the CHValidator interface.
type CHValidatorImpl struct {
}

// NewCHValidator returns a new CHValidator implementation.
func NewCHValidator() CHValidator {
	return CHValidatorImpl{}
}

// ValidateRequestAgainstOpenApiSpec takes a request and an openAPI spec location (string relative path) and uses the
// spec to validate the provided request. If any validation errors are found, then they are formatted and returned to the
// caller. If any errors are encountered while attempting to validate, they are handled and also returned to the caller.
func (chv CHValidatorImpl) ValidateRequestAgainstOpenApiSpec(httpReq *http.Request, openApiSpec string) ([]byte, error) {

	// Get the Open API 3 validation schema location.
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	abs, err := callFilepathAbs(openApiSpec)
	if err != nil {
		log.Error(fmt.Errorf("error occured while retrieving absolute path of validation schema file: %s", err))
		return nil, err
	}
	log.Info(fmt.Sprintf("Retrieved absolute path of validation schema: %s", abs))

	// Load the validation schema.
	doc, err := loader.LoadFromFile(abs)
	if err != nil {
		log.Error(fmt.Errorf("unable to open Open API spec: %s", openApiSpec), nil)
		return nil, err
	} else {
		if err := doc.Validate(ctx); err != nil {
			log.Error(fmt.Errorf("error occured while trying to call kin-openAPI validation method: %s", err))
			return nil, err
		}

		// Initialise router to later retrieve routes to validate against.
		r, err := callNewRouter(doc)
		if err != nil {
			log.Error(fmt.Errorf("error occured while initialising router for validation: %s", err))
			return nil, err
		}

		// Find routes using the given http request.
		route, pathParams, err := r.FindRoute(httpReq)
		if err != nil {
			log.Error(fmt.Errorf("error occured while finding routes for given http request: %s", err))
			return nil, err
		}

		// Enable MultiError's to be returned so that if more than one error is found, it will return them all.
		opts := &openapi3filter.Options{MultiError: true}

		// Validate request parameters
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    httpReq,
			PathParams: pathParams,
			Route:      route,
			Options:    opts,
		}

		// Switch off the addition of schema error details to the returned error. This stops the OpenApi schema being added to errors.
		openapi3.SchemaErrorDetailsDisabled = true

		log.Info("Validating request...", nil)
		if err := callOpenApiFilterValidateRequest(ctx, requestValidationInput); err != nil {
			// If errors are found in the request format them and return them.
			return callFormatError(err), nil
		}

		// If we reach this point, then no validation errors were found.
		return nil, nil
	}
}

func formatError(err error) []byte {
	var errorsArr []models.CHError

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
			err := models.CHError{
				Error:        reason,
				ErrorValues:  map[string]interface{}{fieldName: fieldValue},
				Location:     jsonPath,
				LocationType: "json-path",
				Type:         "ch:validation",
			}
			errorsArr = append(errorsArr, err)
		}
	}

	// If no errors were found then we can return nil here.
	if len(errorsArr) == 0 {
		return nil
	}

	// If errors do exist, format the array into a JSON object for better viewing.
	mr, err := json.Marshal(errorsArr)
	if err != nil {
		log.Error(fmt.Errorf("error occured while formatting CHError array into JSON object: %s", err))
		return nil
	}

	return mr
}
