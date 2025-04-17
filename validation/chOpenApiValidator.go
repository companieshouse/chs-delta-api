package validation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs-delta-api/models"
	"github.com/companieshouse/chs.go/log"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	router "github.com/getkin/kin-openapi/routers/gorillamux"
)

const (
	jsonPath         = "json-path"
	chValidationType = "ch:validation"
)

// Variables used for unit testing and mocking external functions/methods.
var (
	callFilepathAbs                  = filepath.Abs
	callNewRouter                    = router.NewRouter
	callOpenApiFilterValidateRequest = openapi3filter.ValidateRequest
	callGetCHErrors                  = getCHErrors
	callFindRoute                    = findRoute
	callGetSchema                    = getSchema
)

// CHValidator defines the interface for the CH Validator.
type CHValidator interface {
	ValidateRequestAgainstOpenApiSpec(httpReq *http.Request, contextId string) ([]byte, error)
}

// CHValidatorImpl is a concrete implementation of the CHValidator interface.
type CHValidatorImpl struct {
	doc         *openapi3.T
	openApiSpec string
}

// NewCHValidator creates a new CHValidator instance.
// It returns a pointer to CHValidatorImpl which is more idiomatic.
func NewCHValidator(openApiSpec string) (CHValidator, error) {

	ctx := context.Background()
	doc, err := callGetSchema(ctx, openApiSpec)
	if err != nil {
		// Failed to retrieve the schema, so return an error.
		return nil, err
	}

	// Successfully created a CHValidator, so return the fully constructed object.
	return &CHValidatorImpl{
		doc:         doc,
		openApiSpec: openApiSpec,
	}, nil
}

// ValidateRequestAgainstOpenApiSpec validates the HTTP request against the provided OpenAPI specification.
// If errors are found, they are formatted and returned as JSON.
func (chv *CHValidatorImpl) ValidateRequestAgainstOpenApiSpec(httpReq *http.Request, contextId string) ([]byte, error) {

	ctx := context.Background()

	// Initialize the router to later retrieve the matching routes.
	r, err := callNewRouter(chv.doc)
	if err != nil {
		log.ErrorC(contextId, err, log.Data{config.MessageKey: "error occurred while initialising router for validation"})
		return nil, err
	}

	// Find the route matching the given HTTP request.
	route, pathParams, err := callFindRoute(r, httpReq)
	if err != nil {
		log.ErrorC(contextId, err, log.Data{config.MessageKey: "error occurred while finding routes for given http request"})
		return nil, err
	}

	// Define options for the request validator. Enable MultiError so that all found errors are returned.
	opts := &openapi3filter.Options{
		MultiError:         true,
		AuthenticationFunc: openapi3filter.NoopAuthenticationFunc, // No-op as external middleware handles security.
	}

	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
		Options:    opts,
	}

	// Disable the addition of schema error details to the returned error to prevent the OpenAPI spec from being exposed.
	openapi3.SchemaErrorDetailsDisabled = true

	log.InfoC(contextId, "Validating request using: ", log.Data{config.OpenApiSpecKey: chv.openApiSpec})
	if err := callOpenApiFilterValidateRequest(ctx, requestValidationInput); err != nil {
		// Validation errors found: format and return them.
		log.InfoC(contextId, "Request validated. Errors found.", nil)
		return callGetCHErrors(contextId, err), nil
	}

	// If no errors were found, return nil.
	log.InfoC(contextId, "Request validated. No errors were found.", nil)
	return nil, nil
}

// getCHErrors formats the validation errors into JSON using CHError.
func getCHErrors(contextId string, err error) []byte {

	// Build up an array of CHError objects.
	errorsArr := make([]models.CHError, 0)

	// If the error is a MultiError, iterate over its inner errors for further processing.
	var mea openapi3.MultiError
	if errors.As(err, &mea) {
		errorsArr = handleMultiError(contextId, &mea, errorsArr)
	} else {
		// Fallback for non-MultiError errors: add a generic validation error.
		errorsArr = append(errorsArr, models.CHError{
			Error:        err.Error(),
			ErrorValues:  nil,
			Location:     "request-body",
			LocationType: jsonPath,
			Type:         chValidationType,
		})
	}

	// Log all errors for debugging purposes.
	var errSB strings.Builder
	for _, e := range errorsArr {
		errSB.WriteString(e.String() + ",")
	}
	log.ErrorC(contextId, errors.New(errSB.String()), log.Data{config.MessageKey: "Logging validation errors: "})

	// Marshal the error array into JSON.
	mr, errMarshal := json.Marshal(errorsArr)
	if errMarshal != nil {
		log.ErrorC(contextId, errMarshal, log.Data{config.MessageKey: "error occurred while formatting CHError array into JSON object"})
		return nil
	}

	return mr
}

