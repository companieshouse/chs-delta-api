package validation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

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

// Used for unit testing and mocking calls to external functions/methods.
var (
	callFilepathAbs                  = filepath.Abs
	callNewRouter                    = router.NewRouter
	callOpenApiFilterValidateRequest = openapi3filter.ValidateRequest
	callGetCHErrors                  = getCHErrors
	callFindRoute                    = findRoute
	callGetSchema					 = getSchema
	callOnce	 sync.Once
)

// CHValidator provides an interface to interact with the CH Validator.
type CHValidator interface {
	ValidateRequestAgainstOpenApiSpec(httpReq *http.Request, openApiSpec, contextId string) ([]byte, error)
}

// CHValidatorImpl is a concrete implementation of the CHValidator interface.
type CHValidatorImpl struct {
	doc 	*openapi3.T
}

// NewCHValidator returns a new CHValidator implementation.
func NewCHValidator() CHValidator {
	return &CHValidatorImpl{}
}

// ValidateRequestAgainstOpenApiSpec takes a request and an openAPI spec location (string relative path) and uses the
// spec to validate the provided request. If any validation errors are found, then they are formatted and returned to the
// caller. If any errors are encountered while attempting to validate, they are handled and also returned to the caller.
func (chv *CHValidatorImpl) ValidateRequestAgainstOpenApiSpec(httpReq *http.Request, openApiSpec, contextId string) ([]byte, error) {

	// Get the Open API 3 validation schema location.
	ctx := context.Background()

	err := callGetSchema(ctx, openApiSpec, contextId, chv)
	if err != nil {
		return nil, err
	}

	// Initialise router to later retrieve routes to validate against.
	r, err := callNewRouter(chv.doc)
	if err != nil {
		log.ErrorC(contextId, err, log.Data{config.MessageKey: "error occurred while initialising router for validation"})
		return nil, err
	}

	// Find routes using the given http request.
	route, pathParams, err := callFindRoute(r, httpReq)
	if err != nil {
		log.ErrorC(contextId, err, log.Data{config.MessageKey: "error occurred while finding routes for given http request"})
		return nil, err
	}

	// Define a set of Options for the request validator to use. Enable MultiError's to be returned so that all
	// found errors are returned at once.
	opts := &openapi3filter.Options{
		MultiError: true,
		// Set AuthenticationFunc to a Noop function as ERIC will handle missing / malformed security elements.
		AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
	}

	// Validate request parameters
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
		Options:    opts,
	}

	// Switch off the addition of schema error details to the returned error. This stops the OpenApi schema being added to errors.
	openapi3.SchemaErrorDetailsDisabled = true

	log.InfoC(contextId, "Validating request using: ", log.Data{config.OpenApiSpecKey: openApiSpec})
	if err := callOpenApiFilterValidateRequest(ctx, requestValidationInput); err != nil {
		// If errors are found in the request format them and return them.
		log.InfoC(contextId, "Request validated. Errors found.", nil)
		return callGetCHErrors(contextId, err), nil
	}

	// If we reach this point, then no validation errors were found.
	log.InfoC(contextId, "Request validated. No errors were found.", nil)
	return nil, nil
}

func getCHErrors(contextId string, err error) []byte {

	// Create an array of CHError to be returned to the user.
	errorsArr := make([]models.CHError, 0)

	// If we are given a MultiError, then range over it to extract the inner errors for further processing.
	var mea openapi3.MultiError
	if ok := errors.As(err, &mea); ok {
		errorsArr = handleMultiError(contextId, &mea, errorsArr)
	}

	// Marshal the built up array and return it.
	mr, err := json.Marshal(errorsArr)
	if err != nil {
		log.ErrorC(contextId, err, log.Data{config.MessageKey: "error occurred while formatting CHError array into JSON object"})
		return nil
	}

	return mr
}

func handleMultiError(contextId string, mea *openapi3.MultiError, errsArray []models.CHError) []models.CHError {

	// err is a *MultiError, and mea is set to the error's value, so range over it to grab each error.
	for _, e := range *mea {
		// Begin comparing the retrieved error (e) to possible returned error types.

		// If we have a request error, pass re into the handleRequestError func to get a formatted CHError back.
		var re *openapi3filter.RequestError
		if ok := errors.As(e, &re); ok {
			errsArray = handleRequestError(contextId, re, errsArray)
		}

		// If we have a schema error, pass se into the handleSchemaError func to get a formatted CHError back.
		var se *openapi3.SchemaError
		if ok := errors.As(e, &se); ok {
			errsArray = append(errsArray, handleSchemaError(se))
		}

		// Between ERIC and security middleware we should never receive one of these errors.
		var sre *openapi3filter.SecurityRequirementsError
		if ok := errors.As(e, &sre); ok {
			log.InfoC(contextId, "Encountered unexpected security error")
		}
	}

	// Return the populated errsArray.
	return errsArray
}

