// Package handlers contains the http handlers which receive requests to be processed by the API.
package handlers

import (
	"net/http"

	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs-delta-api/helpers"
	"github.com/companieshouse/chs-delta-api/services"
	"github.com/companieshouse/chs-delta-api/validation"
	"github.com/companieshouse/chs.go/authentication"
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
	appRouter.HandleFunc("/delta/officers", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.OfficerDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("officer-delta")
	appRouter.HandleFunc("/delta/officers/validate", NewDeltaHandler(kSvc, h, chv, cfg, true, false, cfg.OfficerDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("officer-delta-validate")
	appRouter.HandleFunc("/delta/insolvency", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.InsolvencyDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("insolvency-delta")
	appRouter.HandleFunc("/delta/insolvency/delete", NewDeltaHandler(kSvc, h, chv, cfg, false, true, cfg.InsolvencyDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("insolvency-delta")
	appRouter.HandleFunc("/delta/insolvency/validate", NewDeltaHandler(kSvc, h, chv, cfg, true, false, cfg.InsolvencyDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("insolvency-delta-validate")
	appRouter.HandleFunc("/delta/charges", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.ChargesDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("charges-delta")
	appRouter.HandleFunc("/delta/charges/delete", NewDeltaHandler(kSvc, h, chv, cfg, false, true, cfg.ChargesDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("charges-delta")
	appRouter.HandleFunc("/delta/charges/validate", NewDeltaHandler(kSvc, h, chv, cfg, true, false, cfg.ChargesDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("charges-delta-validate")
	appRouter.HandleFunc("/delta/disqualification", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.DisqualifiedDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("disqualified-officer-delta")
	appRouter.HandleFunc("/delta/disqualification/delete", NewDeltaHandler(kSvc, h, chv, cfg, false, true, cfg.DisqualifiedDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("disqualified-officer-delta")
	appRouter.HandleFunc("/delta/disqualification/validate", NewDeltaHandler(kSvc, h, chv, cfg, true, false, cfg.DisqualifiedDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("disqualified-officer-delta-validate")
	appRouter.Use(userAuthInterceptor.UserAuthenticationIntercept)
	return nil
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
