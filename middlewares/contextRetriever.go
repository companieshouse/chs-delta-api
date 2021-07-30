package middlewares

import (
	"context"
	"fmt"
	"github.com/companieshouse/chs.go/log"
	"net/http"
)
var (
	callGetRequestFromHeader = getRequestIdFromHeader
)
const XRequestId = "X-Request-Id"

func  WithContextID(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		contextId, err := callGetRequestFromHeader(r)
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
