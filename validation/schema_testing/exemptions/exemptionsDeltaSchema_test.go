package exemptions

import (
	"bytes"
	"github.com/companieshouse/chs-delta-api/validation"
	"github.com/companieshouse/chs-delta-api/validation/schema_testing/common"
	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"testing"
)

const (
	requestBodiesLocation = "./request_bodies/"
	okRequestBodyLocation = requestBodiesLocation + "valid_request.json"

	responseBodiesLocation = "./response_bodies/"

	companyExemptionsEndpoint = "/delta/company-exemptions"
	apiSpecLocation           = "../../../apispec/api-spec.yml"
	contextId                 = "contextId"
	methodPost                = "POST"
)

func TestUnitCompanyExemptionsDeltaSchemaNoErrors(t *testing.T) {

	Convey("Given a valid company exemptions delta request body has been specified", t, func() {

		okRequestBody := common.ReadRequestBody(okRequestBodyLocation)

		r := httptest.NewRequest(methodPost, companyExemptionsEndpoint, bytes.NewBuffer(okRequestBody))
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
