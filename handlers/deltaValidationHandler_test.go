package handlers

import (
	"bytes"
	"errors"
	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs-delta-api/handlers/common"
	hMocks "github.com/companieshouse/chs-delta-api/helpers/mocks"
	chvMocks "github.com/companieshouse/chs-delta-api/validation/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	officersValidateEndpoint = "/delta/officers/validate"
)

// TestNewDeltaValidationHandler asserts that the constructor for the DeltaValidationHandler returns a fully configured handler.
func TestNewDeltaValidationHandler(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("When I call the constructor, then it returns me a valid DeltaValidationHandler", t, func() {

		h := hMocks.NewMockHelper(mockCtrl)
		chv := chvMocks.NewMockCHValidator(mockCtrl)

		config.CallValidateConfig = func(cfg *config.Config) error {
			return nil
		}
		cfg, _ := config.Get()

		deltaValidationHandler := NewDeltaValidationHandler(h, chv, cfg)

		So(deltaValidationHandler, ShouldNotBeNil)

		So(deltaValidationHandler.chv, ShouldNotBeNil)
		So(deltaValidationHandler.chv, ShouldEqual, chv)

		So(deltaValidationHandler.h, ShouldNotBeNil)
		So(deltaValidationHandler.h, ShouldEqual, h)
	})
}

// TestDeltaValidationHandlerErrorsCallingValidation asserts that the DeltaValidationHandler returns an internal error status when
// call to validate the request fails (internal failure such as failure to open schema, not a user validation failure).
func TestDeltaValidationHandlerErrorsCallingValidation(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the officer delta validate endpoint", t, func() {

		req := httptest.NewRequest(common.PostRestMethod, officersValidateEndpoint, bytes.NewBuffer([]byte(common.TestRequestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the kafka service fails to send", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)

			common.InitEnv()
			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewDeltaValidationHandler(h, chv, cfg)

			h.EXPECT().GetRequestIdFromHeader(req).Return(common.TestContextId)
			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, common.TestContextId).Return(nil, errors.New("error"))

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 500 and an error returned", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
			})

			common.DestroyEnv()
		})
	})
}

// TestDeltaValidationHandlerSuccessfulValidation asserts that the DeltaValidationHandler returns an ok status when validation passes
func TestDeltaValidationHandlerSuccessfulValidation(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the officer delta validate endpoint", t, func() {

		req := httptest.NewRequest(common.PostRestMethod, officersValidateEndpoint, bytes.NewBuffer([]byte(common.TestRequestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the kafka service fails to send", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)

			common.InitEnv()
			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewDeltaValidationHandler(h, chv, cfg)

			h.EXPECT().GetRequestIdFromHeader(req).Return(common.TestContextId)
			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, common.TestContextId).Return(nil, nil)

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 200 and no error returned", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})

			common.DestroyEnv()
		})
	})
}

// TestDeltaValidationHandlerFailsValidation asserts that the DeltaValidationHandler returns a bad request status when validation fails
func TestDeltaValidationHandlerFailsValidation(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("Given a HTTP POST request via the officer delta endpoint", t, func() {

		req := httptest.NewRequest(common.PostRestMethod, officersValidateEndpoint, bytes.NewBuffer([]byte(requestBody)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router, but the request body fails validation", func() {

			h := hMocks.NewMockHelper(mockCtrl)
			chv := chvMocks.NewMockCHValidator(mockCtrl)

			common.InitEnv()
			config.CallValidateConfig = func(cfg *config.Config) error {
				return nil
			}
			cfg, _ := config.Get()

			handler := NewDeltaValidationHandler(h, chv, cfg)

			errBytes := []byte("error string")
			h.EXPECT().GetRequestIdFromHeader(req).Return(common.TestContextId)
			chv.EXPECT().ValidateRequestAgainstOpenApiSpec(req, handler.cfg.OpenApiSpec, common.TestContextId).Return(errBytes, nil)

			handler.ServeHTTP(resp, req)

			Convey("Then the response should be 400 and an error returned", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})

			common.DestroyEnv()
		})
	})
}