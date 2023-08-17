package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	cURLDebug = false
)

func EnablecURLDebug() {
	cURLDebug = true
}

func DisablecURLDebug() {
	cURLDebug = false
}

// printcURL prints the cURL equivalent request to stderr.
// It returns an error if the body of the request cannot
// be read.
// The caller MUST cancel the request if there is an error.
func printcURL(req *http.Request) error {
	if !cURLDebug {
		return nil
	}
	var (
		command string
		b       []byte
		err     error
	)

	if req.URL != nil {
		command = fmt.Sprintf("curl -X %s %s", req.Method, req.URL.String())
	}

	if req.Body != nil {
		b, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
		command += fmt.Sprintf(" -d %q", string(b))
	}

	fmt.Fprintf(os.Stderr, "cURL Command: %s\n", command)

	// reset body
	body := bytes.NewBuffer(b)
	req.Body = ioutil.NopCloser(body)

	return nil
}
