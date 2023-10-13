package officers

import (
	"bytes"
	"github.com/companieshouse/chs-delta-api/validation"
	"github.com/companieshouse/chs-delta-api/validation/schema_testing/common"
	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"testing"
)

const (
	requestBodiesLocation            = "./request_bodies/"
	okRequestBodyLocation            = requestBodiesLocation + "ok_request_body"
	typeErrorRequestBodyLocation     = requestBodiesLocation + "type_error_request_body"
	requiredErrorRequestBodyLocation = requestBodiesLocation + "required_error_request_body"
	enumErrorRequestBodyLocation     = requestBodiesLocation + "enum_error_request_body"
	maxPropertiesRequestBodyLocation = requestBodiesLocation + "max_properties_request_body"

	responseBodiesLocation                 = "./response_bodies/"
	typeErrorResponseBodyLocation          = responseBodiesLocation + "type_error_response_body"
	requiredErrorResponseBodyLocation      = responseBodiesLocation + "required_error_response_body"
	enumErrorResponseBodyLocation          = responseBodiesLocation + "enum_error_response_body"
	noRequestBodyErrorResponseBodyLocation = responseBodiesLocation + "no_request_body_error_response_body"
	maxPropertiesResponseBodyLocation      = responseBodiesLocation + "max_properties_response_body"
	deleteRequestBodyLocation              = requestBodiesLocation + "delete_request_body"

	officersEndpoint       = "/delta/officers"
	officersDeleteEndpoint = "/delta/officers/delete"
	apiSpecLocation        = "../../../apispec/api-spec.yml"
	contextId              = "contextId"
	methodPost             = "POST"
)

// TestUnitOfficerDeltaSchemaNoErrors asserts that when a valid request body is given which matches the schema, then no
// errors are returned.
func TestUnitOfficerDeltaSchemaNoErrors(t *testing.T) {

	Convey("Given I want to test the officers-delta API schema", t, func() {

		okRequestBody := common.ReadRequestBody(okRequestBodyLocation)

		r := httptest.NewRequest(methodPost, officersEndpoint, bytes.NewBuffer(okRequestBody))
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

// TestUnitOfficerDeltaSchemaTypeErrors asserts that when an invalid request body is given with type errors (int provided
// instead of string), then an errors array is returned.
func TestUnitOfficerDeltaSchemaTypeErrors(t *testing.T) {

	Convey("Given I want to test the officers-delta API schema for type assertions", t, func() {

		typeErrorRequestBody := common.ReadRequestBody(typeErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, officersEndpoint, bytes.NewBuffer(typeErrorRequestBody))
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

// TestUnitOfficerDeltaSchemaRequiredErrors asserts that when an invalid request body is given with missing mandatory values.
// then an errors array is returned, stating that required values are missing.
func TestUnitOfficerDeltaSchemaRequiredErrors(t *testing.T) {

	Convey("Given I want to test the officers-delta API schema to assert mandatory validation is working correctly", t, func() {

		mandatoryErrorsRequestBody := common.ReadRequestBody(requiredErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, officersEndpoint, bytes.NewBuffer(mandatoryErrorsRequestBody))
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

// TestUnitOfficerDeltaSchemaEnumErrors asserts that when an invalid request body is given with incorrect ENUM values.
// then an errors array is returned, stating that given values are incorrect.
func TestUnitOfficerDeltaSchemaEnumErrors(t *testing.T) {

	Convey("Given I want to test the officers-delta API schema to assert ENUM validation is working correctly", t, func() {

		enumErrorsRequestBody := common.ReadRequestBody(enumErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, officersEndpoint, bytes.NewBuffer(enumErrorsRequestBody))
		r = common.SetHeaders(r)

		Convey("When I call to validate the request body, providing an valid request with incorrect ENUM values", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then I am given an errors array response as validation errors have been found", func() {
				enumErrorsResponseBody := common.ReadRequestBody(enumErrorResponseBodyLocation)
				match := common.CompareActualToExpected(validationErrs, enumErrorsResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}

// TestUnitOfficerDeltaSchemaNoRequestBodyError asserts that when a missing request body is given,
// then an error is returned, stating that request body is missing.
func TestUnitOfficerDeltaSchemaNoRequestBodyError(t *testing.T) {

	Convey("Given I want to test the officers-delta API schema to assert validation is working correctly", t, func() {

		r := httptest.NewRequest(methodPost, officersEndpoint, bytes.NewBuffer(nil))
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

// TestOfficerDeltaSchemaMaxPropertiesError asserts that when an invalid request body is given which breaks the allowed
// bounds on the maxProperties field for the Identification object, then an errors array is returned.
func TestUnitOfficerDeltaSchemaMaxPropertiesError(t *testing.T) {

	Convey("Given I want to test the officers-delta API schema for maxProperty constraints", t, func() {

		maxPropertiesErrorRequestBody := common.ReadRequestBody(maxPropertiesRequestBodyLocation)

		r := httptest.NewRequest(methodPost, officersEndpoint, bytes.NewBuffer(maxPropertiesErrorRequestBody))
		r = common.SetHeaders(r)

		Convey("When I call to validate the request body, providing an valid request with maxProperty constraint errors", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then I am given an errors array response as validation errors have been found", func() {
				maxPropertiesErrorResponseBody := common.ReadRequestBody(maxPropertiesResponseBodyLocation)
				match := common.CompareActualToExpected(validationErrs, maxPropertiesErrorResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}

// TestUnitOfficersDeleteDeltaSchemaNoErrors asserts that when a valid request body is given which matches the schema, then no
// errors are returned.
func TestUnitOfficersDeleteDeltaSchemaNoErrors(t *testing.T) {

	Convey("Given I want to test the officers-delete-delta API schema", t, func() {

		deleteRequestBody := common.ReadRequestBody(deleteRequestBodyLocation)

		r := httptest.NewRequest(methodPost, officersDeleteEndpoint, bytes.NewBuffer(deleteRequestBody))
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
