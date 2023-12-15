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
	requestBodiesLocation = "./request_bodies/"
	okRequestBodyLocation = requestBodiesLocation + "valid_request.json"

	filingHistoryEndpoint = "/delta/filing-history"
	apiSpecLocation       = "../../../apispec/api-spec.yml"
	contextId             = "contextId"
	methodPost            = "POST"
)

func TestUnitFilingHistoryDeltaSchemaNoErrors(t *testing.T) {

	convey.Convey("Given a valid filing history delta request body has been specified", t, func() {

		requestBody := common.ReadRequestBody(okRequestBodyLocation)

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
