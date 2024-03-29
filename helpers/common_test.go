package helpers

import (
	"bytes"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	requestExample = `{"test": "example"}`
	contextId      = "contextId"
)

// TestUnitGetDataFromRequestSuccess asserts that a data string is returned with no errors when given a valid request.
func TestUnitGetDataFromRequestSuccess(t *testing.T) {

	Convey("Given I pass a request into the GetDataFromRequest function", t, func() {
		h := NewHelper()
		reqBody := http.Request{Body: ioutil.NopCloser(bytes.NewReader([]byte(requestExample)))}
		data, err := h.GetDataFromRequest(&reqBody, contextId)

		Convey("Then I am given a string version of my request back with no error", func() {
			So(data, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}

// TestUnitGetDataFromRequestError asserts that when reading of the request fails, it returns an empty string and error.
func TestUnitGetDataFromRequestError(t *testing.T) {

	Convey("Given I pass a request into the GetDataFromRequest function", t, func() {

		callReadAll = func(r io.Reader) ([]byte, error) {
			return nil, errors.New("error getting data from request")
		}

		reqBody := http.Request{Body: ioutil.NopCloser(bytes.NewReader([]byte(requestExample)))}

		h := NewHelper()
		data, err := h.GetDataFromRequest(&reqBody, contextId)

		Convey("Then I am given a an error", func() {

			So(data, ShouldEqual, "")
			So(err, ShouldNotBeNil)
		})
	})
}

// TestUnitGetRequestIdFromHeaderError asserts that request id is not set.
func TestUnitGetRequestIdFromHeaderError(t *testing.T) {

	Convey("Given I try to get the request id from header and X-Request-Id is not provided", t, func() {

		reqBody := http.Request{Body: ioutil.NopCloser(bytes.NewReader([]byte(requestExample)))}

		h := NewHelper()
		data := h.GetRequestIdFromHeader(&reqBody)

		Convey("Then I am given an error", func() {

			So(data, ShouldNotBeNil)
			So(data, ShouldNotEqual, contextId)
		})
	})
}

// TestUnitGetRequestIdFromHeaderSuccess asserts that request id is set and is successfully extracted.
func TestUnitGetRequestIdFromHeaderSuccess(t *testing.T) {

	Convey("Given I try to get the request id from header and X-Request-Id is provided", t, func() {

		reqBody, _ := http.NewRequest(http.MethodGet, "http://www.companieshouse.gov.uk", nil)
		reqBody.Header.Set(xRequestId, contextId)

		h := NewHelper()
		data := h.GetRequestIdFromHeader(reqBody)

		Convey("Then I am given a request id", func() {

			So(data, ShouldNotBeNil)
			So(data, ShouldEqual, contextId)

		})
	})
}
