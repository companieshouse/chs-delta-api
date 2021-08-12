// Package handlers contains the http handlers which receive requests to be processed by the API.
package handlers

import (
	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs-delta-api/helpers"
	"github.com/companieshouse/chs-delta-api/services"
	"github.com/companieshouse/chs-delta-api/validation"
	"github.com/companieshouse/chs.go/authentication"
	"net/http"

	"github.com/companieshouse/chs.go/log"
	"github.com/gorilla/mux"
)

var (
	callNewCHValidator = validation.NewCHValidator
)

// Register defines all REST endpoints for the API.
func Register(mainRouter *mux.Router, cfg *config.Config, kSvc services.KafkaService) error {

	// Initialise all services and components needed to run chs-delta-api correctly.
	h := helpers.NewHelper()

	// Init the CHValidator service and handle any errors that come back.
	chv, err := callNewCHValidator(cfg.OpenApiSpec)
	if err != nil {
		return err
	}

	// Init the Kafka service and handle any errors that come back.
	if err := kSvc.Init(cfg); err != nil {
		return err
	}

	userAuthInterceptor := &authentication.UserAuthenticationInterceptor{
		AllowAPIKeyUser:                true,
		RequireElevatedAPIKeyPrivilege: true,
	}

	// Register endpoints for service.
	mainRouter.HandleFunc("/delta/healthcheck", healthCheck).Methods(http.MethodGet).Name("healthcheck")
	mainRouter.Use(log.Handler)

	appRouter := mainRouter.PathPrefix("").Subrouter()
	appRouter.HandleFunc("/delta/officers", NewOfficerDeltaHandler(kSvc, h, chv, cfg, false).ServeHTTP).Methods(http.MethodPost).Name("officer-delta")
	appRouter.HandleFunc("/delta/officers/validate", NewOfficerDeltaHandler(kSvc, h, chv, cfg, true).ServeHTTP).Methods(http.MethodPost).Name("officer-delta-validate")
	appRouter.Use(userAuthInterceptor.UserAuthenticationIntercept)

	return nil
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
