// Package version implements etcd version parsing and contains latest version
// information.
package version

import (
	"fmt"
	"strings"

	"github.com/coreos/go-semver/semver"
)

var (
	// MinClusterVersion is the min cluster version this etcd binary is compatible with.
	MinClusterVersion = "3.0.0"
	Version           = "3.4.13"
	APIVersion        = "unknown"

	// Git SHA Value will be set during build
	GitSHA = "Not provided (use ./build instead of go build)"
)

func init() {
	ver, err := semver.NewVersion(Version)
	if err == nil {
		APIVersion = fmt.Sprintf("%d.%d", ver.Major, ver.Minor)
	}
}

type Versions struct {
	Server  string `json:"etcdserver"`
	Cluster string `json:"etcdcluster"`
	// TODO: raft state machine version
}

// Cluster only keeps the major.minor.
func Cluster(v string) string {
	vs := strings.Split(v, ".")
	if len(vs) <= 2 {
		return v
	}
	return fmt.Sprintf("%s.%s", vs[0], vs[1])
}
