package integration

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/oldnicke/etcd/etcdserver/api/v3rpc/rpctypes"
	pb "github.com/oldnicke/etcd/etcdserver/etcdserverpb"
	"github.com/oldnicke/etcd/pkg/testutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestV3MaintenanceDefragmentInflightRange ensures inflight range requests
// does not panic the mvcc backend while defragment is running.
func TestV3MaintenanceDefragmentInflightRange(t *testing.T) {
	defer testutil.AfterTest(t)
	clus := NewClusterV3(t, &ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	cli := clus.RandClient()
	kvc := toGRPC(cli).KV
	if _, err := kvc.Put(context.Background(), &pb.PutRequest{Key: []byte("foo"), Value: []byte("bar")}); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		kvc.Range(ctx, &pb.RangeRequest{Key: []byte("foo")})
	}()

	mvc := toGRPC(cli).Maintenance
	mvc.Defragment(context.Background(), &pb.DefragmentRequest{})
	cancel()

	<-donec
}

// TestV3KVInflightRangeRequests ensures that inflight requests
// (sent before server shutdown) are gracefully handled by server-side.
// They are either finished or canceled, but never crash the backend.
// See https://github.com/etcd-io/etcd/issues/7322 for more detail.
func TestV3KVInflightRangeRequests(t *testing.T) {
	defer testutil.AfterTest(t)
	clus := NewClusterV3(t, &ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	cli := clus.RandClient()
	kvc := toGRPC(cli).KV

	if _, err := kvc.Put(context.Background(), &pb.PutRequest{Key: []byte("foo"), Value: []byte("bar")}); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	reqN := 10 // use 500+ for fast machine
	var wg sync.WaitGroup
	wg.Add(reqN)
	for i := 0; i < reqN; i++ {
		go func() {
			defer wg.Done()
			_, err := kvc.Range(ctx, &pb.RangeRequest{Key: []byte("foo"), Serializable: true}, grpc.FailFast(false))
			if err != nil {
				errCode := status.Convert(err).Code()
				errDesc := rpctypes.ErrorDesc(err)
				if err != nil && !(errDesc == context.Canceled.Error() || errCode == codes.Canceled || errCode == codes.Unavailable) {
					t.Errorf("inflight request should be canceled with '%v' or code Canceled or Unavailable, got '%v' with code '%s'", context.Canceled.Error(), errDesc, errCode)
				}
			}
		}()
	}

	clus.Members[0].Stop(t)
	cancel()

	wg.Wait()
}
