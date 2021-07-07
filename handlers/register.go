// Package handlers contains the http handlers which receive requests to be processed by the API.
package handlers

import (
	"net/http"

	"github.com/companieshouse/chs.go/log"
	"github.com/gorilla/mux"
)

// Register defines all REST endpoints for the API.
func Register(mainRouter *mux.Router) {
	mainRouter.HandleFunc("/delta/healthcheck", healthCheck).Methods(http.MethodGet).Name("healthcheck")
	mainRouter.HandleFunc("/delta/officer-delta", NewOfficerDeltaHandler().ServeHTTP).Methods(http.MethodPost).Name("officers")
	mainRouter.Use(log.Handler)
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
