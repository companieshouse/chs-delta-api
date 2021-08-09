package officers

import (
	"bytes"
	"encoding/json"
	"github.com/companieshouse/chs-delta-api/models"
	"github.com/companieshouse/chs-delta-api/validation"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	requestBodiesLocation            = "./request_bodies/"
	okRequestBodyLocation            = requestBodiesLocation + "ok_request_body"
	typeErrorRequestBodyLocation     = requestBodiesLocation + "type_error_request_body"
	requiredErrorRequestBodyLocation = requestBodiesLocation + "required_error_request_body"
	enumErrorRequestBodyLocation     = requestBodiesLocation + "enum_error_request_body"

	responseBodiesLocation                 = "./response_bodies/"
	typeErrorResponseBodyLocation          = responseBodiesLocation + "type_error_response_body"
	requiredErrorResponseBodyLocation      = responseBodiesLocation + "required_error_response_body"
	enumErrorResponseBodyLocation          = responseBodiesLocation + "enum_error_response_body"
	noRequestBodyErrorResponseBodyLocation = responseBodiesLocation + "no_request_body_error_response_body"

	officersEndpoint = "/delta/officers"
	apiSpecLocation  = "../../../apispec/api-spec.yml"

	contextId       = "contextId"
	xRequestId      = "X-Request-Id"
	contentType     = "Content-Type"
	applicationJson = "application/json"
)

// TestOfficerDeltaSchemaNoErrors asserts that when a valid request body is given which matches the schema, then no
// errors are returned.
func TestOfficerDeltaSchemaNoErrors(t *testing.T) {

	Convey("Given I want to test the officers-delta API schema", t, func() {

		okRequestBody := readRequestBody(okRequestBodyLocation)

		r := httptest.NewRequest("POST", officersEndpoint, bytes.NewBuffer(okRequestBody))
		r = setHeaders(r)

		Convey("When I call to validate the request body, providing a valid request", func() {

			chv := validation.NewCHValidator()

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, apiSpecLocation, contextId)

			Convey("Then I am given a nil response as no validation errors are returned", func() {
				So(validationErrs, ShouldBeNil)
			})
		})
	})
}

