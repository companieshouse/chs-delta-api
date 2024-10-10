package acspprofile

import (
	"bytes"
	"github.com/companieshouse/chs-delta-api/validation"
	"github.com/companieshouse/chs-delta-api/validation/schema_testing/common"
	"github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"testing"
)

const (
	requestBodiesLocation                  = "./request_bodies/"
	validRequestBody                       = requestBodiesLocation + "valid_request.json"
	missingFieldsRequestBody               = requestBodiesLocation + "missing_required_fields_request.json"
	invalidDataTypeRequestBody             = requestBodiesLocation + "invalid_data_type_request.json"
	fieldsFailLengthConstraintsRequestBody = requestBodiesLocation + "fields_fail_length_constraints_request.json"

	responseBodiesLocation                  = "./response_bodies/"
	missingFieldsResponseBody               = responseBodiesLocation + "missing_required_fields_response.json"
	invalidDataTypeResponseBody             = responseBodiesLocation + "invalid_data_type_response.json"
	fieldsFailLengthConstraintsResponseBody = responseBodiesLocation + "fields_fail_length_constraints_response.json"

	acspProfileEndpoint = "/delta/acsp"
	apiSpecLocation     = "../../../apispec/api-spec.yml"
	contextId           = "contextId"
	methodPost          = "POST"
)

func TestUnitAcspProfileDeltaSchemaNoErrors(t *testing.T) {
	convey.Convey("Given an valid acsp profile delta request body has been specified", t, func() {
		requestBody := common.ReadRequestBody(validRequestBody)

		r := httptest.NewRequest(methodPost, acspProfileEndpoint, bytes.NewBuffer(requestBody))
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

func TestUnitAcspProfileDeltaSchemaMissingRequiredFields(t *testing.T) {
	convey.Convey("Given an acsp profile delta request body is missing all required fields", t, func() {
		requestBody := common.ReadRequestBody(missingFieldsRequestBody)
		r := httptest.NewRequest(methodPost, acspProfileEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		convey.Convey("When the request is validated", func() {
			chv, _ := validation.NewCHValidator(apiSpecLocation)
			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			convey.Convey("Then all validation errors should be returned", func() {
				responseBody := common.ReadRequestBody(missingFieldsResponseBody)
				match := common.CompareActualToExpected(validationErrs, responseBody)

				convey.So(validationErrs, convey.ShouldNotBeNil)
				convey.So(match, convey.ShouldBeTrue)
			})
		})
	})
}

func TestUnitAcspProfileDeltaSchemaWithInvalidDataType(t *testing.T) {
	convey.Convey("Given an acsp profile delta request body contains invalid fields", t, func() {
		requestBody := common.ReadRequestBody(invalidDataTypeRequestBody)
		r := httptest.NewRequest(methodPost, acspProfileEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		convey.Convey("When the request is validated", func() {
			chv, _ := validation.NewCHValidator(apiSpecLocation)
			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			convey.Convey("Then all validation errors should be returned", func() {
				responseBody := common.ReadRequestBody(invalidDataTypeResponseBody)
				match := common.CompareActualToExpected(validationErrs, responseBody)

				convey.So(validationErrs, convey.ShouldNotBeNil)
				convey.So(match, convey.ShouldBeTrue)
			})
		})
	})
}

func TestUnitAcspProfileDeltaSchemaWhereFieldsFailLengthConstraints(t *testing.T) {
	convey.Convey("Given an acsp profile delta request body contains fields which fail length constraints", t, func() {
		requestBody := common.ReadRequestBody(fieldsFailLengthConstraintsRequestBody)
		r := httptest.NewRequest(methodPost, acspProfileEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		convey.Convey("When the request is validated", func() {
			chv, _ := validation.NewCHValidator(apiSpecLocation)
			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			convey.Convey("Then all validation errors should be returned", func() {
				responseBody := common.ReadRequestBody(fieldsFailLengthConstraintsResponseBody)
				match := common.CompareActualToExpected(validationErrs, responseBody)

				convey.So(validationErrs, convey.ShouldNotBeNil)
				convey.So(match, convey.ShouldBeTrue)
			})
		})
	})
}
