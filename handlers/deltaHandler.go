package handlers

import (
	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs-delta-api/helpers"
	"github.com/companieshouse/chs-delta-api/services"
	"github.com/companieshouse/chs-delta-api/validation"
	"github.com/companieshouse/chs.go/log"
	"net/http"
)

// DeltaHandler offers a handler by which to publish a chs-delta onto the a chosen delta kafka topic.
type DeltaHandler struct {
	kSvc             services.KafkaService
	h                helpers.Helper
	chv              validation.CHValidator
	cfg              *config.Config
	doValidationOnly bool
	isDelete         bool
	topic            string
}

// NewDeltaHandler returns an DeltaHandler.
func NewDeltaHandler(kSvc services.KafkaService, h helpers.Helper, chv validation.CHValidator,
	cfg *config.Config, doValidationOnly bool, isDelete bool, topic string) *DeltaHandler {
	return &DeltaHandler{
		kSvc:             kSvc,
		h:                h,
		chv:              chv,
		cfg:              cfg,
		doValidationOnly: doValidationOnly,
		isDelete:         isDelete,
		topic:            topic,
	}
}

// ServeHTTP accepts an incoming Delta request via a POST method, validates it and if doValidationOnly flag is set to
// false, passes it to a Kafka service for further processing along with a chosen chs-delta topic. If errors are
// encountered then they will be returned via the ResponseWriter.
func (kp *DeltaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	contextId := kp.h.GetRequestIdFromHeader(r)
	log.InfoC(contextId, "processing delta", nil)

	// Validate against the openAPI 3 spec before progressing any further.
	errValidation, err := kp.chv.ValidateRequestAgainstOpenApiSpec(r, contextId)
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

	// We only send to Kafka if doValidationOnly is false.
	if !kp.doValidationOnly {
		// Get request body and marshal into a string, ready for publishing.
		data, err := kp.h.GetDataFromRequest(r, contextId)
		if err != nil {
			log.ErrorC(contextId, err, log.Data{config.MessageKey: "error getting data from request"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Send data string to Kafka service for publishing.
		if err := kp.kSvc.SendMessage(kp.topic, data, contextId, kp.isDelete); err != nil {
			log.ErrorC(contextId, err, log.Data{config.TopicKey: kp.topic, config.MessageKey: "error sending the message to the given kafka topic"})
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	log.InfoC(contextId, "Successfully processed delta", nil)
	w.WriteHeader(http.StatusOK)
}
