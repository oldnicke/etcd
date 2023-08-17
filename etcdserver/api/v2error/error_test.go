package v2error

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestErrorWriteTo(t *testing.T) {
	for k := range errors {
		err := NewError(k, "", 1)
		rr := httptest.NewRecorder()
		err.WriteTo(rr)

		if err.StatusCode() != rr.Code {
			t.Errorf("HTTP status code %d, want %d", rr.Code, err.StatusCode())
		}

		gbody := strings.TrimSuffix(rr.Body.String(), "\n")
		if err.toJsonString() != gbody {
			t.Errorf("HTTP body %q, want %q", gbody, err.toJsonString())
		}

		wheader := http.Header(map[string][]string{
			"Content-Type": {"application/json"},
			"X-Etcd-Index": {"1"},
		})

		if !reflect.DeepEqual(wheader, rr.HeaderMap) {
			t.Errorf("HTTP headers %v, want %v", rr.HeaderMap, wheader)
		}
	}

}