// TestOfficerDeltaSchemaTypeErrors asserts that when an invalid request body is given with type errors (int provided
// instead of string), then an errors array is returned.
func TestOfficerDeltaSchemaTypeErrors(t *testing.T) {

	Convey("Given I want to test the officers-delta API schema for type assertions", t, func() {

		typeErrorRequestBody := readRequestBody(typeErrorRequestBodyLocation)

		r := httptest.NewRequest("POST", officersEndpoint, bytes.NewBuffer(typeErrorRequestBody))
		r = setHeaders(r)

		Convey("When I call to validate the request body, providing an valid request with type errors", func() {

			chv := validation.NewCHValidator()

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, apiSpecLocation, contextId)

			Convey("Then I am given an errors array response as validation errors have been found", func() {
				typeErrorResponseBody := readRequestBody(typeErrorResponseBodyLocation)
				match := compareActualToExpected(validationErrs, typeErrorResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}

// TestOfficerDeltaSchemaRequiredErrors asserts that when an invalid request body is given with missing mandatory values.
// then an errors array is returned, stating that required values are missing.
func TestOfficerDeltaSchemaRequiredErrors(t *testing.T) {

	Convey("Given I want to test the officers-delta API schema to assert mandatory validation is working correctly", t, func() {

		mandatoryErrorsRequestBody := readRequestBody(requiredErrorRequestBodyLocation)

		r := httptest.NewRequest("POST", officersEndpoint, bytes.NewBuffer(mandatoryErrorsRequestBody))
		r = setHeaders(r)

		Convey("When I call to validate the request body, providing an valid request with missing mandatory values", func() {

			chv := validation.NewCHValidator()

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, apiSpecLocation, contextId)

			Convey("Then I am given an errors array response as validation errors have been found", func() {
				mandatoryErrorsResponseBody := readRequestBody(requiredErrorResponseBodyLocation)
				match := compareActualToExpected(validationErrs, mandatoryErrorsResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}

// TestOfficerDeltaSchemaEnumErrors asserts that when an invalid request body is given with incorrect ENUM values.
// then an errors array is returned, stating that given values are incorrect.
func TestOfficerDeltaSchemaEnumErrors(t *testing.T) {

	Convey("Given I want to test the officers-delta API schema to assert ENUM validation is working correctly", t, func() {

		enumErrorsRequestBody := readRequestBody(enumErrorRequestBodyLocation)

		r := httptest.NewRequest("POST", officersEndpoint, bytes.NewBuffer(enumErrorsRequestBody))
		r = setHeaders(r)

		Convey("When I call to validate the request body, providing an valid request with incorrect ENUM values", func() {

			chv := validation.NewCHValidator()

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, apiSpecLocation, contextId)

			Convey("Then I am given an errors array response as validation errors have been found", func() {
				enumErrorsResponseBody := readRequestBody(enumErrorResponseBodyLocation)
				match := compareActualToExpected(validationErrs, enumErrorsResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}

// TestOfficerDeltaSchemaNoRequestBodyError asserts that when a missing request body is given,
// then an error is returned, stating that request body is missing.
func TestOfficerDeltaSchemaNoRequestBodyError(t *testing.T) {

	Convey("Given I want to test the officers-delta API schema to assert validation is working correctly", t, func() {

		r := httptest.NewRequest("POST", officersEndpoint, bytes.NewBuffer(nil))
		r = setHeaders(r)

		Convey("When I call to validate the request body, providing an empty request body", func() {

			chv := validation.NewCHValidator()

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, apiSpecLocation, contextId)

			Convey("Then I am given an error saying no request body provided", func() {
				noRequestBodyErrorsResponseBody := readRequestBody(noRequestBodyErrorResponseBodyLocation)
				match := compareActualToExpected(validationErrs, noRequestBodyErrorsResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}

// setHeaders sets the required headers for the unit testing to correctly function.
func setHeaders(r *http.Request) *http.Request {
	r.Header.Set(xRequestId, contextId)
	r.Header.Set(contentType, applicationJson)

	return r
}

// readRequestBody takes a file location in the format of a string relative path and reads in the contents from the file,
// removing any extra tabs, spaces and special characters using json.Compact. Returns a []byte version of the formatted file read in.
func readRequestBody(fl string) []byte {

	// Read in contents and covert it to a string for further processing.
	raw, _ := os.ReadFile(fl)

	buffer := new(bytes.Buffer)
	_ = json.Compact(buffer, raw)

	raw = buffer.Bytes()
	// Convert back to an []byte and return.
	return raw
}

// compareActualToExpected takes actual and expected json (as byte arrays) and compares them to see if they match. Ordering of response
// isn't always guaranteed so using this function to match them without having to worry about the order changing.
func compareActualToExpected(actual, expected []byte) bool {

	// Define 2 model CHError arrays to hold actual and expected responses.
	var actualErrArr *[]models.CHError
	var expectedErrArr *[]models.CHError

	// Convert responses to JSON object arrays.
	_ = json.Unmarshal(actual, &actualErrArr)
	_ = json.Unmarshal(expected, &expectedErrArr)

	// If the arrays don't match in length, they can't be a complete match so return false.
	if len(*actualErrArr) != len(*expectedErrArr) {
		return false
	}

	// create a map of CHError.Location (string) -> int to compare if both arrays are completely equal.
	diff := make(map[string]int, len(*expectedErrArr))

	// Range over the expected response array and add them to the newly created map.
	for _, ee := range *expectedErrArr {
		// 0 value for int is 0, so just increment a counter for the string
		diff[ee.Location]++
	}

	// Finally range over the actual response array errors and check that they exist in the expected errors map.
	for _, ae := range *actualErrArr {
		// If the error is not in diff bail out early as they can't be a match.
		if _, ok := diff[ae.Location]; !ok {
			return false
		}
		diff[ae.Location] -= 1
		if diff[ae.Location] == 0 {
			delete(diff, ae.Location)
		}
	}

	// Return if length of remaining map values is equal to 0. If it is, we have a match.
	return len(diff) == 0
}
