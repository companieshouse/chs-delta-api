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
	requestBody = `{
  "officers": [{
    "company_number": "09876543",
    "changed_at": "20176543171003950844",
    "kind": "DIR",
    "internal_id": "3002598737",
    "appointment_date": "20111103",
    "title": "Mr",
    "corporate_ind": "N",
    "surname": "SMITH",
    "forename": "JOHN",
    "middle_name": "Peter",
    "date_of_birth": "19850630",
    "service_address_same_as_registered_address": "Y",
    "usual_residential_address_same_as_registered_address": "Y",
    "secure_director" : "Y",
    "nationality": "British",
    "officer_id" : "1234567890",
    "occupation": "Lawyer",
    "secure": "243",
    "officer_detail_id": "3456251385",
    "officer_role": "Director",
    "usual_residential_country": "United Kingdom",
    "previous_name_array": {
      "previous_surname": "BURCH",
      "previous_forename": "VALERIE JEAN",
      "previous_timestamp": "20091101072217613702"
    },
    "identification": {
      "EEA": {
        "place_registered": "United Kingdom",
        "registration_number": "38298",
        "legal_authority": "Chapter 32",
        "legal_form": "Hong Kong"
      }
    },

    "service_address": {
      "premise": "2pm 0",
      "address_line_1": "tall Passage",
      "address_line_2": "......",
      "locality": "Cardiff",
      "care_of": "",
      "region": "",
      "po_box": "",
      "supplied_company_name": "",
      "country": "United Kingdom",
      "postal_code": "CF2 1B6",
      "usual_country_of_residence": "United Kingdom"
    },
    "usual_residential_address": {
      "premise": "2pm 0",
      "address_line_1": "tall Passage",
      "address_line_2": "",
      "locality": "Cardiff",
      "care_of": "",
      "region": "",
      "po_box": "",
      "supplied_company_name": "",
      "country": "United Kingdom",
      "postal_code": "CF2 1B6",
      "usual_country_of_residence": "United Kingdom"
    }
  }],
  "CreatedTime": "21-JUL-21 11.20.00.000000",
  "delta_at": "20140925171003950844"
}`

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

		req := httptest.NewRequest("POST", "/delta/officers", bytes.NewBuffer([]byte(requestBody)))

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

		req := httptest.NewRequest("POST", "/delta/officers", bytes.NewBuffer([]byte(requestBody)))

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

		req := httptest.NewRequest("POST", "/delta/officers", bytes.NewBuffer([]byte(requestBody)))

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
