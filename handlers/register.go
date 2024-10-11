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
	appRouter.HandleFunc("/delta/officers", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.OfficerDeltaTopic, "internal_id").ServeHTTP).Methods(http.MethodPost).Name("officer-delta")
	appRouter.HandleFunc("/delta/officers/delete", NewDeltaHandler(kSvc, h, chv, cfg, false, true, cfg.OfficerDeltaTopic, "internal_id").ServeHTTP).Methods(http.MethodPost).Name("officer-delta")
	appRouter.HandleFunc("/delta/officers/validate", NewDeltaHandlerValidate(kSvc, h, chv, cfg, true, false, cfg.OfficerDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("officer-delta-validate")
	appRouter.HandleFunc("/delta/insolvency", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.InsolvencyDeltaTopic, "company_number").ServeHTTP).Methods(http.MethodPost).Name("insolvency-delta")
	appRouter.HandleFunc("/delta/insolvency/delete", NewDeltaHandler(kSvc, h, chv, cfg, false, true, cfg.InsolvencyDeltaTopic, "company_number").ServeHTTP).Methods(http.MethodPost).Name("insolvency-delta")
	appRouter.HandleFunc("/delta/insolvency/validate", NewDeltaHandlerValidate(kSvc, h, chv, cfg, true, false, cfg.InsolvencyDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("insolvency-delta-validate")
	appRouter.HandleFunc("/delta/charges", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.ChargesDeltaTopic, "id").ServeHTTP).Methods(http.MethodPost).Name("charges-delta")
	appRouter.HandleFunc("/delta/charges/delete", NewDeltaHandler(kSvc, h, chv, cfg, false, true, cfg.ChargesDeltaTopic, "charges_id").ServeHTTP).Methods(http.MethodPost).Name("charges-delta")
	appRouter.HandleFunc("/delta/charges/validate", NewDeltaHandlerValidate(kSvc, h, chv, cfg, true, false, cfg.ChargesDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("charges-delta-validate")
	appRouter.HandleFunc("/delta/disqualification", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.DisqualifiedDeltaTopic, "officer_id").ServeHTTP).Methods(http.MethodPost).Name("disqualified-officer-delta")
	appRouter.HandleFunc("/delta/disqualification/delete", NewDeltaHandler(kSvc, h, chv, cfg, false, true, cfg.DisqualifiedDeltaTopic, "officer_id").ServeHTTP).Methods(http.MethodPost).Name("disqualified-officer-delta")
	appRouter.HandleFunc("/delta/disqualification/validate", NewDeltaHandlerValidate(kSvc, h, chv, cfg, true, false, cfg.DisqualifiedDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("disqualified-officer-delta-validate")
	appRouter.HandleFunc("/delta/company", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.CompanyDeltaTopic, "company_number").ServeHTTP).Methods(http.MethodPost).Name("company-delta")
	appRouter.HandleFunc("/delta/company/delete", NewDeltaHandler(kSvc, h, chv, cfg, false, true, cfg.CompanyDeltaTopic, "company_number").ServeHTTP).Methods(http.MethodPost).Name("company-delta")
	appRouter.HandleFunc("/delta/company/validate", NewDeltaHandlerValidate(kSvc, h, chv, cfg, true, false, cfg.CompanyDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("company-delta-validate")
	appRouter.HandleFunc("/delta/exemption", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.ExemptionDeltaTopic, "company_number").ServeHTTP).Methods(http.MethodPost).Name("exemption-delta")
	appRouter.HandleFunc("/delta/exemption/delete", NewDeltaHandler(kSvc, h, chv, cfg, false, true, cfg.ExemptionDeltaTopic, "company_number").ServeHTTP).Methods(http.MethodPost).Name("exemption-delta")
	appRouter.HandleFunc("/delta/exemption/validate", NewDeltaHandlerValidate(kSvc, h, chv, cfg, true, false, cfg.ExemptionDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("exemption-delta-validate")
	appRouter.HandleFunc("/delta/psc-statement", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.PscStatementDeltaTopic, "psc_statement_id").ServeHTTP).Methods(http.MethodPost).Name("psc-statement-delta")
	appRouter.HandleFunc("/delta/psc-statement/validate", NewDeltaHandlerValidate(kSvc, h, chv, cfg, true, false, cfg.PscStatementDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("psc-statement-delta-validate")
	appRouter.HandleFunc("/delta/psc-statement/delete", NewDeltaHandler(kSvc, h, chv, cfg, false, true, cfg.PscStatementDeltaTopic, "psc_statement_id").ServeHTTP).Methods(http.MethodPost).Name("psc-statement-delta")
	appRouter.HandleFunc("/delta/pscs", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.PscDeltaTopic, "psc_id").ServeHTTP).Methods(http.MethodPost).Name("psc-delta")
	appRouter.HandleFunc("/delta/pscs/validate", NewDeltaHandlerValidate(kSvc, h, chv, cfg, true, false, cfg.PscDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("psc-delta-validate")
	appRouter.HandleFunc("/delta/pscs/delete", NewDeltaHandler(kSvc, h, chv, cfg, false, true, cfg.PscDeltaTopic, "psc_id").ServeHTTP).Methods(http.MethodPost).Name("psc-delta-delete")
	appRouter.HandleFunc("/delta/filing-history", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.FilingHistoryDeltaTopic, "entity_id").ServeHTTP).Methods(http.MethodPost).Name("filing-history-delta")
	appRouter.HandleFunc("/delta/filing-history/validate", NewDeltaHandlerValidate(kSvc, h, chv, cfg, true, false, cfg.FilingHistoryDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("filing-history-delta-validate")
	appRouter.HandleFunc("/delta/filing-history/delete", NewDeltaHandler(kSvc, h, chv, cfg, false, true, cfg.FilingHistoryDeltaTopic, "entity_id").ServeHTTP).Methods(http.MethodPost).Name("filing-history-delete-delta")
	// appRouter.HandleFunc("/delta/document-store", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.DocumentStoreDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("document-store-delta")
	// appRouter.HandleFunc("/delta/document-store/validate", NewDeltaHandler(kSvc, h, chv, cfg, true, false, cfg.DocumentStoreDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("document-store-delta-validate")
	appRouter.HandleFunc("/delta/registers", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.RegistersDeltaTopic, "company_number").ServeHTTP).Methods(http.MethodPost).Name("registers-delta")
	appRouter.HandleFunc("/delta/registers/delete", NewDeltaHandler(kSvc, h, chv, cfg, false, true, cfg.RegistersDeltaTopic, "company_number").ServeHTTP).Methods(http.MethodPost).Name("registers-delta-delete")
	appRouter.HandleFunc("/delta/registers/validate", NewDeltaHandlerValidate(kSvc, h, chv, cfg, true, false, cfg.RegistersDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("registers-delta-validate")
	appRouter.HandleFunc("/delta/acsp", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.AcspProfileDeltaTopic, "acsp_number").ServeHTTP).Methods(http.MethodPost).Name("acsp-profile-delta")
	appRouter.HandleFunc("/delta/acsp/validate", NewDeltaHandlerValidate(kSvc, h, chv, cfg, true, false, cfg.AcspProfileDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("acsp-profile-delta-validate")
	appRouter.Use(userAuthInterceptor.UserAuthenticationIntercept)

	// TODO: move these back to appRouter when CHIPS image-sender service has been updated to allow an aPI key to be configured to its calls here
	mainRouter.HandleFunc("/delta/document-store", NewDeltaHandler(kSvc, h, chv, cfg, false, false, cfg.DocumentStoreDeltaTopic, "transaction_id").ServeHTTP).Methods(http.MethodPost).Name("document-store-delta")
	mainRouter.HandleFunc("/delta/document-store/validate", NewDeltaHandlerValidate(kSvc, h, chv, cfg, true, false, cfg.DocumentStoreDeltaTopic).ServeHTTP).Methods(http.MethodPost).Name("document-store-delta-validate")

	return nil
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
