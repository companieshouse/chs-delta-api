package filinghistory

import (
	"bytes"
	"github.com/companieshouse/chs-delta-api/validation"
	"github.com/companieshouse/chs-delta-api/validation/schema_testing/common"
	"github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"testing"
)

const (
	requestBodiesLocation                = "./request_bodies/"
	validRequestBody                     = requestBodiesLocation + "valid_request.json"
	missingFieldsRequestBody             = requestBodiesLocation + "missing_required_fields_request.json"
	missingFilingHistoryArrayRequestBody = requestBodiesLocation + "missing_filing_history_array_request.json"
	invalidDataTypeRequestBody           = requestBodiesLocation + "invalid_data_type_request.json"
	exceedsMaxLengthRequestBody          = requestBodiesLocation + "fields_exceeds_max_length_request.json"
	validDeleteRequestBody               = requestBodiesLocation + "valid_request_delete.json"
	missingFieldsDeleteRequestBody       = requestBodiesLocation + "missing_required_fields_delete_request.json"
	exceedsMaxLengthDeleteRequestBody    = requestBodiesLocation + "fields_exceeds_max_length_delete_request.json"
	invalidDataTypeDeleteRequestBody     = requestBodiesLocation + "invalid_data_type_delete_request.json"

	responseBodiesLocation                = "./response_bodies/"
	missingFieldsResponseBody             = responseBodiesLocation + "missing_required_fields_response.json"
	missingFieldsDeleteResponseBody       = responseBodiesLocation + "missing_required_fields_delete_response.json"
	missingFilingHistoryArrayResponseBody = responseBodiesLocation + "missing_filing_history_array_response.json"
	invalidDataTypeResponseBody           = responseBodiesLocation + "invalid_data_type_response.json"
	exceedsMaxLengthResponseBody          = responseBodiesLocation + "fields_exceeds_max_length_response.json"
	exceedsMaxLengthDeleteResponseBody    = responseBodiesLocation + "fields_exceeds_max_length_delete_response.json"
	invalidDataTypeDeleteResponseBody     = responseBodiesLocation + "invalid_data_type_delete_response.json"

	filingHistoryEndpoint       = "/delta/filing-history"
	filingHistoryDeleteEndpoint = "/delta/filing-history/delete"
	apiSpecLocation             = "../../../ecs-image-build/apispec/api-spec.yml"
	contextId                   = "contextId"
	methodPost                  = "POST"
)

