package httpproxy

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/http2"
)

const (
	// DefaultMaxIdleConnsPerHost indicates the default maximum idle connection
	// count maintained between proxy and each member. We set it to 128 to
	// let proxy handle 128 concurrent requests in long term smoothly.
	// If the number of concurrent requests is bigger than this value,
	// proxy needs to create one new connection when handling each request in
	// the delta, which is bad because the creation consumes resource and
	// may eat up ephemeral ports.
	DefaultMaxIdleConnsPerHost = 128
)

// GetProxyURLs is a function which should return the current set of URLs to
// which client requests should be proxied. This function will be queried
// periodically by the proxy Handler to refresh the set of available
// backends.
type GetProxyURLs func() []string

// NewHandler creates a new HTTP handler, listening on the given transport,
// which will proxy requests to an etcd cluster.
// The handler will periodically update its view of the cluster.
func NewHandler(t *http.Transport, urlsFunc GetProxyURLs, failureWait time.Duration, refreshInterval time.Duration) http.Handler {
	if t.TLSClientConfig != nil {
		// Enable http2, see Issue 5033.
		err := http2.ConfigureTransport(t)
		if err != nil {
			plog.Infof("Error enabling Transport HTTP/2 support: %v", err)
		}
	}

	p := &reverseProxy{
		director:  newDirector(urlsFunc, failureWait, refreshInterval),
		transport: t,
	}

	mux := http.NewServeMux()
	mux.Handle("/", p)
	mux.HandleFunc("/v2/config/local/proxy", p.configHandler)

	return mux
}

// NewReadonlyHandler wraps the given HTTP handler to allow only GET requests
func NewReadonlyHandler(hdlr http.Handler) http.Handler {
	readonly := readonlyHandlerFunc(hdlr)
	return http.HandlerFunc(readonly)
}

func readonlyHandlerFunc(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		next.ServeHTTP(w, req)
	}
}

func (p *reverseProxy) configHandler(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r.Method, "GET") {
		return
	}

	eps := p.director.endpoints()
	epstr := make([]string, len(eps))
	for i, e := range eps {
		epstr[i] = e.URL.String()
	}

	proxyConfig := struct {
		Endpoints []string `json:"endpoints"`
	}{
		Endpoints: epstr,
	}

	json.NewEncoder(w).Encode(proxyConfig)
}

// allowMethod verifies that the given method is one of the allowed methods,
// and if not, it writes an error to w.  A boolean is returned indicating
// whether or not the method is allowed.
func allowMethod(w http.ResponseWriter, m string, ms ...string) bool {
	for _, meth := range ms {
		if m == meth {
			return true
		}
	}
	w.Header().Set("Allow", strings.Join(ms, ","))
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	return false
}
