package helpers

import (
	"io/ioutil"
	"net/http"
)

// Used for unit testing. Allows for redirecting to stubbed functions to assert correct behaviour.
var (
	callReadAll = ioutil.ReadAll
)

// Helper contains a list of all common functions.
type Helper interface {
	GetDataFromRequest(r *http.Request) (string, error)
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

	// Retrieve the request body.
	data, err := callReadAll(r.Body)
	if err != nil {
		return "", err
	}

	// Convert the request body into a string and pass it to the Kafka Service for publishing.
	strData := string(data)
	return strData, nil
}