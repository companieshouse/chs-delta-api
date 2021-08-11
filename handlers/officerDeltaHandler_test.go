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
	"os"
	"testing"
)

const (
	requestBody = `{"dummy" : "request"}`
	contextId   = "contextId"
	postMethod  = "POST"
	endPoint    = "/delta/officers"
)

func initEnv() {
	_ = os.Setenv("BIND_ADDR", "bind_addr")
	_ = os.Setenv("KAFKA_BROKER_ADDR", "kafka_broker_addr,kafka_broker_addr")
	_ = os.Setenv("SCHEMA_REGISTRY_URL", "schema_registry_url")
	_ = os.Setenv("OFFICER_DELTA_TOPIC", "officer_delta_topic")
	_ = os.Setenv("OPEN_API_SPEC", "open_api_spec")
}

func destroyEnv() {
	os.Clearenv()
}

// TestNewOfficerDeltaHandler asserts that the constructor for the OfficerDeltaHandler returns a fully configured handler.
func TestNewOfficerDeltaHandler(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("When I call the constructor, then it returns me a valid OfficerDeltaHandler", t, func() {

		svc := sMocks.NewMockKafkaService(mockCtrl)
		h := hMocks.NewMockHelper(mockCtrl)
		chv := chvMocks.NewMockCHValidator(mockCtrl)
		doValidationOnly := false

		config.CallValidateConfig = func(cfg *config.Config) error {
			return nil
		}
		cfg, _ := config.Get()

		officerHandler := NewOfficerDeltaHandler(svc, h, chv, cfg, doValidationOnly)

		So(officerHandler, ShouldNotBeNil)

		So(officerHandler.kSvc, ShouldNotBeNil)
		So(officerHandler.kSvc, ShouldEqual, svc)

		So(officerHandler.chv, ShouldNotBeNil)
		So(officerHandler.chv, ShouldEqual, chv)

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

		req := httptest.NewRequest(postMethod, endPoint, bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the kafka service fails to send", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)
			doValidationOnly := false

			initEnv()
			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewOfficerDeltaHandler(svc, h, chv, cfg, doValidationOnly)

			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, contextId).Return(nil, nil)
			h.EXPECT().GetDataFromRequest(req, contextId).Return("", errors.New("error converting request body"))
			h.EXPECT().GetRequestIdFromHeader(req).Return(contextId)

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 500 and an error returned", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
			})

			destroyEnv()
		})
	})
}

// TestOfficerDeltaHandlerSuccessfullySends asserts that you can send a REST request onto a kafka topic with no errors.
func TestOfficerDeltaHandlerSuccessfullySends(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the officer delta endpoint", t, func() {

		req := httptest.NewRequest(postMethod, endPoint, bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)
			doValidationOnly := false

			initEnv()
			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewOfficerDeltaHandler(svc, h, chv, cfg, doValidationOnly)

			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, contextId).Return(nil, nil)
			h.EXPECT().GetDataFromRequest(req, contextId).Return(requestBody, nil)
			h.EXPECT().GetRequestIdFromHeader(req).Return(contextId)
			svc.EXPECT().SendMessage(handler.cfg.OfficerDeltaTopic, requestBody, contextId).Return(nil)

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})

			destroyEnv()
		})
	})
}

// TestOfficerDeltaHandlerFailsSend asserts that the officerDeltaHandler returns an internal error status when sending fails.
func TestOfficerDeltaHandlerFailsSend(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the officer delta endpoint", t, func() {

		req := httptest.NewRequest(postMethod, endPoint, bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the kafka service fails to send", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)
			doValidationOnly := false

			initEnv()
			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewOfficerDeltaHandler(svc, h, chv, cfg, doValidationOnly)

			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, contextId).Return(nil, nil)
			h.EXPECT().GetDataFromRequest(req, contextId).Return(requestBody, nil)
			h.EXPECT().GetRequestIdFromHeader(req).Return(contextId)
			svc.EXPECT().SendMessage(handler.cfg.OfficerDeltaTopic, requestBody, contextId).Return(errors.New("error sending message"))

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 500 and an error returned", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
			})

			destroyEnv()
		})
	})
}

// TestOfficerDeltaHandlerErrorsCallingValidation asserts that the officerDeltaHandler returns an internal error status when
// call to validate the request fails (internal failure such as failure to open schema, not a user validation failure).
func TestOfficerDeltaHandlerErrorsCallingValidation(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the officer delta endpoint", t, func() {

		req := httptest.NewRequest(postMethod, endPoint, bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the kafka service fails to send", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)
			doValidationOnly := false

			initEnv()
			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewOfficerDeltaHandler(svc, h, chv, cfg, doValidationOnly)

			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, contextId).Return(nil, errors.New("error"))
			h.EXPECT().GetRequestIdFromHeader(req).Return(contextId)

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 500 and an error returned", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
			})

			destroyEnv()
		})
	})
}

// TestOfficerDeltaHandlerFailsValidation asserts that the officerDeltaHandler returns a bad request status when validation fails
func TestOfficerDeltaHandlerFailsValidation(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the officer delta endpoint", t, func() {

		req := httptest.NewRequest(postMethod, endPoint, bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the request body fails validation", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)
			doValidationOnly := false

			initEnv()
			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewOfficerDeltaHandler(svc, h, chv, cfg, doValidationOnly)

			errBytes := []byte("error string")
			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, contextId).Return(errBytes, nil)
			h.EXPECT().GetRequestIdFromHeader(req).Return(contextId)

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 400 and an error returned", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})

			destroyEnv()
		})
	})
}

// TestOfficerDeltaHandlerValidatesOnly asserts that the officerDeltaHandler returns a bad request status when validation fails
func TestOfficerDeltaHandlerValidatesOnly(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the officer delta endpoint", t, func() {

		req := httptest.NewRequest(postMethod, endPoint, bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the request body fails validation", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			svc := sMocks.NewMockKafkaService(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)
			doValidationOnly := true

			initEnv()
			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewOfficerDeltaHandler(svc, h, chv, cfg, doValidationOnly)

			h.EXPECT().GetRequestIdFromHeader(req).Return(contextId)
			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, contextId).Return(nil, nil)
			h.EXPECT().GetDataFromRequest(req, contextId).Times(0)
			svc.EXPECT().SendMessage(handler.cfg.OfficerDeltaTopic, requestBody, contextId).Times(0)

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 200", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})

			destroyEnv()
		})
	})
}