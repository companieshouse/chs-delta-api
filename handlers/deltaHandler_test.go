package handlers

import (
	"bytes"
	"errors"
	"github.com/companieshouse/chs-delta-api/config"
	hMocks "github.com/companieshouse/chs-delta-api/helpers/mocks"
	sMocks "github.com/companieshouse/chs-delta-api/services/mocks"
	chvMocks "github.com/companieshouse/chs-delta-api/validation/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	requestBody = `{"dummy" : "request"}`
	contextId   = "contextId"
	postMethod  = "POST"
	topic		= "topic"
	endPoint    = "/delta/delta"
	doValidationOnly = true
)

// TestUnitNewDeltaHandler asserts that the constructor for the DeltaHandler returns a fully configured handler.
func TestUnitNewDeltaHandler(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("When I call the constructor, then it returns me a valid DeltaHandler", t, func() {

		svc := sMocks.NewMockKafkaService(mockCtrl)
		h := hMocks.NewMockHelper(mockCtrl)
		chv := chvMocks.NewMockCHValidator(mockCtrl)

		config.CallValidateConfig = func(cfg *config.Config) error {
			return nil
		}
		cfg, _ := config.Get()

		deltaHandler := NewDeltaHandler(svc, h, chv, cfg, !doValidationOnly, topic)

		So(deltaHandler, ShouldNotBeNil)

		So(deltaHandler.kSvc, ShouldNotBeNil)
		So(deltaHandler.kSvc, ShouldEqual, svc)

		So(deltaHandler.chv, ShouldNotBeNil)
		So(deltaHandler.chv, ShouldEqual, chv)

		So(deltaHandler.h, ShouldNotBeNil)
		So(deltaHandler.h, ShouldEqual, h)
	})
}

// TestUnitDeltaHandlerFailsRequestBodyRetrieval asserts that when converting the request body fails, errors are
// handled correctly and returned to the user with the correct status.
func TestUnitDeltaHandlerFailsRequestBodyRetrieval(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the delta endpoint", t, func() {

		req := httptest.NewRequest(postMethod, endPoint, bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the kafka service fails to send", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)

			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewDeltaHandler(svc, h, chv, cfg, !doValidationOnly, topic)

			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, contextId).Return(nil, nil)
			h.EXPECT().GetDataFromRequest(req, contextId).Return("", errors.New("error converting request body"))
			h.EXPECT().GetRequestIdFromHeader(req).Return(contextId)

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 500 and an error returned", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
			})

		})
	})
}

// TestUnitDeltaHandlerSuccessfullySends asserts that you can send a REST request onto a kafka topic with no errors.
func TestUnitDeltaHandlerSuccessfullySends(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the delta endpoint", t, func() {

		req := httptest.NewRequest(postMethod, endPoint, bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)

			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewDeltaHandler(svc, h, chv, cfg, !doValidationOnly, topic)

			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, contextId).Return(nil, nil)
			h.EXPECT().GetDataFromRequest(req, contextId).Return(requestBody, nil)
			h.EXPECT().GetRequestIdFromHeader(req).Return(contextId)
			svc.EXPECT().SendMessage(handler.topic, requestBody, contextId).Return(nil)

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

// TestUnitDeltaHandlerFailsSend asserts that the DeltaHandler returns an internal error status when sending fails.
func TestUnitDeltaHandlerFailsSend(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the delta endpoint", t, func() {

		req := httptest.NewRequest(postMethod, endPoint, bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the kafka service fails to send", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)

			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewDeltaHandler(svc, h, chv, cfg, !doValidationOnly, topic)

			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, contextId).Return(nil, nil)
			h.EXPECT().GetDataFromRequest(req, contextId).Return(requestBody, nil)
			h.EXPECT().GetRequestIdFromHeader(req).Return(contextId)
			svc.EXPECT().SendMessage(handler.topic, requestBody, contextId).Return(errors.New("error sending message"))

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 500 and an error returned", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

// TestUnitDeltaHandlerErrorsCallingValidation asserts that the DeltaHandler returns an internal error status when
// call to validate the request fails (internal failure such as failure to open schema, not a user validation failure).
func TestUnitDeltaHandlerErrorsCallingValidation(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the delta endpoint", t, func() {

		req := httptest.NewRequest(postMethod, endPoint, bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the kafka service fails to send", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)

			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewDeltaHandler(svc, h, chv, cfg, !doValidationOnly, topic)

			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, contextId).Return(nil, errors.New("error"))
			h.EXPECT().GetRequestIdFromHeader(req).Return(contextId)

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 500 and an error returned", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

// TestUnitDeltaHandlerFailsValidation asserts that the DeltaHandler returns a bad request status when validation fails
func TestUnitDeltaHandlerFailsValidation(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the delta endpoint", t, func() {

		req := httptest.NewRequest(postMethod, endPoint, bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the request body fails validation", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)

			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewDeltaHandler(svc, h, chv, cfg, !doValidationOnly, topic)

			errBytes := []byte("error string")
			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, contextId).Return(errBytes, nil)
			h.EXPECT().GetRequestIdFromHeader(req).Return(contextId)

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 400 and an error returned", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})
}

// TestUnitDeltaHandlerValidatesOnly asserts that the DeltaHandler returns a bad request status when validation fails
func TestUnitDeltaHandlerValidatesOnly(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the delta endpoint", t, func() {

		req := httptest.NewRequest(postMethod, endPoint, bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the request body fails validation", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)

			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewDeltaHandler(svc, h, chv, cfg, doValidationOnly, topic)

			h.EXPECT().GetRequestIdFromHeader(req).Return(contextId)
			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, contextId).Return(nil, nil)
			h.EXPECT().GetDataFromRequest(req, contextId).Times(0)
			svc.EXPECT().SendMessage(handler.topic, requestBody, contextId).Times(0)

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}
