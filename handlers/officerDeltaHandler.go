package handlers

import (
	"github.com/companieshouse/chs-delta-api/services"
	"io/ioutil"
	"net/http"
)

// OfficerDeltaHandler offers a handler by which to publish a message onto a kafka topic.
type OfficerDeltaHandler struct {
	KSvc services.KafkaServiceImpl
}

// NewOfficerDeltaHandler returns an OfficerDeltaHandler.
func NewOfficerDeltaHandler(kSvc services.KafkaServiceImpl) *OfficerDeltaHandler {
	return &OfficerDeltaHandler{KSvc: kSvc}
}

// ServeHTTP accepts an incoming OfficerDelta request via a POST method, validates it
// and then passes it to a Kafka service for further processing. If errors are encountered
// then they will be returned via the ResponseWriter.
func (kp *OfficerDeltaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	deltaData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Converting body: " + err.Error()))
		return
	}
	conData := string(deltaData)
	if err := kp.KSvc.SendMessage("deltaTopic", conData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Send message: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
