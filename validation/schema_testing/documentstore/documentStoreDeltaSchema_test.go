package documentstore

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
	typeErrorRequestBodyLocation       = requestBodiesLocation + "type_error_request_body"
	requiredErrorRequestBodyLocation   = requestBodiesLocation + "required_error_request_body"
	dateLengthErrorRequestBodyLocation = requestBodiesLocation + "date_length_error_request_body"
	regexErrorRequestBodyLocation      = requestBodiesLocation + "regex_error_request_body"

	responseBodiesLocation                 = "./response_bodies/"
	typeErrorResponseBodyLocation          = responseBodiesLocation + "type_error_response_body"
	requiredErrorResponseBodyLocation      = responseBodiesLocation + "required_error_response_body"
	noRequestBodyErrorResponseBodyLocation = responseBodiesLocation + "no_request_body_error_response_body"
	dateLengthErrorResponseBodyLocation    = responseBodiesLocation + "date_length_error_response_body"
	regexErrorResponseBodyLocation         = responseBodiesLocation + "regex_error_response_body"

	docStoreEndpoint = "/delta/document-store"
	apiSpecLocation  = "../../../ecs-image-build/apispec/api-spec.yml"
	contextId        = "contextId"
	methodPost       = "POST"
)

// TestUnitDocStoreDeltaSchemaNoErrors asserts that when a valid request body is given which matches the schema, then no
// errors are returned.
func TestUnitDocStoreDeltaSchemaNoErrors(t *testing.T) {

	Convey("Given I want to test the doc-store-delta API schema", t, func() {

		okRequestBody := common.ReadRequestBody(okRequestBodyLocation)

		r := httptest.NewRequest(methodPost, docStoreEndpoint, bytes.NewBuffer(okRequestBody))
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

// TestUnitDocStoreDeltaSchemaTypeErrors asserts that when an invalid request body is given with type errors (int provided
// instead of string), then an errors array is returned.
func TestUnitDocStoreDeltaSchemaTypeErrors(t *testing.T) {

	Convey("Given I want to test the doc-store-delta API schema for type assertions", t, func() {

		typeErrorRequestBody := common.ReadRequestBody(typeErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, docStoreEndpoint, bytes.NewBuffer(typeErrorRequestBody))
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

// TestUnitDocStoreDeltaSchemaRequiredErrors asserts that when an invalid request body is given with missing mandatory values,
// then an errors array is returned, stating that required values are missing.
func TestUnitDocStoreDeltaSchemaRequiredErrors(t *testing.T) {

	Convey("Given I want to test the doc-store-delta API schema to assert mandatory validation is working correctly", t, func() {

		mandatoryErrorsRequestBody := common.ReadRequestBody(requiredErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, docStoreEndpoint, bytes.NewBuffer(mandatoryErrorsRequestBody))
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

// TestUnitDocStoreDeltaDateLengthErrors asserts that when a request body is given with dates which are not of length 8
// then an errors array is returned.
// NOTE: there is no validation on the format of the dates in the spec, only properties asserted are type: string and [min|max]Length: 8
func TestUnitDocStoreDeltaDateLengthErrors(t *testing.T) {

	Convey("Given I want to test the doc-store-delta API schema to assert mandatory validation is working correctly", t, func() {
		dateErrorsRequestBody := common.ReadRequestBody(dateLengthErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, docStoreEndpoint, bytes.NewBuffer(dateErrorsRequestBody))
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

// TestUnitDocStoreDeltaSchemaNoRequestBodyError asserts that when a missing request body is given,
// then an error is returned, stating that request body is missing.
func TestUnitDocStoreDeltaSchemaNoRequestBodyError(t *testing.T) {

	Convey("Given I want to test the doc-store-delta API schema to assert validation is working correctly", t, func() {

		r := httptest.NewRequest(methodPost, docStoreEndpoint, bytes.NewBuffer(nil))
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

// TestUnitDocStoreDeltaSchemaRegexErrors asserts that when an invalid request body is given with incorrect Regular Expression pattern values,
// then an errors array is returned, stating that required values are incorrect.
func TestUnitDocStoreDeltaSchemaRegexErrors(t *testing.T) {

	Convey("Given I want to test the doc-store-delta API schema to assert Regular Expression validation is working correctly", t, func() {

		mandatoryErrorsRequestBody := common.ReadRequestBody(regexErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, docStoreEndpoint, bytes.NewBuffer(mandatoryErrorsRequestBody))
		r = common.SetHeaders(r)

		Convey("When I call to validate the request body, providing an invalid request with not maching Regular Expression values", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then I am given an errors array response as validation errors have been found", func() {
				mandatoryErrorsResponseBody := common.ReadRequestBody(regexErrorResponseBodyLocation)
				match := common.CompareActualToExpected(validationErrs, mandatoryErrorsResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}
