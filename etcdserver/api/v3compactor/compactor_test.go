package v3compactor

import (
	"context"
	"sync/atomic"

	pb "oldnicke/etcd/etcdserver/etcdserverpb"
	"oldnicke/etcd/pkg/testutil"
)

type fakeCompactable struct {
	testutil.Recorder
}

func (fc *fakeCompactable) Compact(ctx context.Context, r *pb.CompactionRequest) (*pb.CompactionResponse, error) {
	fc.Record(testutil.Action{Name: "c", Params: []interface{}{r}})
	return &pb.CompactionResponse{}, nil
}

type fakeRevGetter struct {
	testutil.Recorder
	rev int64
}

func (fr *fakeRevGetter) Rev() int64 {
	fr.Record(testutil.Action{Name: "g"})
	rev := atomic.AddInt64(&fr.rev, 1)
	return rev
}

func (fr *fakeRevGetter) SetRev(rev int64) {
	atomic.StoreInt64(&fr.rev, rev)
}