func handleRequestError(contextId string, re *openapi3filter.RequestError, errsArray []models.CHError) []models.CHError {

	// It is possible that the RequestError contains a MultiError if more than 1 error has been retuned inside of it.
	var mea openapi3.MultiError
	if ok := errors.As(re.Err, &mea); ok {
		return handleMultiError(contextId, &mea, errsArray)
	}

	// If it isn't a MultiError then we can begin the extract the error contents straight away.
	var se *openapi3.SchemaError
	if ok := errors.As(re.Err, &se); ok {
		errsArray = append(errsArray, handleSchemaError(se))
	}

	// If it is neither a MultiError or a SchemaError then check for a ParseError (malformed JSON).
	var pe *openapi3filter.ParseError
	if ok := errors.As(re.Err, &pe); ok {
		errsArray = append(errsArray, handleParseError(pe))
	}

	// The error could be that we are missing the request body entirely, so account for this too.
	if ok := errors.Is(re.Err, openapi3filter.ErrInvalidRequired); ok {
		errsArray = append(errsArray, models.CHError{
			Error:        re.Err.Error(),
			ErrorValues:  nil,
			Location:     "request-body",
			LocationType: jsonPath,
			Type:         chValidationType,
		})
	}

	// Return populated errsArray with new errors added.
	return errsArray
}

func handleSchemaError(se *openapi3.SchemaError) models.CHError {

	reason := strings.Replace(se.Reason, "\"", "'", -1)
	path := strings.Join(se.JSONPointer(), ".")
	fieldName := se.JSONPointer()[len(se.JSONPointer())-1]

	// Switch over validation error for fieldValue to replace required with an empty string. Without this the
	// error simply returns nothing when a required error is found, as it returns what the user gave (nothing).
	fieldValue := ""
	switch se.SchemaField {
	case "required":
		fieldValue = ""
	default:
		fieldValue = fmt.Sprintf("%v", se.Value)
	}

	// Construct and return a CHError using the extracted data.
	return models.CHError{
		Error:        reason,
		ErrorValues:  map[string]interface{}{fieldName: fieldValue},
		Location:     path,
		LocationType: jsonPath,
		Type:         chValidationType,
	}
}

func handleParseError(pe *openapi3filter.ParseError) models.CHError {

	// Construct and return a CHError to handle a ParseError.
	return models.CHError{
		Error:        pe.Cause.Error(),
		ErrorValues:  map[string]interface{}{},
		Location:     "request-body",
		LocationType: jsonPath,
		Type:         chValidationType,
	}
}

// findRoute is used to add an abstraction layer for unit testing. Allowing us to mock the returns for external methods.
func findRoute(r routers.Router, req *http.Request) (route *routers.Route, pathParams map[string]string, err error) {
	return r.FindRoute(req)
}

func getSchema(ctx context.Context, openApiSpec, contextId string, chv *CHValidatorImpl) (err error) {
	callOnce.Do(func (){
		chv.doc , err = loadSchemaFromFile(ctx, openApiSpec, contextId)
		if err != nil {
			log.ErrorC(contextId, err, log.Data{config.OpenApiSpecKey: openApiSpec, config.MessageKey: "unable to open Open API spec"})
		} else {
			if err = chv.doc.Validate(ctx); err != nil {
				log.ErrorC(contextId, err, log.Data{config.MessageKey: "error occurred while trying to call kin-openAPI validation method"})
			}
		}
	})

	return err
}

func loadSchemaFromFile(ctx context.Context, openApiSpec, contextId string) (*openapi3.T, error){
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	abs, err := callFilepathAbs(openApiSpec)
	if err != nil {
		log.ErrorC(contextId, err, log.Data{config.MessageKey: "error occurred while retrieving absolute path of validation schema file"})
		return nil, err
	}
	log.InfoC(contextId, fmt.Sprintf("Retrieved absolute path of validation schema "), log.Data{config.SchemaAbsolutePathKey: abs})

	// Load the validation schema.
	return loader.LoadFromFile(abs)
}