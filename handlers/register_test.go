package handlers

import (
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
