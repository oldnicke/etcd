package httptypes

import (
	"encoding/json"
	"net/http"

	"github.com/coreos/pkg/capnslog"
)

var (
	plog = capnslog.NewPackageLogger("oldnicke/etcd", "etcdserver/api/v2http/httptypes")
)

type HTTPError struct {
	Message string `json:"message"`
	// Code is the HTTP status code
	Code int `json:"-"`
}

func (e HTTPError) Error() string {
	return e.Message
}

func (e HTTPError) WriteTo(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	b, err := json.Marshal(e)
	if err != nil {
		plog.Panicf("marshal HTTPError should never fail (%v)", err)
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

func NewHTTPError(code int, m string) *HTTPError {
	return &HTTPError{
		Message: m,
		Code:    code,
	}
}
