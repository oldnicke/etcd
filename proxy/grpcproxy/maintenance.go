package grpcproxy

import (
	"context"
	"io"

	"github.com/oldnicke/etcd/clientv3"
	pb "github.com/oldnicke/etcd/etcdserver/etcdserverpb"
)

type maintenanceProxy struct {
	client *clientv3.Client
}

func NewMaintenanceProxy(c *clientv3.Client) pb.MaintenanceServer {
	return &maintenanceProxy{
		client: c,
	}
}

func (mp *maintenanceProxy) Defragment(ctx context.Context, dr *pb.DefragmentRequest) (*pb.DefragmentResponse, error) {
	conn := mp.client.ActiveConnection()
	return pb.NewMaintenanceClient(conn).Defragment(ctx, dr)
}

func (mp *maintenanceProxy) Snapshot(sr *pb.SnapshotRequest, stream pb.Maintenance_SnapshotServer) error {
	conn := mp.client.ActiveConnection()
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	ctx = withClientAuthToken(ctx, stream.Context())

	sc, err := pb.NewMaintenanceClient(conn).Snapshot(ctx, sr)
	if err != nil {
		return err
	}

	for {
		rr, err := sc.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		err = stream.Send(rr)
		if err != nil {
			return err
		}
	}
}

func (mp *maintenanceProxy) Hash(ctx context.Context, r *pb.HashRequest) (*pb.HashResponse, error) {
	conn := mp.client.ActiveConnection()
	return pb.NewMaintenanceClient(conn).Hash(ctx, r)
}

func (mp *maintenanceProxy) HashKV(ctx context.Context, r *pb.HashKVRequest) (*pb.HashKVResponse, error) {
	conn := mp.client.ActiveConnection()
	return pb.NewMaintenanceClient(conn).HashKV(ctx, r)
}

func (mp *maintenanceProxy) Alarm(ctx context.Context, r *pb.AlarmRequest) (*pb.AlarmResponse, error) {
	conn := mp.client.ActiveConnection()
	return pb.NewMaintenanceClient(conn).Alarm(ctx, r)
}

func (mp *maintenanceProxy) Status(ctx context.Context, r *pb.StatusRequest) (*pb.StatusResponse, error) {
	conn := mp.client.ActiveConnection()
	return pb.NewMaintenanceClient(conn).Status(ctx, r)
}

func (mp *maintenanceProxy) MoveLeader(ctx context.Context, r *pb.MoveLeaderRequest) (*pb.MoveLeaderResponse, error) {
	conn := mp.client.ActiveConnection()
	return pb.NewMaintenanceClient(conn).MoveLeader(ctx, r)
}
