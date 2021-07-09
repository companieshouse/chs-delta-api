package helpers

import (
	"io/ioutil"
	"net/http"
)

// GetDataFromRequest will try to retrieve the Body from a given request and convert it into a string.
func GetDataFromRequest(r *http.Request) (string, error) {

	// Retrieve the request body.
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	// Convert the request body into a string and pass it to the Kafka Service for publishing.
	strData := string(data)
	return strData, nil
}