package handlers

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	requestBody = `{
        "officers" : [
			{
				"company_number" : "99999999",
				"company_name" : "Test company 1",
				"forename" : "Test forename 1",
				"middle_name" : "Test middle name 1",
				"surname" : "Test surname 1",
				"age" : 25
			},
			{
				"company_number" : "55555555",
				"company_name" : "Test company 2",
				"forename" : "Test forename 2",
				"middle_name" : "Test middle name 2",
				"surname" : "Test surname 2",
				"age" : 39
			}
    	],
    	"CreatedTime" : "07-JUN-21 15.26.17.000000",
    	"delta_at" : "20140925171003950844"
	}`
)

func TestNewOfficerDeltaHandler(t *testing.T) {

	Convey("When I call the constructor, then it returns me a valid OfficerDeltaHandler", t, func() {

		officerHandler := NewOfficerDeltaHandler()

		So(officerHandler, ShouldNotBeNil)
	})
}

func TestOfficerDeltaHandlerWithCorrectRoute(t *testing.T) {
	Convey("Given a HTTP POST request via the officer delta endpoint", t, func() {

		req := httptest.NewRequest("POST", "/delta/officers", bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router", func() {

			handler := NewOfficerDeltaHandler()
			handler.HandleOfficerDelta(resp, req)

			Convey("Then the response should be 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}
