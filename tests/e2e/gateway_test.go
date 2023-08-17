package e2e

import (
	"os"
	"strings"
	"testing"

	"oldnicke/etcd/pkg/expect"
)

var (
	defaultGatewayEndpoint = "127.0.0.1:23790"
)

func TestGateway(t *testing.T) {
	ec, err := newEtcdProcessCluster(&configNoTLS)
	if err != nil {
		t.Fatal(err)
	}
	defer ec.Stop()

	eps := strings.Join(ec.EndpointsV3(), ",")

	p := startGateway(t, eps)
	defer p.Stop()

	os.Setenv("ETCDCTL_API", "3")
	defer os.Unsetenv("ETCDCTL_API")

	err = spawnWithExpect([]string{ctlBinPath, "--endpoints=" + defaultGatewayEndpoint, "put", "foo", "bar"}, "OK\r\n")
	if err != nil {
		t.Errorf("failed to finish put request through gateway: %v", err)
	}
}

func startGateway(t *testing.T, endpoints string) *expect.ExpectProcess {
	p, err := expect.NewExpect(binPath, "gateway", "--endpoints="+endpoints, "start")
	if err != nil {
		t.Fatal(err)
	}
	_, err = p.Expect("ready to proxy client requests")
	if err != nil {
		t.Fatal(err)
	}
	return p
}
