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
	kSvc             services.KafkaService
	h                helpers.Helper
	chv              validation.CHValidator
	cfg              *config.Config
	doValidationOnly bool
}

// NewOfficerDeltaHandler returns an OfficerDeltaHandler.
func NewOfficerDeltaHandler(kSvc services.KafkaService, h helpers.Helper, chv validation.CHValidator, cfg *config.Config, doValidationOnly bool) *OfficerDeltaHandler {
	return &OfficerDeltaHandler{kSvc: kSvc, h: h, chv: chv, cfg: cfg, doValidationOnly: doValidationOnly}
}

// ServeHTTP accepts an incoming OfficerDelta request via a POST method, validates it
// and then passes it to a Kafka service for further processing along with an officers-delta topic. If errors are
// encountered then they will be returned via the ResponseWriter.
func (kp *OfficerDeltaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	contextId := kp.h.GetRequestIdFromHeader(r)

	log.InfoC(contextId, fmt.Sprintf("Using the open api spec: "), log.Data{config.OpenApiSpecKey: kp.cfg.OpenApiSpec})

	// Validate against the open API 3 spec before progressing any further.
	errValidation, err := kp.chv.ValidateRequestAgainstOpenApiSpec(r, kp.cfg.OpenApiSpec, contextId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.ErrorC(contextId, err, log.Data{config.MessageKey: "error occurred while trying to validate request"})
		return
	} else if errValidation != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errValidation)
		if err != nil {
			log.ErrorC(contextId, err, log.Data{config.MessageKey: "error occurred while trying to write response"})
		}

		return
	}

	if !kp.doValidationOnly {
		// Get request body and marshal into a string, ready for publishing.
		data, err := kp.h.GetDataFromRequest(r, contextId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Send data string to Kafka service for publishing.
		if err := kp.kSvc.SendMessage(kp.cfg.OfficerDeltaTopic, data, contextId); err != nil {
			log.ErrorC(contextId, err, log.Data{config.TopicKey: kp.cfg.OfficerDeltaTopic, config.MessageKey: "error sending the message to the given kafka topic"})
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
