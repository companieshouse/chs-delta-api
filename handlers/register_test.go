package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs-delta-api/services/mocks"
	"github.com/companieshouse/chs-delta-api/validation"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	. "github.com/smartystreets/goconvey/convey"
)

// TestUnitHealthCheck asserts that the healthcheck endpoint correctly returns 200 when called.
func TestUnitHealthCheck(t *testing.T) {
	Convey("When I call the healthcheck endpoint, then I am given a 200 status", t, func() {
		w := httptest.ResponseRecorder{}
		healthCheck(&w, nil)
		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

// TestUnitRegister asserts that all routes are correctly registered and can be called.
func TestUnitRegister(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("When we call the register function then all routes are registered", t, func() {
		router := mux.NewRouter()

		callNewCHValidator = func(openApiSpec string) (validation.CHValidator, error) {
			return &validation.CHValidatorImpl{}, nil
		}

		config.CallValidateConfig = func(cfg *config.Config) error {
			return nil
		}
		cfg, _ := config.Get()
		kSvc := mocks.NewMockKafkaService(mockCtrl)

		kSvc.EXPECT().Init(cfg).Return(nil)

		err := Register(router, cfg, kSvc)
		So(router.GetRoute("healthcheck"), ShouldNotBeNil)
		So(router.GetRoute("officer-delta"), ShouldNotBeNil)
		So(router.GetRoute("officer-delta-validate"), ShouldNotBeNil)
		So(router.GetRoute("insolvency-delta"), ShouldNotBeNil)
		So(router.GetRoute("insolvency-delta-validate"), ShouldNotBeNil)
		So(router.GetRoute("charges-delta"), ShouldNotBeNil)
		So(router.GetRoute("charges-delta-validate"), ShouldNotBeNil)
		So(router.GetRoute("disqualified-officer-delta"), ShouldNotBeNil)
		So(router.GetRoute("disqualified-officer-delta-validate"), ShouldNotBeNil)
		So(err, ShouldBeNil)
	})

}
