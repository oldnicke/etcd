package e2e

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"oldnicke/etcd/clientv3"
)

func TestCtlV3EndpointHealth(t *testing.T) { testCtl(t, endpointHealthTest, withQuorum()) }
func TestCtlV3EndpointStatus(t *testing.T) { testCtl(t, endpointStatusTest, withQuorum()) }
func TestCtlV3EndpointHashKV(t *testing.T) { testCtl(t, endpointHashKVTest, withQuorum()) }

func endpointHealthTest(cx ctlCtx) {
	if err := ctlV3EndpointHealth(cx); err != nil {
		cx.t.Fatalf("endpointStatusTest ctlV3EndpointHealth error (%v)", err)
	}
}

func ctlV3EndpointHealth(cx ctlCtx) error {
	cmdArgs := append(cx.PrefixArgs(), "endpoint", "health")
	lines := make([]string, cx.epc.cfg.clusterSize)
	for i := range lines {
		lines[i] = "is healthy"
	}
	return spawnWithExpects(cmdArgs, lines...)
}

func endpointStatusTest(cx ctlCtx) {
	if err := ctlV3EndpointStatus(cx); err != nil {
		cx.t.Fatalf("endpointStatusTest ctlV3EndpointStatus error (%v)", err)
	}
}

func ctlV3EndpointStatus(cx ctlCtx) error {
	cmdArgs := append(cx.PrefixArgs(), "endpoint", "status")
	var eps []string
	for _, ep := range cx.epc.EndpointsV3() {
		u, _ := url.Parse(ep)
		eps = append(eps, u.Host)
	}
	return spawnWithExpects(cmdArgs, eps...)
}

func endpointHashKVTest(cx ctlCtx) {
	if err := ctlV3EndpointHashKV(cx); err != nil {
		cx.t.Fatalf("endpointHashKVTest ctlV3EndpointHashKV error (%v)", err)
	}
}

func ctlV3EndpointHashKV(cx ctlCtx) error {
	eps := cx.epc.EndpointsV3()

	// get latest hash to compare
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   eps,
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		cx.t.Fatal(err)
	}
	defer cli.Close()
	hresp, err := cli.HashKV(context.TODO(), eps[0], 0)
	if err != nil {
		cx.t.Fatal(err)
	}

	cmdArgs := append(cx.PrefixArgs(), "endpoint", "hashkv")
	var ss []string
	for _, ep := range cx.epc.EndpointsV3() {
		u, _ := url.Parse(ep)
		ss = append(ss, fmt.Sprintf("%s, %d", u.Host, hresp.Hash))
	}
	return spawnWithExpects(cmdArgs, ss...)
}
