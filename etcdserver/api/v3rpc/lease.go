package v3rpc

import (
	"context"
	"io"

	"github.com/oldnicke/etcd/etcdserver"
	"github.com/oldnicke/etcd/etcdserver/api/v3rpc/rpctypes"
	pb "github.com/oldnicke/etcd/etcdserver/etcdserverpb"
	"github.com/oldnicke/etcd/lease"

	"go.uber.org/zap"
)

type LeaseServer struct {
	lg  *zap.Logger
	hdr header
	le  etcdserver.Lessor
}

func NewLeaseServer(s *etcdserver.EtcdServer) pb.LeaseServer {
	return &LeaseServer{lg: s.Cfg.Logger, le: s, hdr: newHeader(s)}
}

func (ls *LeaseServer) LeaseGrant(ctx context.Context, cr *pb.LeaseGrantRequest) (*pb.LeaseGrantResponse, error) {
	resp, err := ls.le.LeaseGrant(ctx, cr)

	if err != nil {
		return nil, togRPCError(err)
	}
	ls.hdr.fill(resp.Header)
	return resp, nil
}

func (ls *LeaseServer) LeaseRevoke(ctx context.Context, rr *pb.LeaseRevokeRequest) (*pb.LeaseRevokeResponse, error) {
	resp, err := ls.le.LeaseRevoke(ctx, rr)
	if err != nil {
		return nil, togRPCError(err)
	}
	ls.hdr.fill(resp.Header)
	return resp, nil
}

func (ls *LeaseServer) LeaseTimeToLive(ctx context.Context, rr *pb.LeaseTimeToLiveRequest) (*pb.LeaseTimeToLiveResponse, error) {
	resp, err := ls.le.LeaseTimeToLive(ctx, rr)
	if err != nil && err != lease.ErrLeaseNotFound {
		return nil, togRPCError(err)
	}
	if err == lease.ErrLeaseNotFound {
		resp = &pb.LeaseTimeToLiveResponse{
			Header: &pb.ResponseHeader{},
			ID:     rr.ID,
			TTL:    -1,
		}
	}
	ls.hdr.fill(resp.Header)
	return resp, nil
}

func (ls *LeaseServer) LeaseLeases(ctx context.Context, rr *pb.LeaseLeasesRequest) (*pb.LeaseLeasesResponse, error) {
	resp, err := ls.le.LeaseLeases(ctx, rr)
	if err != nil && err != lease.ErrLeaseNotFound {
		return nil, togRPCError(err)
	}
	if err == lease.ErrLeaseNotFound {
		resp = &pb.LeaseLeasesResponse{
			Header: &pb.ResponseHeader{},
			Leases: []*pb.LeaseStatus{},
		}
	}
	ls.hdr.fill(resp.Header)
	return resp, nil
}

func (ls *LeaseServer) LeaseKeepAlive(stream pb.Lease_LeaseKeepAliveServer) (err error) {
	errc := make(chan error, 1)
	go func() {
		errc <- ls.leaseKeepAlive(stream)
	}()
	select {
	case err = <-errc:
	case <-stream.Context().Done():
		// the only server-side cancellation is noleader for now.
		err = stream.Context().Err()
		if err == context.Canceled {
			err = rpctypes.ErrGRPCNoLeader
		}
	}
	return err
}

func (ls *LeaseServer) leaseKeepAlive(stream pb.Lease_LeaseKeepAliveServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			if isClientCtxErr(stream.Context().Err(), err) {
				if ls.lg != nil {
					ls.lg.Debug("failed to receive lease keepalive request from gRPC stream", zap.Error(err))
				} else {
					plog.Debugf("failed to receive lease keepalive request from gRPC stream (%q)", err.Error())
				}
			} else {
				if ls.lg != nil {
					ls.lg.Warn("failed to receive lease keepalive request from gRPC stream", zap.Error(err))
				} else {
					plog.Warningf("failed to receive lease keepalive request from gRPC stream (%q)", err.Error())
				}
				streamFailures.WithLabelValues("receive", "lease-keepalive").Inc()
			}
			return err
		}

		// Create header before we sent out the renew request.
		// This can make sure that the revision is strictly smaller or equal to
		// when the keepalive happened at the local server (when the local server is the leader)
		// or remote leader.
		// Without this, a lease might be revoked at rev 3 but client can see the keepalive succeeded
		// at rev 4.
		resp := &pb.LeaseKeepAliveResponse{ID: req.ID, Header: &pb.ResponseHeader{}}
		ls.hdr.fill(resp.Header)

		ttl, err := ls.le.LeaseRenew(stream.Context(), lease.LeaseID(req.ID))
		if err == lease.ErrLeaseNotFound {
			err = nil
			ttl = 0
		}

		if err != nil {
			return togRPCError(err)
		}

		resp.TTL = ttl
		err = stream.Send(resp)
		if err != nil {
			if isClientCtxErr(stream.Context().Err(), err) {
				if ls.lg != nil {
					ls.lg.Debug("failed to send lease keepalive response to gRPC stream", zap.Error(err))
				} else {
					plog.Debugf("failed to send lease keepalive response to gRPC stream (%q)", err.Error())
				}
			} else {
				if ls.lg != nil {
					ls.lg.Warn("failed to send lease keepalive response to gRPC stream", zap.Error(err))
				} else {
					plog.Warningf("failed to send lease keepalive response to gRPC stream (%q)", err.Error())
				}
				streamFailures.WithLabelValues("send", "lease-keepalive").Inc()
			}
			return err
		}
	}
}
