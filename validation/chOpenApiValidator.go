package validation

import (
	"context"
	"fmt"
	"github.com/companieshouse/chs.go/log"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	lazyrouter "github.com/getkin/kin-openapi/routers/legacy"
	"net/http"
	"path/filepath"
)

func ValidateRequestAgainstOpenApiSpec(httpReq *http.Request, openApiSpec string) []byte {

	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	abs, _ := filepath.Abs(openApiSpec)
	log.Info(fmt.Sprintf("Absolute path: %s", abs))
	doc, _ := loader.LoadFromFile(abs)
	if doc != nil {
		_ = doc.Validate(ctx)
		router, _ := lazyrouter.NewRouter(doc)

		// Find route
		route, pathParams, _ := router.FindRoute(httpReq)

		opts := &openapi3filter.Options{MultiError: true}

		// Validate request parameters
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    httpReq,
			PathParams: pathParams,
			Route:      route,
			Options:    opts,
		}

		// Switch off the addition of schema error details to the returned error. This stops the OpenApi schema being added to errors
		openapi3.SchemaErrorDetailsDisabled = true

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			return FormatError(err)
		}
	} else {
		log.Error(fmt.Errorf("unable to open Open API spec: %s", openApiSpec), nil)
	}

	return nil
}