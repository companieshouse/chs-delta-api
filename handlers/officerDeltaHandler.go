package handlers

import (
	"github.com/companieshouse/chs-delta-api/helpers"
	"github.com/companieshouse/chs-delta-api/services"
	"net/http"
)

const (
	OfficersTopic = "officers-delta"
)

// OfficerDeltaHandler offers a handler by which to publish an office-delta onto the officer-delta kafka topic.
type OfficerDeltaHandler struct {
	kSvc services.KafkaService
	h helpers.Helper
}

// NewOfficerDeltaHandler returns an OfficerDeltaHandler.
func NewOfficerDeltaHandler(kSvc services.KafkaService, h helpers.Helper) *OfficerDeltaHandler {
	return &OfficerDeltaHandler{kSvc: kSvc, h: h}
}

// ServeHTTP accepts an incoming OfficerDelta request via a POST method, validates it
// and then passes it to a Kafka service for further processing along with an officers-delta topic. If errors are
// encountered then they will be returned via the ResponseWriter.
func (kp *OfficerDeltaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Get request body.
	data, err := kp.h.GetDataFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send message to Kafka service for publishing.
	if err := kp.kSvc.SendMessage(OfficersTopic, data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
