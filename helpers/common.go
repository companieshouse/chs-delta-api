package helpers

import (
	"context"
	"fmt"
	"github.com/companieshouse/chs.go/log"
	"io/ioutil"
	"net/http"
)

// Used for unit testing. Allows for redirecting to stubbed functions to assert correct behaviour.
var (
	callReadAll = ioutil.ReadAll
)

const XRequestId = "X-Request-Id"

// Helper contains a list of all common functions.
type Helper interface {
	GetDataFromRequest(r *http.Request) (string, error)
	WithContextID(next http.Handler) http.Handler
}

// Impl directly implements the Helper interface.
type Impl struct {
}

// NewHelper Returns an Impl.
func NewHelper() Impl {
	return Impl{}
}

// GetDataFromRequest will try to retrieve the Body from a given request and convert it into a string.
func (h Impl) GetDataFromRequest(r *http.Request) (string, error) {

	contextId := r.Context().Value(XRequestId).(string)

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

func (h Impl) WithContextID(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		contextId, err := getRequestIdFromHeader(r)
		if err != nil {
			log.Error(err)
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, XRequestId, contextId)

		next.ServeHTTP(w, r.Clone(ctx))
	})
}

func getRequestIdFromHeader(r *http.Request) (string, error) {
	requestID := r.Header.Get(XRequestId)
	if requestID == "" {
		return "", fmt.Errorf("unable to extract request ID")
	}
	return requestID, nil
}