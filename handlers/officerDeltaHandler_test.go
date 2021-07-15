package handlers

import (
	"bytes"
	"errors"
	hMocks "github.com/companieshouse/chs-delta-api/helpers/mocks"
	sMocks "github.com/companieshouse/chs-delta-api/services/mocks"
	"github.com/golang/mock/gomock"
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

	topic = "officers-delta"
)

// TestNewOfficerDeltaHandler asserts that the constructor for the OfficerDeltaHandler returns a fully configured handler.
func TestNewOfficerDeltaHandler(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("When I call the constructor, then it returns me a valid OfficerDeltaHandler", t, func() {

		svc := sMocks.NewMockKafkaService(mockCtrl)
		h := hMocks.NewMockHelper(mockCtrl)

		officerHandler := NewOfficerDeltaHandler(svc, h)

		So(officerHandler, ShouldNotBeNil)

		So(officerHandler.kSvc, ShouldNotBeNil)
		So(officerHandler.kSvc, ShouldEqual, svc)

		So(officerHandler.h, ShouldNotBeNil)
		So(officerHandler.h, ShouldEqual, h)
	})
}

// TestOfficerDeltaHandlerFailsRequestBodyRetrieval asserts that when converting the request body fails, errors are
// handled correctly and returned to the user with the correct status.
func TestOfficerDeltaHandlerFailsRequestBodyRetrieval(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the officer delta endpoint", t, func() {

		req := httptest.NewRequest("POST", "/delta/officers", bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the kafka service fails to send", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)

			h.EXPECT().GetDataFromRequest(req).Return("", errors.New("error converting request body"))

			handler := NewOfficerDeltaHandler(svc, h)
			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 500 and an error returned", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

// TestOfficerDeltaHandlerSuccessfullySends asserts that you can send a REST request onto a kafka topic with no errors.
func TestOfficerDeltaHandlerSuccessfullySends(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the officer delta endpoint", t, func() {

		req := httptest.NewRequest("POST", "/delta/officers", bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)

			h.EXPECT().GetDataFromRequest(req).Return(requestBody, nil)
			svc.EXPECT().SendMessage(topic, requestBody).Return(nil)

			handler := NewOfficerDeltaHandler(svc, h)
			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

// TestOfficerDeltaHandlerFailsSend asserts that the officerDeltaHandler returns a bad request status when sending fails.
func TestOfficerDeltaHandlerFailsSend(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the officer delta endpoint", t, func() {

		req := httptest.NewRequest("POST", "/delta/officers", bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the kafka service fails to send", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)

			h.EXPECT().GetDataFromRequest(req).Return(requestBody, nil)
			svc.EXPECT().SendMessage(topic, requestBody).Return(errors.New("error sending message"))

			handler := NewOfficerDeltaHandler(svc, h)
			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 500 and an error returned", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
