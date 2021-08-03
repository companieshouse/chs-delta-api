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

// TestGetDataFromRequestSuccess asserts that a data string is returned with no errors when given a valid request.
func TestGetDataFromRequestSuccess(t *testing.T) {

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

// TestGetDataFromRequestError asserts that when reading of the request fails, it returns an empty string and error.
func TestGetDataFromRequestError(t *testing.T) {

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

// TestGetRequestIdFromHeaderError asserts that request id is not set.
func TestGetRequestIdFromHeaderError(t *testing.T) {

	Convey("Given I try to get the request id from header and X-Request-Id is not provided", t, func() {

		reqBody := http.Request{Body: ioutil.NopCloser(bytes.NewReader([]byte(requestExample)))}

		h := NewHelper()
		data, err := h.GetRequestIdFromHeader(&reqBody)

		Convey("Then I am given an error", func() {

			So(data, ShouldEqual, contextId)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unable to extract request ID")
		})
	})
}

// TestGetRequestIdFromHeaderSuccess asserts that request id is set and is successfully extracted.
func TestGetRequestIdFromHeaderSuccess(t *testing.T) {

	Convey("Given I try to get the request id from header and X-Request-Id is provided", t, func() {

		reqBody, _ := http.NewRequest(http.MethodGet, "http://www.companieshouse.gov.uk", nil)
		reqBody.Header.Set(xRequestId, contextId)

		h := NewHelper()
		data, err := h.GetRequestIdFromHeader(reqBody)

		Convey("Then I am given a request id", func() {

			So(data, ShouldEqual, contextId)
			So(err, ShouldBeNil)
		})
	})
}