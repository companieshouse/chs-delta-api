package handlers

import (
	"net/http"
)

// OfficerDeltaHandler offers a handler by which to publish a message onto a kafka topic.
type OfficerDeltaHandler struct {
}

// NewOfficerDeltaHandler returns an OfficerDeltaHandler.
func NewOfficerDeltaHandler() *OfficerDeltaHandler {
	return &OfficerDeltaHandler{}
}

// ServeHTTP accepts an incoming OfficerDelta request via a POST method, validates it
// and then passes it to a Kafka service for further processing. If errors are encountered
// then they will be returned via the ResponseWriter.
func (kp *OfficerDeltaHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {

	w.WriteHeader(http.StatusOK)
}