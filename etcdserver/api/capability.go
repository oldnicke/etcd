package api

import (
	"sync"

	"go.uber.org/zap"
	"oldnicke/etcd/version"

	"github.com/coreos/go-semver/semver"
	"github.com/coreos/pkg/capnslog"
)

type Capability string

const (
	AuthCapability  Capability = "auth"
	V3rpcCapability Capability = "v3rpc"
)

var (
	plog = capnslog.NewPackageLogger("oldnicke/etcd", "etcdserver/api")

	// capabilityMaps is a static map of version to capability map.
	capabilityMaps = map[string]map[Capability]bool{
		"3.0.0": {AuthCapability: true, V3rpcCapability: true},
		"3.1.0": {AuthCapability: true, V3rpcCapability: true},
		"3.2.0": {AuthCapability: true, V3rpcCapability: true},
		"3.3.0": {AuthCapability: true, V3rpcCapability: true},
		"3.4.0": {AuthCapability: true, V3rpcCapability: true},
	}

	enableMapMu sync.RWMutex
	// enabledMap points to a map in capabilityMaps
	enabledMap map[Capability]bool

	curVersion *semver.Version
)

func init() {
	enabledMap = map[Capability]bool{
		AuthCapability:  true,
		V3rpcCapability: true,
	}
}

// UpdateCapability updates the enabledMap when the cluster version increases.
func UpdateCapability(lg *zap.Logger, v *semver.Version) {
	if v == nil {
		// if recovered but version was never set by cluster
		return
	}
	enableMapMu.Lock()
	if curVersion != nil && !curVersion.LessThan(*v) {
		enableMapMu.Unlock()
		return
	}
	curVersion = v
	enabledMap = capabilityMaps[curVersion.String()]
	enableMapMu.Unlock()

	if lg != nil {
		lg.Info(
			"enabled capabilities for version",
			zap.String("cluster-version", version.Cluster(v.String())),
		)
	} else {
		plog.Infof("enabled capabilities for version %s", version.Cluster(v.String()))
	}
}

func IsCapabilityEnabled(c Capability) bool {
	enableMapMu.RLock()
	defer enableMapMu.RUnlock()
	if enabledMap == nil {
		return false
	}
	return enabledMap[c]
}

func EnableCapability(c Capability) {
	enableMapMu.Lock()
	defer enableMapMu.Unlock()
	enabledMap[c] = true
}
