package exemption

import (
	"bytes"
	"github.com/companieshouse/chs-delta-api/validation"
	"github.com/companieshouse/chs-delta-api/validation/schema_testing/common"
	"github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"testing"
)

const (
	requestBodiesLocation                           = "./request_bodies/"
	okRequestBodyLocation                           = requestBodiesLocation + "valid_request.json"
	missingTopLevelFieldsRequestBodyLocation        = requestBodiesLocation + "missing_top_level_fields_request.json"
	missingNestedExemptionFieldsRequestBodyLocation = requestBodiesLocation + "missing_description_items_fields_request.json"
	invalidDataTypeRequestBodyLocation              = requestBodiesLocation + "invalid_data_type_request.json"

	responseBodiesLocation                           = "./response_bodies/"
	missingTopLevelFieldsResponseBodyLocation        = responseBodiesLocation + "missing_top_level_fields_response.json"
	missingNestedExemptionFieldsResponseBodyLocation = responseBodiesLocation + "missing_description_items_fields_response.json"
	invalidDataTypeResponseBodyLocation              = responseBodiesLocation + "invalid_data_type_response.json"

	exemptionEndpoint = "/delta/exemption"
	apiSpecLocation   = "../../../apispec/api-spec.yml"
	contextId         = "contextId"
	methodPost        = "POST"
)

func TestUnitExemptionDeltaSchemaNoErrors(t *testing.T) {

	convey.Convey("Given a valid exemption delta request body has been specified", t, func() {

		requestBody := common.ReadRequestBody(okRequestBodyLocation)

		r := httptest.NewRequest(methodPost, exemptionEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		convey.Convey("When the request is validated", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			convey.Convey("Then the returned value should be nil", func() {
				convey.So(validationErrs, convey.ShouldBeNil)
			})
		})
	})
}

func TestUnitExemptionDeltaSchemaRaisesErrorsIfTopLevelPropertiesAbsent(t *testing.T) {
	convey.Convey("Given company_number and exemption fields are absent from an exemption delta request", t, func() {
		requestBody := common.ReadRequestBody(missingTopLevelFieldsRequestBodyLocation)

		r := httptest.NewRequest(methodPost, exemptionEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		convey.Convey("When the request is validated", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			convey.Convey("Then all validation errors should be returned", func() {
				responseBody := common.ReadRequestBody(missingTopLevelFieldsResponseBodyLocation)
				match := common.CompareActualToExpected(validationErrs, responseBody)

				convey.So(validationErrs, convey.ShouldNotBeNil)
				convey.So(match, convey.ShouldBeTrue)
			})
		})
	})
}

func TestUnitExemptionDeltaSchemaRaisesErrorsIfPropertiesAbsentFromExemptionTypes(t *testing.T) {
	convey.Convey("Given description and items fields are absent from individual exemptions in the request", t, func() {
		requestBody := common.ReadRequestBody(missingNestedExemptionFieldsRequestBodyLocation)

		r := httptest.NewRequest(methodPost, exemptionEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		convey.Convey("When the request is validated", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			convey.Convey("Then all validation errors should be returned", func() {
				responseBody := common.ReadRequestBody(missingNestedExemptionFieldsResponseBodyLocation)
				match := common.CompareActualToExpected(validationErrs, responseBody)

				convey.So(validationErrs, convey.ShouldNotBeNil)
				convey.So(match, convey.ShouldBeTrue)
			})
		})
	})
}

func TestUnitExemptionDeltaSchemaRaisesErrorIfInvalidDataTypesProvided(t *testing.T) {
	convey.Convey("Given the request contains invalid data types for expected fields", t, func() {
		requestBody := common.ReadRequestBody(invalidDataTypeRequestBodyLocation)

		r := httptest.NewRequest(methodPost, exemptionEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		convey.Convey("When the request is validated", func() {

			chv, _ := validation.NewCHValidator(apiSpecLocation)

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			convey.Convey("Then all validation errors should be returned", func() {
				responseBody := common.ReadRequestBody(invalidDataTypeResponseBodyLocation)
				match := common.CompareActualToExpected(validationErrs, responseBody)

				convey.So(validationErrs, convey.ShouldNotBeNil)
				convey.So(match, convey.ShouldBeTrue)
			})
		})
	})
}
