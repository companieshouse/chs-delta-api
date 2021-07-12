package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitHealthCheck(t *testing.T) {
	Convey("When I call the healthcheck endpoint, then I am given a 200 status", t, func() {
		w := httptest.ResponseRecorder{}
		healthCheck(&w, nil)
		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

func TestRegister(t *testing.T) {
	Convey("When we call the register function then all routes are registered", t, func() {
		router := mux.NewRouter()
		Register(router)
		So(router.GetRoute("healthcheck"), ShouldNotBeNil)
		So(router.GetRoute("officer-delta"), ShouldNotBeNil)
	})
}