package v2http

import (
	"fmt"
	"net/http"

	"github.com/oldnicke/etcd/etcdserver/api"
	"github.com/oldnicke/etcd/etcdserver/api/v2http/httptypes"
)

func authCapabilityHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !api.IsCapabilityEnabled(api.AuthCapability) {
			notCapable(w, r, api.AuthCapability)
			return
		}
		fn(w, r)
	}
}

func notCapable(w http.ResponseWriter, r *http.Request, c api.Capability) {
	herr := httptypes.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Not capable of accessing %s feature during rolling upgrades.", c))
	if err := herr.WriteTo(w); err != nil {
		plog.Debugf("error writing HTTPError (%v) to %s", err, r.RemoteAddr)
	}
}
