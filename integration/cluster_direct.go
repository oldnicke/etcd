//go:build !cluster_proxy
// +build !cluster_proxy

package integration

import (
	"github.com/oldnicke/etcd/clientv3"
	"github.com/oldnicke/etcd/etcdserver/api/v3election/v3electionpb"
	"github.com/oldnicke/etcd/etcdserver/api/v3lock/v3lockpb"
	pb "github.com/oldnicke/etcd/etcdserver/etcdserverpb"
)

func toGRPC(c *clientv3.Client) grpcAPI {
	return grpcAPI{
		pb.NewClusterClient(c.ActiveConnection()),
		pb.NewKVClient(c.ActiveConnection()),
		pb.NewLeaseClient(c.ActiveConnection()),
		pb.NewWatchClient(c.ActiveConnection()),
		pb.NewMaintenanceClient(c.ActiveConnection()),
		pb.NewAuthClient(c.ActiveConnection()),
		v3lockpb.NewLockClient(c.ActiveConnection()),
		v3electionpb.NewElectionClient(c.ActiveConnection()),
	}
}

func newClientV3(cfg clientv3.Config) (*clientv3.Client, error) {
	return clientv3.New(cfg)
}
