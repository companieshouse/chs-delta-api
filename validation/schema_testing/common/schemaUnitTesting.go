// common contains common functions which can be used when unit testing schemas to make them more compact and easy to implement.
package common

import (
	"bytes"
	"encoding/json"
	"github.com/companieshouse/chs-delta-api/models"
	"net/http"
	"os"
)

const (
	xRequestId      = "X-Request-Id"
	contentType     = "Content-Type"
	applicationJson = "application/json"
	contextId       = "contextId"
)

// SetHeaders sets the required headers for the unit testing schemas to correctly function.
// We require a contextId to be set in the Header of the request as it will be later pulled and used as a contextId.
// We require a contentType to be set (application/json) for the kin-openAPI library to validate the request.
func SetHeaders(r *http.Request) *http.Request {
	r.Header.Set(xRequestId, contextId)
	r.Header.Set(contentType, applicationJson)

	return r
}

// ReadRequestBody takes a file location in the format of a string relative path and reads in the contents from the file,
// removing any extra tabs, spaces and special characters using json.Compact(). Returns a []byte version of the formatted file read in.
func ReadRequestBody(fl string) []byte {

	// Read in contents and covert it to a string for further processing.
	raw, _ := os.ReadFile(fl)

	buffer := new(bytes.Buffer)
	_ = json.Compact(buffer, raw)

	raw = buffer.Bytes()
	// Convert back to an []byte and return.
	return raw
}

// CompareActualToExpected takes actual and expected json (as byte arrays) and compares them to see if they match. Ordering of response
// isn't always guaranteed when calling the kin-openAPI validator so using this function allows you to match the response
// the library gives you with an expected response without worrying about ordering.
func CompareActualToExpected(actual, expected []byte) bool {

	// Define 2 model CHError arrays to hold actual and expected responses.
	var actualErrArr *[]models.CHError
	var expectedErrArr *[]models.CHError

	// Convert responses to JSON object arrays.
	_ = json.Unmarshal(actual, &actualErrArr)
	_ = json.Unmarshal(expected, &expectedErrArr)

	// If the arrays don't match in length, they can't be a complete match so return false.
	if len(*actualErrArr) != len(*expectedErrArr) {
		return false
	}

	// create a map of CHError.Location (string) -> int to compare if both arrays are completely equal.
	diff := make(map[string]int, len(*expectedErrArr))

	// Range over the expected response array and add them to the newly created map.
	for _, ee := range *expectedErrArr {
		// 0 value for int is 0, so just increment a counter for the string
		diff[ee.Location]++
	}

	// Finally range over the actual response array errors and check that they exist in the expected errors map.
	for _, ae := range *actualErrArr {
		// If the error is not in diff bail out early as they can't be a match.
		if _, ok := diff[ae.Location]; !ok {
			return false
		}
		diff[ae.Location] -= 1
		if diff[ae.Location] == 0 {
			delete(diff, ae.Location)
		}
	}

	// Return if length of remaining map values is equal to 0. If it is, we have a match.
	return len(diff) == 0
}
