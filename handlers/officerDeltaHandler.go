package handlers

import (
	"fmt"
	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs-delta-api/helpers"
	"github.com/companieshouse/chs-delta-api/services"
	"github.com/companieshouse/chs-delta-api/validation"
	"github.com/companieshouse/chs.go/log"
	"net/http"
)

// OfficerDeltaHandler offers a handler by which to publish an office-delta onto the officer-delta kafka topic.
type OfficerDeltaHandler struct {
	kSvc services.KafkaService
	h helpers.Helper
	cfg *config.Config
}

// NewOfficerDeltaHandler returns an OfficerDeltaHandler.
func NewOfficerDeltaHandler(kSvc services.KafkaService, h helpers.Helper, cfg *config.Config) *OfficerDeltaHandler {
	return &OfficerDeltaHandler{kSvc: kSvc, h: h, cfg: cfg}
}

// ServeHTTP accepts an incoming OfficerDelta request via a POST method, validates it
// and then passes it to a Kafka service for further processing along with an officers-delta topic. If errors are
// encountered then they will be returned via the ResponseWriter.
func (kp *OfficerDeltaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Info(fmt.Sprintf("Open API spec to use: %s", kp.cfg.OpenApiSpec), nil)

	errValidation := validation.ValidateRequestAgainstOpenApiSpec(r, kp.cfg.OpenApiSpec)
	if errValidation != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errValidation)
		return
	}

	// Get request body.
	data, err := kp.h.GetDataFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send message to Kafka service for publishing.
	if err := kp.kSvc.SendMessage(kp.cfg.OfficerDeltaTopic, data); err != nil {
		log.Error(fmt.Errorf("error sending the message to the given kafka topic %s: %s", kp.cfg.OfficerDeltaTopic, err), nil)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
