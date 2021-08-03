package helpers

import (
	"fmt"
	"github.com/companieshouse/chs-delta-api/config"
	"github.com/companieshouse/chs.go/log"
	"io/ioutil"
	"net/http"
	mrand "math/rand"
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
		log.ErrorC(contextId, err, log.Data{config.MessageKey : "error while retrieving the Body from a given request and converting it into a string"})

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
		return generateContextId()
	}
	return requestID
}

func generateContextId() string {
	urlbase64 := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")
	b := make([]rune, 28)
	for i := range b {
		b[i] = urlbase64[mrand.Intn(len(urlbase64))]
	}
	return string(b)
}