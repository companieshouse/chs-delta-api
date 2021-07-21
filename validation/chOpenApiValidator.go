package validation

import (
	"context"
	"fmt"
	"github.com/companieshouse/chs.go/log"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	router "github.com/getkin/kin-openapi/routers/gorillamux"
	"net/http"
	"path/filepath"
)

func ValidateRequestAgainstOpenApiSpec(httpReq *http.Request, openApiSpec string) []byte {

	// Get the Open API 3 validation schema location.
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	abs, err := filepath.Abs(openApiSpec)
	if err != nil {
		log.Error(fmt.Errorf("error occured while retrieving absolute path of validation schema file: %s", err))
		return nil
	}
	log.Info(fmt.Sprintf("Retrieved absolute path of validation schema: %s", abs))

	// Load the validation schema.
	doc, _ := loader.LoadFromFile(abs)
	if doc != nil {
		if err := doc.Validate(ctx); err != nil {
			log.Error(fmt.Errorf("error occured while trying to call kin-openAPI validation method: %s", err))
			return nil
		}

		// Initialise router to later retrieve routes to validate against.
		r, err := router.NewRouter(doc)
		if err != nil {
			log.Error(fmt.Errorf("error occured while initialising router for validation: %s", err))
			return nil
		}

		// Find routes using the given http request.
		route, pathParams, err := r.FindRoute(httpReq)
		if err != nil {
			log.Error(fmt.Errorf("error occured while finding routes for given http request: %s", err))
			return nil
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
		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			// If errors are found in the request format them and return them.
			return FormatError(err)
		}
	} else {
		log.Error(fmt.Errorf("unable to open Open API spec: %s", openApiSpec), nil)
	}

	return nil
}
