package insolvency

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/companieshouse/chs-delta-api/validation"
	"github.com/companieshouse/chs-delta-api/validation/schema_testing/common"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	requestBodiesLocation            = "./request_bodies/"
	okRequestBodyLocation            = requestBodiesLocation + "ok_request_body"
	typeErrorRequestBodyLocation     = requestBodiesLocation + "type_error_request_body"
	requiredErrorRequestBodyLocation = requestBodiesLocation + "required_error_request_body"

	responseBodiesLocation                 = "./response_bodies/"
	typeErrorResponseBodyLocation          = responseBodiesLocation + "type_error_response_body"
	requiredErrorResponseBodyLocation      = responseBodiesLocation + "required_error_response_body"
	noRequestBodyErrorResponseBodyLocation = responseBodiesLocation + "no_request_body_error_response_body"

	insolvencyEndpoint = "/delta/insolvency"
	apiSpecLocation    = "../../../apispec/api-spec.yml"
	contextId          = "contextId"
	methodPost         = "POST"
)

// TestUnitInsolvencyDeltaSchemaNoErrors asserts that when a valid request body is given which matches the schema, then no
// errors are returned.
func TestUnitInsolvencyDeltaSchemaNoErrors(t *testing.T) {

	Convey("Given I want to test the insolvency-delta API schema", t, func() {

		okRequestBody := common.ReadRequestBody(okRequestBodyLocation)

		r := httptest.NewRequest(methodPost, insolvencyEndpoint, bytes.NewBuffer(okRequestBody))
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

// TestUnitInsolvencyDeltaSchemaTypeErrors asserts that when an invalid request body is given with type errors (int provided
// instead of string), then an errors array is returned.
func TestUnitInsolvencyDeltaSchemaTypeErrors(t *testing.T) {

	Convey("Given I want to test the insolvency-delta API schema for type assertions", t, func() {

		typeErrorRequestBody := common.ReadRequestBody(typeErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, insolvencyEndpoint, bytes.NewBuffer(typeErrorRequestBody))
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

// TestUnitInsolvencyDeltaSchemaRequiredErrors asserts that when an invalid request body is given with missing mandatory values.
// then an errors array is returned, stating that required values are missing.
func TestUnitInsolvencyDeltaSchemaRequiredErrors(t *testing.T) {

	Convey("Given I want to test the insolvency-delta API schema to assert mandatory validation is working correctly", t, func() {

		mandatoryErrorsRequestBody := common.ReadRequestBody(requiredErrorRequestBodyLocation)

		r := httptest.NewRequest(methodPost, insolvencyEndpoint, bytes.NewBuffer(mandatoryErrorsRequestBody))
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

// TestUnitInsolvencyDeltaSchemaNoRequestBodyError asserts that when a missing request body is given,
// then an error is returned, stating that request body is missing.
func TestUnitInsolvencyDeltaSchemaNoRequestBodyError(t *testing.T) {

	Convey("Given I want to test the insolvency-delta API schema to assert validation is working correctly", t, func() {

		r := httptest.NewRequest(methodPost, insolvencyEndpoint, bytes.NewBuffer(nil))
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
