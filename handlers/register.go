// Package handlers contains the http handlers which receive requests to be processed by the API.
package handlers

import (
	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs-delta-api/services"
	"net/http"

	"github.com/companieshouse/chs.go/log"
	"github.com/gorilla/mux"
)

// Register defines all REST endpoints for the API.
func Register(mainRouter *mux.Router, cfg *config.Config) error {

	kSvc, err := services.NewKafkaService(cfg)
	if err != nil {
		return err
	}

	mainRouter.HandleFunc("/delta/healthcheck", healthCheck).Methods(http.MethodGet).Name("healthcheck")
	mainRouter.HandleFunc("/delta/officers", NewOfficerDeltaHandler().ServeHTTP).Methods(http.MethodPost).Name("officer-delta")
	mainRouter.HandleFunc("/delta/officer-delta", NewOfficerDeltaHandler(kSvc).ServeHTTP).Methods(http.MethodPost).Name("officers")
	mainRouter.Use(log.Handler)

	return nil
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
