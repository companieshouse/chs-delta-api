package handlers

import (
	"github.com/companieshouse/chs-delta-api/helpers"
	"github.com/companieshouse/chs-delta-api/services"
	"net/http"
)

const (
	officersTopic = "officers-delta"
)

// OfficerDeltaHandler offers a handler by which to publish an office-delta onto the officer-delta kafka topic.
type OfficerDeltaHandler struct {
	KSvc services.KafkaService
	h helpers.Helper
}

// NewOfficerDeltaHandler returns an OfficerDeltaHandler.
func NewOfficerDeltaHandler(kSvc services.KafkaService, h helpers.Helper) *OfficerDeltaHandler {
	return &OfficerDeltaHandler{KSvc: kSvc, h: h}
}

// ServeHTTP accepts an incoming OfficerDelta request via a POST method, validates it
// and then passes it to a Kafka service for further processing along with an officers-delta topic. If errors are
// encountered then they will be returned via the ResponseWriter.
func (kp *OfficerDeltaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Get request body.
	data, err := kp.h.GetDataFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error retrieving request body: " + err.Error())) // TODO: Temp until we add the CH errors object
		return
	}

	// Send message to Kafka service for publishing.
	if err := kp.KSvc.SendMessage(officersTopic, data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Send message: " + err.Error())) // TODO: Temp until we add the CH errors object
		return
	}

	w.WriteHeader(http.StatusOK)
}
