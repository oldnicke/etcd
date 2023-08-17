package etcdserver

import (
	"reflect"
	"testing"

	"go.etcd.io/etcd/pkg/types"
	"go.etcd.io/etcd/version"

	"github.com/coreos/go-semver/semver"
	"go.uber.org/zap"
)

var testLogger = zap.NewExample()

func TestDecideClusterVersion(t *testing.T) {
	tests := []struct {
		vers  map[string]*version.Versions
		wdver *semver.Version
	}{
		{
			map[string]*version.Versions{"a": {Server: "2.0.0"}},
			semver.Must(semver.NewVersion("2.0.0")),
		},
		// unknown
		{
			map[string]*version.Versions{"a": nil},
			nil,
		},
		{
			map[string]*version.Versions{"a": {Server: "2.0.0"}, "b": {Server: "2.1.0"}, "c": {Server: "2.1.0"}},
			semver.Must(semver.NewVersion("2.0.0")),
		},
		{
			map[string]*version.Versions{"a": {Server: "2.1.0"}, "b": {Server: "2.1.0"}, "c": {Server: "2.1.0"}},
			semver.Must(semver.NewVersion("2.1.0")),
		},
		{
			map[string]*version.Versions{"a": nil, "b": {Server: "2.1.0"}, "c": {Server: "2.1.0"}},
			nil,
		},
	}

	for i, tt := range tests {
		dver := decideClusterVersion(testLogger, tt.vers)
		if !reflect.DeepEqual(dver, tt.wdver) {
			t.Errorf("#%d: ver = %+v, want %+v", i, dver, tt.wdver)
		}
	}
}

func TestIsCompatibleWithVers(t *testing.T) {
	tests := []struct {
		vers       map[string]*version.Versions
		local      types.ID
		minV, maxV *semver.Version
		wok        bool
	}{
		// too low
		{
			map[string]*version.Versions{
				"a": {Server: "2.0.0", Cluster: "not_decided"},
				"b": {Server: "2.1.0", Cluster: "2.1.0"},
				"c": {Server: "2.1.0", Cluster: "2.1.0"},
			},
			0xa,
			semver.Must(semver.NewVersion("2.0.0")), semver.Must(semver.NewVersion("2.0.0")),
			false,
		},
		{
			map[string]*version.Versions{
				"a": {Server: "2.1.0", Cluster: "not_decided"},
				"b": {Server: "2.1.0", Cluster: "2.1.0"},
				"c": {Server: "2.1.0", Cluster: "2.1.0"},
			},
			0xa,
			semver.Must(semver.NewVersion("2.0.0")), semver.Must(semver.NewVersion("2.1.0")),
			true,
		},
		// too high
		{
			map[string]*version.Versions{
				"a": {Server: "2.2.0", Cluster: "not_decided"},
				"b": {Server: "2.0.0", Cluster: "2.0.0"},
				"c": {Server: "2.0.0", Cluster: "2.0.0"},
			},
			0xa,
			semver.Must(semver.NewVersion("2.1.0")), semver.Must(semver.NewVersion("2.2.0")),
			false,
		},
		// cannot get b's version, expect ok
		{
			map[string]*version.Versions{
				"a": {Server: "2.1.0", Cluster: "not_decided"},
				"b": nil,
				"c": {Server: "2.1.0", Cluster: "2.1.0"},
			},
			0xa,
			semver.Must(semver.NewVersion("2.0.0")), semver.Must(semver.NewVersion("2.1.0")),
			true,
		},
		// cannot get b and c's version, expect not ok
		{
			map[string]*version.Versions{
				"a": {Server: "2.1.0", Cluster: "not_decided"},
				"b": nil,
				"c": nil,
			},
			0xa,
			semver.Must(semver.NewVersion("2.0.0")), semver.Must(semver.NewVersion("2.1.0")),
			false,
		},
	}

	for i, tt := range tests {
		ok := isCompatibleWithVers(testLogger, tt.vers, tt.local, tt.minV, tt.maxV)
		if ok != tt.wok {
			t.Errorf("#%d: ok = %+v, want %+v", i, ok, tt.wok)
		}
	}
}
