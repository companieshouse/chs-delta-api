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
    	"created_at" : "01/01/2000",
    	"delta_at" : "01/01/2000"
	}`

	badRequestBody = `{
        "officers" : [
			{
				"company_number" : "99999999",
				"company_name" : "Test company 1",
				"forename" : "Test forename 1",
				"middle_name" : "Test middle name 1",
				"surname" : "Test surname 1",
				"age" : "25"
			}
    	],
	}`
)

func TestNewOfficersHandler(t *testing.T) {

	Convey("Then the constructor returns me a valid OfficerHandler", t, func() {

		officerHandler := NewOfficerHandler()

		So(officerHandler, ShouldNotBeNil)
	})
}

func TestOfficersHandlerWithCorrectRoute(t *testing.T) {
	Convey("Given a HTTP request for /delta/officers", t, func() {

		req := httptest.NewRequest("POST", "/delta/officers", bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router", func() {

			handler := NewOfficerHandler()
			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func TestOfficersHandlerWithBadRequest(t *testing.T) {
	Convey("Given a HTTP request for /delta/officers", t, func() {

		req := httptest.NewRequest("POST", "/delta/officers", bytes.NewBuffer([]byte(badRequestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router", func() {

			handler := NewOfficerHandler()
			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 400", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})
}