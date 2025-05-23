package company

import (
	"bytes"
	"github.com/companieshouse/chs-delta-api/validation"
	"github.com/companieshouse/chs-delta-api/validation/schema_testing/common"
	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"testing"
)

const (
	requestBodiesLocation             = "./request_bodies/"
	okRequestBodyLocation             = requestBodiesLocation + "valid_request.json"
	invalidEnumRequestLocation        = requestBodiesLocation + "invalid_enum_request.json"
	invalidValueFormatRequestLocation = requestBodiesLocation + "invalid_value_format_request.json"
	invalidDataTypeRequestLocation    = requestBodiesLocation + "invalid_data_type_request.json"
	invalidComplexDataTypeRequest     = requestBodiesLocation + "invalid_complex_data_type_request.json"
	deleteRequestBodyLocation         = requestBodiesLocation + "delete_ok_request_body.json"
	deleteBadRequestBodyLocation      = requestBodiesLocation + "delete_bad_request_body.json"

	responseBodiesLocation             = "./response_bodies/"
	invalidEnumResponseLocation        = responseBodiesLocation + "invalid_enum_error_response.json"
	invalidValueFormatResponseLocation = responseBodiesLocation + "invalid_value_format_error_response.json"
	invalidDataTypeResponseLocation    = responseBodiesLocation + "invalid_data_type_error_response.json"
	invalidComplexDataTypeResponse     = responseBodiesLocation + "invalid_complex_data_type_error_response.json"
	deleteErrorResponseBodyLocation    = responseBodiesLocation + "delete_error_response_body"

	companyEndpoint       = "/delta/company"
	companyDeleteEndpoint = "/delta/company/delete"
	apiSpecLocation       = "../../../ecs-image-build/apispec/api-spec.yml"
	contextId             = "contextId"
	methodPost            = "POST"
)

func TestUnitCompanyDeltaSchemaNoErrors(t *testing.T) {

	Convey("Given a valid company delta request body has been specified", t, func() {

		okRequestBody := common.ReadRequestBody(okRequestBodyLocation)

		r := httptest.NewRequest(methodPost, companyEndpoint, bytes.NewBuffer(okRequestBody))
		r = common.SetHeaders(r)

		Convey("When the request is validated", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then the returned value should be nil", func() {
				So(validationErrs, ShouldBeNil)
			})
		})
	})
}

func TestUnitCompanyDeltaSchemaReturnsErrorIfInvalidEnumValueSpecified(t *testing.T) {

	Convey("Given an invalid enum value in a company delta request body has been specified", t, func() {

		okRequestBody := common.ReadRequestBody(invalidEnumRequestLocation)

		r := httptest.NewRequest(methodPost, companyEndpoint, bytes.NewBuffer(okRequestBody))
		r = common.SetHeaders(r)

		Convey("When the request is validated", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then errors should be returned", func() {
				typeErrorResponseBody := common.ReadRequestBody(invalidEnumResponseLocation)
				match := common.CompareActualToExpected(validationErrs, typeErrorResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}

func TestUnitCompanyDeltaSchemaReturnsErrorIfInvalidValueFormatSpecified(t *testing.T) {

	Convey("Given values in a company delta request body do not match expected constraints", t, func() {

		okRequestBody := common.ReadRequestBody(invalidValueFormatRequestLocation)

		r := httptest.NewRequest(methodPost, companyEndpoint, bytes.NewBuffer(okRequestBody))
		r = common.SetHeaders(r)

		Convey("When the request is validated", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then errors should be returned", func() {
				typeErrorResponseBody := common.ReadRequestBody(invalidValueFormatResponseLocation)
				match := common.CompareActualToExpected(validationErrs, typeErrorResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}

func TestUnitCompanyDeltaSchemaReturnsErrorIfInvalidDataTypesSpecified(t *testing.T) {

	Convey("Given values in a company delta request body do not match expected data types", t, func() {

		okRequestBody := common.ReadRequestBody(invalidDataTypeRequestLocation)

		r := httptest.NewRequest(methodPost, companyEndpoint, bytes.NewBuffer(okRequestBody))
		r = common.SetHeaders(r)

		Convey("When the request is validated", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then errors should be returned", func() {
				typeErrorResponseBody := common.ReadRequestBody(invalidDataTypeResponseLocation)
				match := common.CompareActualToExpected(validationErrs, typeErrorResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}

func TestUnitCompanyDeltaSchemaReturnsErrorIfInvalidComplexDataTypesSpecified(t *testing.T) {

	Convey("Given values in a company delta request body do not match expected data types", t, func() {

		okRequestBody := common.ReadRequestBody(invalidComplexDataTypeRequest)

		r := httptest.NewRequest(methodPost, companyEndpoint, bytes.NewBuffer(okRequestBody))
		r = common.SetHeaders(r)

		Convey("When the request is validated", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then errors should be returned", func() {
				typeErrorResponseBody := common.ReadRequestBody(invalidComplexDataTypeResponse)
				match := common.CompareActualToExpected(validationErrs, typeErrorResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}

func TestUnitCompanyDeleteDeltaSchemaNoErrors(t *testing.T) {

	Convey("Given I want to test the company-delete-delta API schema", t, func() {

		deleteRequestBody := common.ReadRequestBody(deleteRequestBodyLocation)

		r := httptest.NewRequest(methodPost, companyDeleteEndpoint, bytes.NewBuffer(deleteRequestBody))
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

func TestUnitCompanyDeleteDeltaSchemaBadRequestError(t *testing.T) {

	Convey("Given I want to test the company-delete-delta API schema", t, func() {

		deleteRequestBody := common.ReadRequestBody(deleteBadRequestBodyLocation)

		r := httptest.NewRequest(methodPost, companyDeleteEndpoint, bytes.NewBuffer(deleteRequestBody))
		r = common.SetHeaders(r)

		Convey("When I call to validate the request body, providing a valid request", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			Convey("Then I am given an errors array response as validation errors have been found", func() {
				mandatoryErrorsResponseBody := common.ReadRequestBody(deleteErrorResponseBodyLocation)
				match := common.CompareActualToExpected(validationErrs, mandatoryErrorsResponseBody)

				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}