func TestUnitFilingHistoryDeltaSchemaNoErrors(t *testing.T) {
	convey.Convey("Given a valid filing history delta request body has been specified", t, func() {
		requestBody := common.ReadRequestBody(validRequestBody)

		r := httptest.NewRequest(methodPost, filingHistoryEndpoint, bytes.NewBuffer(requestBody))
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

func TestUnitFilingHistoryDeltaSchemaMissingRequiredFields(t *testing.T) {
	convey.Convey("Given a filing history delta request body is missing all required fields", t, func() {
		requestBody := common.ReadRequestBody(missingFieldsRequestBody)
		r := httptest.NewRequest(methodPost, filingHistoryEndpoint, bytes.NewBuffer(requestBody))
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

func TestUnitFilingHistoryDeltaSchemaMissingTopLevelField(t *testing.T) {
	convey.Convey("Given a filing history delta request body is missing the top level filing history array", t, func() {
		requestBody := common.ReadRequestBody(missingFilingHistoryArrayRequestBody)
		r := httptest.NewRequest(methodPost, filingHistoryEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		convey.Convey("When the request is validated", func() {
			chv, _ := validation.NewCHValidator(apiSpecLocation)
			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			convey.Convey("Then all validation errors should be returned", func() {
				responseBody := common.ReadRequestBody(missingFilingHistoryArrayResponseBody)
				match := common.CompareActualToExpected(validationErrs, responseBody)

				convey.So(validationErrs, convey.ShouldNotBeNil)
				convey.So(match, convey.ShouldBeTrue)
			})
		})
	})
}

func TestUnitFilingHistoryDeltaSchemaWithInvalidDataType(t *testing.T) {
	convey.Convey("Given a filing history delta request body contains fields which are int instead of string", t, func() {
		requestBody := common.ReadRequestBody(invalidDataTypeRequestBody)
		r := httptest.NewRequest(methodPost, filingHistoryEndpoint, bytes.NewBuffer(requestBody))
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

func TestUnitFilingHistoryDeltaSchemaWhereFieldsExceedMaxLength(t *testing.T) {
	convey.Convey("Given a filing history delta request body contains fields which exceeds its maximum length", t, func() {
		requestBody := common.ReadRequestBody(exceedsMaxLengthRequestBody)
		r := httptest.NewRequest(methodPost, filingHistoryEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		convey.Convey("When the request is validated", func() {
			chv, _ := validation.NewCHValidator(apiSpecLocation)
			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			convey.Convey("Then all validation errors should be returned", func() {
				responseBody := common.ReadRequestBody(exceedsMaxLengthResponseBody)
				match := common.CompareActualToExpected(validationErrs, responseBody)

				convey.So(validationErrs, convey.ShouldNotBeNil)
				convey.So(match, convey.ShouldBeTrue)
			})
		})
	})
}

func TestUnitFilingHistoryDeleteDeltaSchemaNoErrors(t *testing.T) {
	convey.Convey("Given a valid filing history delete delta request body has been specified", t, func() {
		requestBody := common.ReadRequestBody(validDeleteRequestBody)

		r := httptest.NewRequest(methodPost, filingHistoryDeleteEndpoint, bytes.NewBuffer(requestBody))
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

func TestUnitFilingHistoryDeleteDeltaSchemasWhereFieldsAreMissing(t *testing.T) {
	convey.Convey("Given a filing history delete delta request body is missing top level fields", t, func() {
		requestBody := common.ReadRequestBody(missingFieldsDeleteRequestBody)
		r := httptest.NewRequest(methodPost, filingHistoryDeleteEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		convey.Convey("When the request is validated", func() {
			chv, _ := validation.NewCHValidator(apiSpecLocation)
			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			convey.Convey("Then all validation errors should be returned", func() {
				responseBody := common.ReadRequestBody(missingFieldsDeleteResponseBody)
				match := common.CompareActualToExpected(validationErrs, responseBody)

				convey.So(validationErrs, convey.ShouldNotBeNil)
				convey.So(match, convey.ShouldBeTrue)
			})
		})
	})
}

func TestUnitFilingHistoryDeleteDeltaSchemasWhereFieldsExceedsMaxLength(t *testing.T) {
	convey.Convey("Given a filing history delete delta request body where entity id is over 10 characters", t, func() {
		requestBody := common.ReadRequestBody(exceedsMaxLengthDeleteRequestBody)
		r := httptest.NewRequest(methodPost, filingHistoryDeleteEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		convey.Convey("When the request is validated", func() {
			chv, _ := validation.NewCHValidator(apiSpecLocation)
			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			convey.Convey("Then all validation errors should be returned", func() {
				responseBody := common.ReadRequestBody(exceedsMaxLengthDeleteResponseBody)
				match := common.CompareActualToExpected(validationErrs, responseBody)

				convey.So(validationErrs, convey.ShouldNotBeNil)
				convey.So(match, convey.ShouldBeTrue)
			})
		})
	})
}

func TestUnitFilingHistoryDeleteDeltaSchemasWithInvalidDataTypeForFields(t *testing.T) {
	convey.Convey("Given a filing history delete delta request body where entity id is not a string", t, func() {
		requestBody := common.ReadRequestBody(invalidDataTypeDeleteRequestBody)
		r := httptest.NewRequest(methodPost, filingHistoryDeleteEndpoint, bytes.NewBuffer(requestBody))
		r = common.SetHeaders(r)

		convey.Convey("When the request is validated", func() {
			chv, _ := validation.NewCHValidator(apiSpecLocation)
			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, contextId)

			convey.Convey("Then all validation errors should be returned", func() {
				responseBody := common.ReadRequestBody(invalidDataTypeDeleteResponseBody)
				match := common.CompareActualToExpected(validationErrs, responseBody)

				convey.So(validationErrs, convey.ShouldNotBeNil)
				convey.So(match, convey.ShouldBeTrue)
			})
		})
	})
}
