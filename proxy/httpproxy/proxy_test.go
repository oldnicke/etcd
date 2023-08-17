package httpproxy

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestReadonlyHandler(t *testing.T) {
	fixture := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	hdlrFunc := readonlyHandlerFunc(http.HandlerFunc(fixture))

	tests := []struct {
		method string
		want   int
	}{
		// GET is only passing method
		{"GET", http.StatusOK},

		// everything but GET is StatusNotImplemented
		{"POST", http.StatusNotImplemented},
		{"PUT", http.StatusNotImplemented},
		{"PATCH", http.StatusNotImplemented},
		{"DELETE", http.StatusNotImplemented},
		{"FOO", http.StatusNotImplemented},
	}

	for i, tt := range tests {
		req, _ := http.NewRequest(tt.method, "http://example.com", nil)
		rr := httptest.NewRecorder()
		hdlrFunc(rr, req)

		if tt.want != rr.Code {
			t.Errorf("#%d: incorrect HTTP status code: method=%s want=%d got=%d", i, tt.method, tt.want, rr.Code)
		}
	}
}

func TestConfigHandlerGET(t *testing.T) {
	var err error
	us := make([]*url.URL, 3)
	us[0], err = url.Parse("http://example1.com")
	if err != nil {
		t.Fatal(err)
	}
	us[1], err = url.Parse("http://example2.com")
	if err != nil {
		t.Fatal(err)
	}
	us[2], err = url.Parse("http://example3.com")
	if err != nil {
		t.Fatal(err)
	}

	rp := reverseProxy{
		director: &director{
			ep: []*endpoint{
				newEndpoint(*us[0], 1*time.Second),
				newEndpoint(*us[1], 1*time.Second),
				newEndpoint(*us[2], 1*time.Second),
			},
		},
	}

	req, _ := http.NewRequest("GET", "http://example.com//v2/config/local/proxy", nil)
	rr := httptest.NewRecorder()
	rp.configHandler(rr, req)

	wbody := "{\"endpoints\":[\"http://example1.com\",\"http://example2.com\",\"http://example3.com\"]}\n"

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != wbody {
		t.Errorf("body = %s, want %s", string(body), wbody)
	}
}
