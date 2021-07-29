// Package handlers contains the http handlers which receive requests to be processed by the API.
package handlers

import (
	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs-delta-api/helpers"
	"github.com/companieshouse/chs-delta-api/services"
	"github.com/companieshouse/chs-delta-api/validation"
	"net/http"

	"github.com/companieshouse/chs.go/log"
	"github.com/gorilla/mux"
)

// Register defines all REST endpoints for the API.
func Register(mainRouter *mux.Router, cfg *config.Config, kSvc services.KafkaService) error {

	// Initialise all services and components needed to run chs-delta-api correctly.
	h := helpers.NewHelper()
	chv := validation.NewCHValidator()
	if err := kSvc.Init(cfg); err != nil {
		return err
	}

	// Register endpoints for service.
	mainRouter.HandleFunc("/delta/healthcheck", healthCheck).Methods(http.MethodGet).Name("healthcheck")
	mainRouter.HandleFunc("/delta/officers", NewOfficerDeltaHandler(kSvc, h, chv, cfg).ServeHTTP).Methods(http.MethodPost).Name("officer-delta")
	mainRouter.Use(log.Handler)

	return nil
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
