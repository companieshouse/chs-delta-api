package charges

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
	deleteRequestBodyLocation          = requestBodiesLocation + "valid_delete_request"
	missingFieldsDeleteRequestBody     = requestBodiesLocation + "missing_required_fields_delete_request.json"
	exceedsMaxLengthDeleteRequestBody  = requestBodiesLocation + "fields_exceeds_max_length_delete_request.json"
	invalidDataTypeDeleteRequestBody   = requestBodiesLocation + "invalid_data_type_delete_request.json"

	responseBodiesLocation                 = "./response_bodies/"
	typeErrorResponseBodyLocation          = responseBodiesLocation + "type_error_response_body"
	requiredErrorResponseBodyLocation      = responseBodiesLocation + "required_error_response_body"
	noRequestBodyErrorResponseBodyLocation = responseBodiesLocation + "no_request_body_error_response_body"
	dateLengthErrorResponseBodyLocation    = responseBodiesLocation + "date_length_error_response_body"
	regexErrorResponseBodyLocation         = responseBodiesLocation + "regex_error_response_body"
	missingFieldsDeleteResponseBody        = responseBodiesLocation + "missing_required_fields_delete_response.json"
	exceedsMaxLengthDeleteResponseBody     = responseBodiesLocation + "fields_exceeds_max_length_delete_response.json"
	invalidDataTypeDeleteResponseBody      = responseBodiesLocation + "invalid_data_type_delete_response.json"

	chargesEndpoint       = "/delta/charges"
	chargesDeleteEndpoint = "/delta/charges/delete"
	apiSpecLocation       = "../../../apispec/api-spec.yml"
	contextId             = "contextId"
	methodPost            = "POST"
)

// TestUnitChargesDeltaSchemaNoErrors asserts that when a valid request body is given which matches the schema, then no
// errors are returned.
func TestUnitChargesDeltaSchemaNoErrors(t *testing.T) {

	Convey("Given I want to test the charges-delta API schema", t, func() {

		okRequestBody := common.ReadRequestBody(okRequestBodyLocation)

		r := httptest.NewRequest(methodPost, chargesEndpoint, bytes.NewBuffer(okRequestBody))
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

// TestUnitChargesDeltaSchemaTypeErrors asserts that when an invalid request body is given with type errors (int provided
// instead of string), then an errors array is returned.
func TestUnitChargesDeltaSchemaTypeErrors(t *testing.T) {

	Convey("Given I want to test the charges-delta API schema for type assertions", t, func() {

		typeErrorRequestBody := common.ReadRequestBody(typeErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, chargesEndpoint, bytes.NewBuffer(typeErrorRequestBody))
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

// TestUnitChargesDeltaSchemaRequiredErrors asserts that when an invalid request body is given with missing mandatory values,
// then an errors array is returned, stating that required values are missing.
func TestUnitChargesDeltaSchemaRequiredErrors(t *testing.T) {

	Convey("Given I want to test the charges-delta API schema to assert mandatory validation is working correctly", t, func() {

		mandatoryErrorsRequestBody := common.ReadRequestBody(requiredErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, chargesEndpoint, bytes.NewBuffer(mandatoryErrorsRequestBody))
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

// TestUnitChargesDeltaDateLengthErrors asserts that when a request body is given with dates which are not of length 8
// then an errors array is returned.
// NOTE: there is no validation on the format of the dates in the spec, only properties asserted are type: string and [min|max]Length: 8
func TestUnitChargesDeltaDateLengthErrors(t *testing.T) {

	Convey("Given I want to test the charges-delta API schema to assert mandatory validation is working correctly", t, func() {
		dateErrorsRequestBody := common.ReadRequestBody(dateLengthErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, chargesEndpoint, bytes.NewBuffer(dateErrorsRequestBody))
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

// TestUnitChargesDeltaSchemaNoRequestBodyError asserts that when a missing request body is given,
// then an error is returned, stating that request body is missing.
func TestUnitChargesDeltaSchemaNoRequestBodyError(t *testing.T) {

	Convey("Given I want to test the charges-delta API schema to assert validation is working correctly", t, func() {

		r := httptest.NewRequest(methodPost, chargesEndpoint, bytes.NewBuffer(nil))
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

// TestUnitChargesDeltaSchemaRegexErrors asserts that when an invalid request body is given with incorrect Regular Expression pattern values,
// then an errors array is returned, stating that required values are incorrect.
func TestUnitChargesDeltaSchemaRegexErrors(t *testing.T) {

	Convey("Given I want to test the charges-delta API schema to assert Regular Expression validation is working correctly", t, func() {

		mandatoryErrorsRequestBody := common.ReadRequestBody(regexErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, chargesEndpoint, bytes.NewBuffer(mandatoryErrorsRequestBody))
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

// TestUnitChargesDeleteDeltaSchemaNoErrors asserts that when a valid request body is given which matches the schema, then no
// errors are returned.
func TestUnitChargesDeleteDeltaSchemaNoErrors(t *testing.T) {

	Convey("Given I want to test the charges-delete-delta API schema", t, func() {

		deleteRequestBody := common.ReadRequestBody(deleteRequestBodyLocation)

		r := httptest.NewRequest(methodPost, chargesDeleteEndpoint, bytes.NewBuffer(deleteRequestBody))
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

func TestUnitFilingHistoryDeleteDeltaSchemasWhereFieldsAreMissing(t *testing.T) {
	Convey("Given a charge delete delta request body is missing top level fields", t, func() {
		requestBody := common.ReadRequestBody(missingFieldsDeleteRequestBody)
		r := httptest.NewRequest(methodPost, chargesDeleteEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		Convey("When the request is validated", func() {
			chv, _ := validation.NewCHValidator(apiSpecLocation)
			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then all validation errors should be returned", func() {
				responseBody := common.ReadRequestBody(missingFieldsDeleteResponseBody)
				match := common.CompareActualToExpected(validationErrs, responseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldBeTrue)
			})
		})
	})
}

func TestUnitFilingHistoryDeleteDeltaSchemasWhereFieldsExceedsMaxLength(t *testing.T) {
	Convey("Given a charge delete delta request body where entity id is over 10 characters", t, func() {
		requestBody := common.ReadRequestBody(exceedsMaxLengthDeleteRequestBody)
		r := httptest.NewRequest(methodPost, chargesDeleteEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		Convey("When the request is validated", func() {
			chv, _ := validation.NewCHValidator(apiSpecLocation)
			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then all validation errors should be returned", func() {
				responseBody := common.ReadRequestBody(exceedsMaxLengthDeleteResponseBody)
				match := common.CompareActualToExpected(validationErrs, responseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldBeTrue)
			})
		})
	})
}

func TestUnitFilingHistoryDeleteDeltaSchemasWithInvalidDataTypeForFields(t *testing.T) {
	Convey("Given a charge delete delta request body where entity id is not a string", t, func() {
		requestBody := common.ReadRequestBody(invalidDataTypeDeleteRequestBody)
		r := httptest.NewRequest(methodPost, chargesDeleteEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		Convey("When the request is validated", func() {
			chv, _ := validation.NewCHValidator(apiSpecLocation)
			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then all validation errors should be returned", func() {
				responseBody := common.ReadRequestBody(invalidDataTypeDeleteResponseBody)
				match := common.CompareActualToExpected(validationErrs, responseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldBeTrue)
			})
		})
	})
}
