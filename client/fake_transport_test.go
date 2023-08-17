package client

import (
	"errors"
	"net/http"
)

func (t *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	select {
	case resp := <-t.respchan:
		return resp, nil
	case err := <-t.errchan:
		return nil, err
	case <-t.startCancel:
	case <-req.Cancel:
	}
	select {
	// this simulates that the request is finished before cancel effects
	case resp := <-t.respchan:
		return resp, nil
	// wait on finishCancel to simulate taking some amount of
	// time while calling CancelRequest
	case <-t.finishCancel:
		return nil, errors.New("cancelled")
	}
}
