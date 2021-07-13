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
)

// TestGetDataFromRequestSuccess asserts that a data string is returned with no errors when given a valid request.
func TestGetDataFromRequestSuccess(t *testing.T) {

	Convey("Given I pass a request into the GetDataFromRequest function", t, func() {
		h := NewHelper()
		reqBody := http.Request{Body: ioutil.NopCloser(bytes.NewReader([]byte(requestExample)))}
		data, err := h.GetDataFromRequest(&reqBody)

		Convey("Then I am given a string version of my request back with no error", func() {
			So(data, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}

//Used in error test for asserting error occurs
type mockRequestError struct{}

// TestGetDataFromRequestError asserts that when reading of the request fails, it returns an empty string and error.
func TestGetDataFromRequestError(t *testing.T) {

	Convey("Given I pass a request into the GetDataFromRequest function", t, func() {

		callReadAll = func (r io.Reader) ([]byte, error) {
			return nil, errors.New("error getting data from request")
		}

		reqBody := http.Request{Body: ioutil.NopCloser(bytes.NewReader([]byte(requestExample)))}

		h := NewHelper()
		data, err := h.GetDataFromRequest(&reqBody)

		Convey("Then I am given a an error", func() {

			So(data, ShouldEqual, "")
			So(err, ShouldNotBeNil)
		})
	})
}
