package handlers

import (
	"encoding/json"
	"github.com/companieshouse/chs-delta-api/models"
	"net/http"
)

// OfficerHandler offers a handler by which to publish a message onto a kafka topic.
type OfficerHandler struct {
}

// NewOfficerHandler returns a OfficerHandler.
func NewOfficerHandler() *OfficerHandler {
	return &OfficerHandler{
	}
}

// ServeHTTP accepts an incoming officer request via a POST method and validates and forwards to the XXX service.
func (kp *OfficerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var officersData models.OfficerDelta

	// If the decoding fails, it will most likely be due to bad data being submitted by the user.
	if err := json.NewDecoder(r.Body).Decode(&officersData) ; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}