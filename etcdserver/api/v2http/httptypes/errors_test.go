package httptypes

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHTTPErrorWriteTo(t *testing.T) {
	err := NewHTTPError(http.StatusBadRequest, "what a bad request you made!")
	rr := httptest.NewRecorder()
	if e := err.WriteTo(rr); e != nil {
		t.Fatalf("HTTPError.WriteTo error (%v)", e)
	}

	wcode := http.StatusBadRequest
	wheader := http.Header(map[string][]string{
		"Content-Type": {"application/json"},
	})
	wbody := `{"message":"what a bad request you made!"}`

	if wcode != rr.Code {
		t.Errorf("HTTP status code %d, want %d", rr.Code, wcode)
	}

	if !reflect.DeepEqual(wheader, rr.HeaderMap) {
		t.Errorf("HTTP headers %v, want %v", rr.HeaderMap, wheader)
	}

	gbody := rr.Body.String()
	if wbody != gbody {
		t.Errorf("HTTP body %q, want %q", gbody, wbody)
	}
}
