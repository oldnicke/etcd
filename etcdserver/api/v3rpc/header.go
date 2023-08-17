package v3rpc

import (
	"github.com/oldnicke/etcd/etcdserver"
	pb "github.com/oldnicke/etcd/etcdserver/etcdserverpb"
)

type header struct {
	clusterID int64
	memberID  int64
	sg        etcdserver.RaftStatusGetter
	rev       func() int64
}

func newHeader(s *etcdserver.EtcdServer) header {
	return header{
		clusterID: int64(s.Cluster().ID()),
		memberID:  int64(s.ID()),
		sg:        s,
		rev:       func() int64 { return s.KV().Rev() },
	}
}

// fill populates pb.ResponseHeader using etcdserver information
func (h *header) fill(rh *pb.ResponseHeader) {
	if rh == nil {
		plog.Panic("unexpected nil resp.Header")
	}
	rh.ClusterId = uint64(h.clusterID)
	rh.MemberId = uint64(h.memberID)
	rh.RaftTerm = h.sg.Term()
	if rh.Revision == 0 {
		rh.Revision = h.rev()
	}
}
