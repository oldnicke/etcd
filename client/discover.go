package client

import (
	"go.etcd.io/etcd/pkg/srv"
)

// Discoverer is an interface that wraps the Discover method.
type Discoverer interface {
	// Discover looks up the etcd servers for the domain.
	Discover(domain string, serviceName string) ([]string, error)
}

type srvDiscover struct{}

// NewSRVDiscover constructs a new Discoverer that uses the stdlib to lookup SRV records.
func NewSRVDiscover() Discoverer {
	return &srvDiscover{}
}

func (d *srvDiscover) Discover(domain string, serviceName string) ([]string, error) {
	srvs, err := srv.GetClient("etcd-client", domain, serviceName)
	if err != nil {
		return nil, err
	}
	return srvs.Endpoints, nil
}
