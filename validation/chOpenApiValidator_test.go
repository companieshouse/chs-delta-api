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
	"net/http/httptest"
	"testing"
)

const (
	requestBody     = `{"dummy" : "request"}`
	apiSpecLocation = "../apispec/api-spec.yml"
)

// TestNewCHValidator asserts that the CHValidator constructor correctly returns a CHValidator.
func TestNewCHValidator(t *testing.T) {
	Convey("When I call to get a new CHValidator", t, func() {
		chv := NewCHValidator()

		Convey("Then I am given a new ChValidator", func() {
			So(chv, ShouldNotBeNil)
		})
	})
}

// TestValidateRequestAgainstOpenApiSpecFailsAbs asserts that when getting the ABS path fails, it handled the error correctly.
func TestValidateRequestAgainstOpenApiSpecFailsAbs(t *testing.T) {
	Convey("When I call to validate a request", t, func() {
		chv := NewCHValidator()

		req := httptest.NewRequest("POST", "/dummy/target", bytes.NewBuffer([]byte(requestBody)))

		errReturned := errors.New("error getting ABS path")
		callFilepathAbs = func(path string) (string, error) {
			return "", errReturned
		}

		valErrs, err := chv.ValidateRequestAgainstOpenApiSpec(req, apiSpecLocation)

		Convey("Then the failure to get the ABS path to spec is handled correctly", func() {
			So(valErrs, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, errReturned)
		})
	})
}

// TestValidateRequestAgainstOpenApiSpecFailsFileOpen asserts that when calling to load the file using an ABS path fails
// it is handled correctly.
func TestValidateRequestAgainstOpenApiSpecFailsFileOpen(t *testing.T) {
	Convey("When I call to validate a request", t, func() {
		chv := NewCHValidator()

		req := httptest.NewRequest("POST", "/dummy/target", bytes.NewBuffer([]byte(requestBody)))

		callFilepathAbs = func(path string) (string, error) {
			return "wrongLocation/toSpec", nil
		}

		valErrs, err := chv.ValidateRequestAgainstOpenApiSpec(req, apiSpecLocation)

		Convey("Then the failure to get the ABS path to spec is handled correctly", func() {
			So(valErrs, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}

// TestValidateRequestAgainstOpenApiSpecFailsToCreateRouter asserts that all errors are handled correctly when
// creating of the Gorilla MUX router fails.
func TestValidateRequestAgainstOpenApiSpecFailsToCreateRouter(t *testing.T) {
	Convey("When I call to validate a request", t, func() {
		chv := NewCHValidator()

		req := httptest.NewRequest("POST", "/dummy/target", bytes.NewBuffer([]byte(requestBody)))

		callFilepathAbs = func(path string) (string, error) {
			return apiSpecLocation, nil
		}

		callNewRouter = func(doc *openapi3.T) (routers.Router, error) {
			return nil, errors.New("error creating router")
		}

		valErrs, err := chv.ValidateRequestAgainstOpenApiSpec(req, apiSpecLocation)

		Convey("Then the failure to get the ABS path to spec is handled correctly", func() {
			So(valErrs, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}

// TestValidateRequestAgainstOpenApiSpecNoErrors assets that when no errors occur we get a nil return.
func TestValidateRequestAgainstOpenApiSpecNoErrors(t *testing.T) {
	Convey("When I call to validate a request", t, func() {
		chv := NewCHValidator()

		// Provide an actual url target to allow Router to correctly be initialised for further unit testing.
		req := httptest.NewRequest("POST", "/delta/officers", bytes.NewBuffer([]byte(requestBody)))

		callFilepathAbs = func(path string) (string, error) {
			return apiSpecLocation, nil
		}

		callNewRouter = router.NewRouter

		callOpenApiFilterValidateRequest = func(ctx context.Context, input *openapi3filter.RequestValidationInput) error {
			return nil
		}

		valErrs, err := chv.ValidateRequestAgainstOpenApiSpec(req, apiSpecLocation)

		Convey("Then the failure to get the ABS path to spec is handled correctly", func() {
			So(valErrs, ShouldBeNil)
			So(err, ShouldBeNil)
		})
	})
}

// TestValidateRequestAgainstOpenApiSpecFindsValErrors asserts that when kin-openAPI finds validation errors, they are
// returned as a formatted byte array.
func TestValidateRequestAgainstOpenApiSpecFindsValErrors(t *testing.T) {
	Convey("When I call to validate a request", t, func() {
		chv := NewCHValidator()

		// Provide an actual url target to allow Router to correctly be initialised for further unit testing.
		req := httptest.NewRequest("POST", "/delta/officers", bytes.NewBuffer([]byte(requestBody)))

		callFilepathAbs = func(path string) (string, error) {
			return apiSpecLocation, nil
		}

		callNewRouter = router.NewRouter

		callOpenApiFilterValidateRequest = openapi3filter.ValidateRequest

		callFormatError = func(err error) []byte {
			return []byte("error while validating")
		}

		valErrs, err := chv.ValidateRequestAgainstOpenApiSpec(req, apiSpecLocation)

		Convey("Then the failure to get the ABS path to spec is handled correctly", func() {
			So(valErrs, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}
