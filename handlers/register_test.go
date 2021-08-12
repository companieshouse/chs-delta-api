package handlers

import (
	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs-delta-api/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"

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
		cfg, _ := config.Get()
		kSvc := mocks.NewMockKafkaService(mockCtrl)

		kSvc.EXPECT().Init(cfg).Return(nil)

		err := Register(router, cfg, kSvc)
		So(router.GetRoute("healthcheck"), ShouldNotBeNil)
		So(router.GetRoute("officer-delta"), ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
