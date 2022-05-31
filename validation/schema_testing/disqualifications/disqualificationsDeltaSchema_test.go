package disqualifications

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/companieshouse/chs-delta-api/validation"
	"github.com/companieshouse/chs-delta-api/validation/schema_testing/common"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	requestBodiesLocation              = "./request_bodies/"
	okRequestBodyLocation              = requestBodiesLocation + "ok_request_body"
	deleteRequestBodyLocation          = requestBodiesLocation + "delete_request_body"
	typeErrorRequestBodyLocation       = requestBodiesLocation + "type_error_request_body"
	requiredErrorRequestBodyLocation   = requestBodiesLocation + "required_error_request_body"
	dateLengthErrorRequestBodyLocation = requestBodiesLocation + "date_length_error_request_body"

	responseBodiesLocation                 = "./response_bodies/"
	typeErrorResponseBodyLocation          = responseBodiesLocation + "type_error_response_body"
	requiredErrorResponseBodyLocation      = responseBodiesLocation + "required_error_response_body"
	noRequestBodyErrorResponseBodyLocation = responseBodiesLocation + "no_request_body_error_response_body"
	dateLengthErrorResponseBodyLocation    = responseBodiesLocation + "date_length_error_response_body"

	disqualificationEndpoint       = "/delta/disqualification"
	disqualificationDeleteEndpoint = "/delta/disqualification/delete"
	apiSpecLocation                = "../../../apispec/api-spec.yml"
	contextId                      = "contextId"
	methodPost                     = "POST"
)

// TestUnitDisqualificationDeltaSchemaNoErrors asserts that when a valid request body is given which matches the schema, then no
// errors are returned.
func TestUnitDisqualificationDeltaSchemaNoErrors(t *testing.T) {

	Convey("Given I want to test the disqualification-delta API schema", t, func() {

		okRequestBody := common.ReadRequestBody(okRequestBodyLocation)

		r := httptest.NewRequest(methodPost, disqualificationEndpoint, bytes.NewBuffer(okRequestBody))
		r = common.SetHeaders(r)

		Convey("When I call to validate the request body, providing a valid request", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then I am given a nil response as no validation errors are returned", func() {
				So(validationErrs, ShouldBeNil)
			})
		})
	})
}

// TestUnitDisqualificationDeleteDeltaSchemaNoErrors asserts that when a valid request body is given which matches the schema, then no
// errors are returned.
func TestUnitDisqualificationDeleteDeltaSchemaNoErrors(t *testing.T) {

	Convey("Given I want to test the disqualification-delete-delta API schema", t, func() {

		deleteRequestBody := common.ReadRequestBody(deleteRequestBodyLocation)

		r := httptest.NewRequest(methodPost, disqualificationDeleteEndpoint, bytes.NewBuffer(deleteRequestBody))
		r = common.SetHeaders(r)

		Convey("When I call to validate the request body, providing a valid request", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then I am given a nil response as no validation errors are returned", func() {
				So(validationErrs, ShouldBeNil)
			})
		})
	})
}

// TestUnitDisqualificationDeltaSchemaTypeErrors asserts that when an invalid request body is given with type errors (int provided
// instead of string), then an errors array is returned.
func TestUnitDisqualificationDeltaSchemaTypeErrors(t *testing.T) {

	Convey("Given I want to test the disqualification-delta API schema for type assertions", t, func() {

		typeErrorRequestBody := common.ReadRequestBody(typeErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, disqualificationEndpoint, bytes.NewBuffer(typeErrorRequestBody))
		r = common.SetHeaders(r)

		Convey("When I call to validate the request body, providing an valid request with type errors", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then I am given an errors array response as validation errors have been found", func() {
				typeErrorResponseBody := common.ReadRequestBody(typeErrorResponseBodyLocation)
				match := common.CompareActualToExpected(validationErrs, typeErrorResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}

// TestUnitDisqualificationDeltaSchemaRequiredErrors asserts that when an invalid request body is given with missing mandatory values,
// then an errors array is returned, stating that required values are missing.
func TestUnitDisqualificationDeltaSchemaRequiredErrors(t *testing.T) {

	Convey("Given I want to test the disqualification-delta API schema to assert mandatory validation is working correctly", t, func() {

		mandatoryErrorsRequestBody := common.ReadRequestBody(requiredErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, disqualificationEndpoint, bytes.NewBuffer(mandatoryErrorsRequestBody))
		r = common.SetHeaders(r)

		Convey("When I call to validate the request body, providing an valid request with missing mandatory values", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then I am given an errors array response as validation errors have been found", func() {
				mandatoryErrorsResponseBody := common.ReadRequestBody(requiredErrorResponseBodyLocation)
				match := common.CompareActualToExpected(validationErrs, mandatoryErrorsResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}

// TestUnitDisqualificationDeltaDateLengthErrors asserts that when a request body is given with dates which are not of length 8
// then an errors array is returned.
// NOTE: there is no validation on the format of the dates in the spec, only properties asserted are type: string and [min|max]Length: 8
func TestUnitDisqualificationDeltaDateLengthErrors(t *testing.T) {

	Convey("Given I want to test the disqualification-delta API schema to assert mandatory validation is working correctly", t, func() {
		dateErrorsRequestBody := common.ReadRequestBody(dateLengthErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, disqualificationEndpoint, bytes.NewBuffer(dateErrorsRequestBody))
		r = common.SetHeaders(r)

		Convey("When I call to validate the request body, providing an valid request with missing mandatory values", func() {
			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then I am given an errors array response as validation errors have been found", func() {
				dateErrorsResponseBody := common.ReadRequestBody(dateLengthErrorResponseBodyLocation)
				match := common.CompareActualToExpected(validationErrs, dateErrorsResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}

// TestUnitDisqualificationDeltaSchemaNoRequestBodyError asserts that when a missing request body is given,
// then an error is returned, stating that request body is missing.
func TestUnitDisqualificationDeltaSchemaNoRequestBodyError(t *testing.T) {

	Convey("Given I want to test the disqualification-delta API schema to assert validation is working correctly", t, func() {

		r := httptest.NewRequest(methodPost, disqualificationEndpoint, bytes.NewBuffer(nil))
		r = common.SetHeaders(r)

		Convey("When I call to validate the request body, providing an empty request body", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then I am given an error saying no request body provided", func() {
				noRequestBodyErrorsResponseBody := common.ReadRequestBody(noRequestBodyErrorResponseBodyLocation)
				match := common.CompareActualToExpected(validationErrs, noRequestBodyErrorsResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}
