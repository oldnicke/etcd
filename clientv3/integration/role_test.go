package integration

import (
	"context"
	"testing"

	"oldnicke/etcd/etcdserver/api/v3rpc/rpctypes"
	"oldnicke/etcd/integration"
	"oldnicke/etcd/pkg/testutil"
)

func TestRoleError(t *testing.T) {
	defer testutil.AfterTest(t)

	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	authapi := clus.RandClient()

	_, err := authapi.RoleAdd(context.TODO(), "test-role")
	if err != nil {
		t.Fatal(err)
	}

	_, err = authapi.RoleAdd(context.TODO(), "test-role")
	if err != rpctypes.ErrRoleAlreadyExist {
		t.Fatalf("expected %v, got %v", rpctypes.ErrRoleAlreadyExist, err)
	}

	_, err = authapi.RoleAdd(context.TODO(), "")
	if err != rpctypes.ErrRoleEmpty {
		t.Fatalf("expected %v, got %v", rpctypes.ErrRoleEmpty, err)
	}
}
