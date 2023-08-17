package grpcproxy

import (
	"context"
	"net"
	"testing"
	"time"

	"oldnicke/etcd/clientv3"
	pb "oldnicke/etcd/etcdserver/etcdserverpb"
	"oldnicke/etcd/integration"
	"oldnicke/etcd/pkg/testutil"

	"google.golang.org/grpc"
)

func TestKVProxyRange(t *testing.T) {
	defer testutil.AfterTest(t)

	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	kvts := newKVProxyServer([]string{clus.Members[0].GRPCAddr()}, t)
	defer kvts.close()

	// create a client and try to get key from proxy.
	cfg := clientv3.Config{
		Endpoints:   []string{kvts.l.Addr().String()},
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	_, err = client.Get(context.Background(), "foo")
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	client.Close()
}

type kvproxyTestServer struct {
	kp     pb.KVServer
	c      *clientv3.Client
	server *grpc.Server
	l      net.Listener
}

func (kts *kvproxyTestServer) close() {
	kts.server.Stop()
	kts.l.Close()
	kts.c.Close()
}

func newKVProxyServer(endpoints []string, t *testing.T) *kvproxyTestServer {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	kvp, _ := NewKvProxy(client)

	kvts := &kvproxyTestServer{
		kp: kvp,
		c:  client,
	}

	var opts []grpc.ServerOption
	kvts.server = grpc.NewServer(opts...)
	pb.RegisterKVServer(kvts.server, kvts.kp)

	kvts.l, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	go kvts.server.Serve(kvts.l)

	return kvts
}
