package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"oldnicke/etcd/clientv3"
	"oldnicke/etcd/pkg/testutil"
)

func TestCtlV3Migrate(t *testing.T) {
	defer testutil.AfterTest(t)

	epc := setupEtcdctlTest(t, &configNoTLS, false)
	defer func() {
		if errC := epc.Close(); errC != nil {
			t.Fatalf("error closing etcd processes (%v)", errC)
		}
	}()

	keys := make([]string, 3)
	vals := make([]string, 3)
	for i := range keys {
		keys[i] = fmt.Sprintf("foo_%d", i)
		vals[i] = fmt.Sprintf("bar_%d", i)
	}
	os.Setenv("ETCDCTL_API", "2")
	for i := range keys {
		if err := etcdctlSet(epc, keys[i], vals[i]); err != nil {
			t.Fatal(err)
		}
	}

	dataDir := epc.procs[0].Config().dataDirPath
	if err := epc.Stop(); err != nil {
		t.Fatalf("error closing etcd processes (%v)", err)
	}

	os.Unsetenv("ETCDCTL_API")
	cx := ctlCtx{
		t:           t,
		cfg:         configNoTLS,
		dialTimeout: 7 * time.Second,
		epc:         epc,
	}
	if err := ctlV3Migrate(cx, dataDir, ""); err != nil {
		t.Fatal(err)
	}

	epc.procs[0].Config().keepDataDir = true
	if err := epc.Restart(); err != nil {
		t.Fatal(err)
	}

	// to ensure revision increment is continuous from migrated v2 data
	if err := ctlV3Put(cx, "test", "value", ""); err != nil {
		t.Fatal(err)
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   epc.EndpointsV3(),
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	resp, err := cli.Get(context.TODO(), "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Kvs) != 1 {
		t.Fatalf("len(resp.Kvs) expected 1, got %+v", resp.Kvs)
	}
	if resp.Kvs[0].CreateRevision != 7 {
		t.Fatalf("resp.Kvs[0].CreateRevision expected 7, got %d", resp.Kvs[0].CreateRevision)
	}
}

func ctlV3Migrate(cx ctlCtx, dataDir, walDir string) error {
	cmdArgs := append(cx.PrefixArgs(), "migrate", "--data-dir", dataDir, "--wal-dir", walDir)
	return spawnWithExpects(cmdArgs, "finished transforming keys")
}
