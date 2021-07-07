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

// HandleOfficerDelta accepts an incoming officer request via a POST method and validates and publishes to an officer-delta kafka topic.
func (kp *OfficerHandler) HandleOfficerDelta(w http.ResponseWriter, r *http.Request) {

	var officersData models.OfficerDelta

	// If the decoding fails, it will most likely be due to bad data being submitted by the user.
	if err := json.NewDecoder(r.Body).Decode(&officersData) ; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return 
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}