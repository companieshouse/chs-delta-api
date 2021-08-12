package validation

import (
	"bytes"
	"context"
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	router "github.com/getkin/kin-openapi/routers/gorillamux"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

const (
	requestBody     = `{"dummy" : "request"}`
	apiSpecLocation = "../apispec/api-spec.yml"
	contextId       = "contextId"
)

// TestUnitNewCHValidator asserts that the CHValidator constructor correctly returns a CHValidator.
func TestUnitNewCHValidator(t *testing.T) {
	Convey("When I call to get a new CHValidator", t, func() {
		chv := NewCHValidator()

		Convey("Then I am given a new ChValidator", func() {
			So(chv, ShouldNotBeNil)
		})
	})
}

// TestUnitValidateRequestAgainstOpenApiSpecFailsAbs asserts that when getting the ABS path fails, it handled the error correctly.
func TestUnitValidateRequestAgainstOpenApiSpecFailsAbs(t *testing.T) {
	Convey("When I call to validate a request", t, func() {
		chv := NewCHValidator()

		req := httptest.NewRequest("POST", "/dummy/target", bytes.NewBuffer([]byte(requestBody)))

		errReturned := errors.New("error getting ABS path")
		callFilepathAbs = func(path string) (string, error) {
			return "", errReturned
		}

		valErrs, err := chv.ValidateRequestAgainstOpenApiSpec(req, apiSpecLocation, contextId)

		Convey("Then the failure to get the ABS path to spec is handled correctly", func() {
			So(valErrs, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, errReturned)
		})
	})
}

// TestUnitValidateRequestAgainstOpenApiSpecFailsFileOpen asserts that when calling to load the file using an ABS path fails
// it is handled correctly.
func TestUnitValidateRequestAgainstOpenApiSpecFailsFileOpen(t *testing.T) {
	//callOnce needs to be reset in order to work with running tests collectively
	callOnce = sync.Once{}
	Convey("When I call to validate a request", t, func() {
		chv := NewCHValidator()
		req := httptest.NewRequest("POST", "/dummy/target", bytes.NewBuffer([]byte(requestBody)))

		callFilepathAbs = func(path string) (string, error) {
			return "wrongLocation/toSpec", nil
		}

		valErrs, err := chv.ValidateRequestAgainstOpenApiSpec(req, apiSpecLocation, contextId)

		Convey("Then the failure to get the ABS path to spec is handled correctly", func() {
			So(valErrs, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}

// TestUnitValidateRequestAgainstOpenApiSpecFailsToCreateRouter asserts that all errors are handled correctly when
// creating of the Gorilla MUX router fails.
func TestUnitValidateRequestAgainstOpenApiSpecFailsToCreateRouter(t *testing.T) {
	Convey("When I call to validate a request", t, func() {
		chv := NewCHValidator()
		req := httptest.NewRequest("POST", "/dummy/target", bytes.NewBuffer([]byte(requestBody)))

		callFilepathAbs = func(path string) (string, error) {
			return apiSpecLocation, nil
		}

		callNewRouter = func(doc *openapi3.T) (routers.Router, error) {
			return nil, errors.New("error creating router")
		}

		valErrs, err := chv.ValidateRequestAgainstOpenApiSpec(req, apiSpecLocation, contextId)

		Convey("Then the failure to get the ABS path to spec is handled correctly", func() {
			So(valErrs, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}

// TestUnitValidateRequestAgainstOpenApiSpecNoErrors assets that when no errors occur we get a nil return.
func TestUnitValidateRequestAgainstOpenApiSpecNoErrors(t *testing.T) {
	//callOnce needs to be reset in order to work with running tests collectively
	callOnce = sync.Once{}
	Convey("When I call to validate a request", t, func() {
		chv := NewCHValidator()
		req := httptest.NewRequest("POST", "/dummy/target", bytes.NewBuffer([]byte(requestBody)))

		callFilepathAbs = func(path string) (string, error) {
			return apiSpecLocation, nil
		}

		callNewRouter = router.NewRouter

		callFindRoute = func(r routers.Router, req *http.Request) (route *routers.Route, pathParams map[string]string, err error) {
			return &routers.Route{}, make(map[string]string, 1), nil
		}

		callOpenApiFilterValidateRequest = func(ctx context.Context, input *openapi3filter.RequestValidationInput) error {
			return nil
		}

		valErrs, err := chv.ValidateRequestAgainstOpenApiSpec(req, apiSpecLocation, contextId)

		Convey("Then the failure to get the ABS path to spec is handled correctly", func() {
			So(valErrs, ShouldBeNil)
			So(err, ShouldBeNil)
		})
	})
}

// TestUnitValidateRequestAgainstOpenApiSpecFindsValErrors asserts that when kin-openAPI finds validation errors, they are
// returned as a formatted byte array.
func TestUnitValidateRequestAgainstOpenApiSpecFindsValErrors(t *testing.T) {
	//callOnce needs to be reset in order to work with running tests collectively
	callOnce = sync.Once{}
	Convey("When I call to validate a request", t, func() {
		chv := NewCHValidator()
		req := httptest.NewRequest("POST", "/dummy/delta", bytes.NewBuffer([]byte(requestBody)))

		callFilepathAbs = func(path string) (string, error) {
			return apiSpecLocation, nil
		}
		callNewRouter = router.NewRouter

		callFindRoute = func(r routers.Router, req *http.Request) (route *routers.Route, pathParams map[string]string, err error) {
			return &routers.Route{}, make(map[string]string, 1), nil
		}

		callOpenApiFilterValidateRequest = func(ctx context.Context, input *openapi3filter.RequestValidationInput) error {
			return errors.New("validation error")
		}

		callGetCHErrors = func(contextId string, err error) []byte {
			return []byte("error while validating")
		}

		valErrs, err := chv.ValidateRequestAgainstOpenApiSpec(req, apiSpecLocation, contextId)

		Convey("Then the failure to get the ABS path to spec is handled correctly", func() {
			So(valErrs, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}
