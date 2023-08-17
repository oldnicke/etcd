package clientv3

import (
	"reflect"
	"testing"

	"oldnicke/etcd/etcdserver/etcdserverpb"
)

func TestCompactOp(t *testing.T) {
	req1 := OpCompact(100, WithCompactPhysical()).toRequest()
	req2 := &etcdserverpb.CompactionRequest{Revision: 100, Physical: true}
	if !reflect.DeepEqual(req1, req2) {
		t.Fatalf("expected %+v, got %+v", req2, req1)
	}
}
