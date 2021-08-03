package helpers

import (
	"fmt"
	"github.com/companieshouse/chs.go/log"
	"io/ioutil"
	"net/http"
)

// Used for unit testing. Allows for redirecting to stubbed functions to assert correct behaviour.
var (
	callReadAll = ioutil.ReadAll
)

const xRequestId = "X-Request-Id"

// Helper contains a list of all common functions.
type Helper interface {
	GetDataFromRequest(r *http.Request, contextId string) (string, error)
	GetRequestIdFromHeader(r *http.Request) string
}

// Impl directly implements the Helper interface.
type Impl struct {
}

// NewHelper Returns an Impl.
func NewHelper() Impl {
	return Impl{}
}

// GetDataFromRequest will try to retrieve the Body from a given request and convert it into a string.
func (h Impl) GetDataFromRequest(r *http.Request, contextId string) (string, error) {

	// Retrieve the request body.
	data, err := callReadAll(r.Body)
	if err != nil {
		log.ErrorC(contextId, fmt.Errorf("error while retrieving the Body from a given request and converting it into a string : %s", err), nil)
		return "", err
	}

	// Convert the request body into a string and pass it to the Kafka Service for publishing.
	strData := string(data)
	return strData, nil
}

//GetRequestIdFromHeader gets X-Request-Id from header and use it as a context id for logging
func (h Impl) GetRequestIdFromHeader(r *http.Request) string {
	requestID := r.Header.Get(xRequestId)
	if requestID == "" {
		log.Error(fmt.Errorf("unable to extract request ID"))
		return "contextId"
	}
	return requestID
}