// handleMultiError iterates over a MultiError and processes each contained error.
func handleMultiError(contextId string, mea *openapi3.MultiError, errsArray []models.CHError) []models.CHError {

	for _, e := range *mea {
		// Check if the error is a RequestError.
		var re *openapi3filter.RequestError
		if errors.As(e, &re) {
			errsArray = handleRequestError(contextId, re, errsArray)
			continue
		}

		// Check if the error is a SchemaError.
		var se *openapi3.SchemaError
		if errors.As(e, &se) {
			errsArray = append(errsArray, handleSchemaError(se))
			continue
		}

		// Check for SecurityRequirementsError (which should normally not occur).
		var sre *openapi3filter.SecurityRequirementsError
		if errors.As(e, &sre) {
			log.InfoC(contextId, "Encountered unexpected security error", nil)
			continue
		}

		// Fallback for unexpected error types.
		errsArray = append(errsArray, models.CHError{
			Error:        e.Error(),
			ErrorValues:  nil,
			Location:     "unknown",
			LocationType: jsonPath,
			Type:         chValidationType,
		})
	}

	return errsArray
}

// handleRequestError processes RequestErrors and extracts meaningful error details.
func handleRequestError(contextId string, re *openapi3filter.RequestError, errsArray []models.CHError) []models.CHError {

	// If RequestError contains a MultiError, process it.
	var mea openapi3.MultiError
	if errors.As(re.Err, &mea) {
		return handleMultiError(contextId, &mea, errsArray)
	}

	// If the error is a SchemaError, format it.
	var se *openapi3.SchemaError
	if errors.As(re.Err, &se) {
		errsArray = append(errsArray, handleSchemaError(se))
		return errsArray
	}

	// If the error is a ParseError (malformed JSON), format it.
	var pe *openapi3filter.ParseError
	if errors.As(re.Err, &pe) {
		errsArray = append(errsArray, handleParseError(pe))
		return errsArray
	}

	// If a required field is missing.
	if errors.Is(re.Err, openapi3filter.ErrInvalidRequired) {
		errsArray = append(errsArray, models.CHError{
			Error:        re.Err.Error(),
			ErrorValues:  nil,
			Location:     "request-body",
			LocationType: jsonPath,
			Type:         chValidationType,
		})
		return errsArray
	}

	// Fallback – append a generic error.
	errsArray = append(errsArray, models.CHError{
		Error:        re.Err.Error(),
		ErrorValues:  nil,
		Location:     "request-body",
		LocationType: jsonPath,
		Type:         chValidationType,
	})
	return errsArray
}

// handleSchemaError processes a SchemaError and returns a formatted CHError.
func handleSchemaError(se *openapi3.SchemaError) models.CHError {

	// Replace double quotes in the error reason with single quotes.
	reason := strings.Replace(se.Reason, "\"", "'", -1)

	// Get the error path from JSON pointer; if empty, use default values.
	pointerParts := se.JSONPointer()
	path := ""
	fieldName := ""
	if len(pointerParts) > 0 {
		path = strings.Join(pointerParts, ".")
		fieldName = pointerParts[len(pointerParts)-1]
	} else {
		path = "request-body"
		fieldName = "unknown"
	}

	// Handle special case for a required error.
	fieldValue := ""
	if se.SchemaField == "required" {
		fieldValue = ""
	} else {
		fieldValue = fmt.Sprintf("%v", se.Value)
	}

	return models.CHError{
		Error:        reason,
		ErrorValues:  map[string]interface{}{fieldName: fieldValue},
		Location:     path,
		LocationType: jsonPath,
		Type:         chValidationType,
	}
}

// handleParseError processes a ParseError (e.g., malformed JSON) and returns a formatted CHError.
func handleParseError(pe *openapi3filter.ParseError) models.CHError {

	return models.CHError{
		Error:        pe.Cause.Error(),
		ErrorValues:  map[string]interface{}{},
		Location:     "request-body",
		LocationType: jsonPath,
		Type:         chValidationType,
	}
}

// getSchema retrieves and validates the OpenAPI3 specification.
func getSchema(ctx context.Context, openApiSpec string) (*openapi3.T, error) {

	log.Info("Retrieving openAPI3 spec")
	doc, err := loadSchemaFromFile(ctx, openApiSpec)
	if err != nil {
		log.Error(err, log.Data{config.OpenApiSpecKey: openApiSpec, config.MessageKey: "unable to open Open API spec"})
	} else {
		if errValidate := doc.Validate(ctx); errValidate != nil {
			log.Error(errValidate, log.Data{config.MessageKey: "error occurred while validating the Open API spec"})
		}
	}

	return doc, err
}

// loadSchemaFromFile loads the OpenAPI3 schema from the filesystem using an absolute path.
func loadSchemaFromFile(ctx context.Context, openApiSpec string) (*openapi3.T, error) {

	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	abs, err := callFilepathAbs(openApiSpec)
	if err != nil {
		log.Error(err, log.Data{config.MessageKey: "error occurred while retrieving absolute path of validation schema file"})
		return nil, err
	}
	log.Info(fmt.Sprintf("Retrieved absolute path of validation schema"), log.Data{config.SchemaAbsolutePathKey: abs})

	return loader.LoadFromFile(abs)
}

// findRoute provides an abstraction layer to allow for easier unit testing.
// It finds the route that matches the given HTTP request.
func findRoute(r routers.Router, req *http.Request) (route *routers.Route, pathParams map[string]string, err error) {
	return r.FindRoute(req)
}
