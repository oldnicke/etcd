package api

import (
	"github.com/oldnicke/etcd/etcdserver/api/membership"
	"github.com/oldnicke/etcd/pkg/types"

	"github.com/coreos/go-semver/semver"
)

// Cluster is an interface representing a collection of members in one etcd cluster.
type Cluster interface {
	// ID returns the cluster ID
	ID() types.ID
	// ClientURLs returns an aggregate set of all URLs on which this
	// cluster is listening for client requests
	ClientURLs() []string
	// Members returns a slice of members sorted by their ID
	Members() []*membership.Member
	// Member retrieves a particular member based on ID, or nil if the
	// member does not exist in the cluster
	Member(id types.ID) *membership.Member
	// Version is the cluster-wide minimum major.minor version.
	Version() *semver.Version
}
