package adapter

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	pb "oldnicke/etcd/etcdserver/etcdserverpb"
)

var errAlreadySentHeader = errors.New("adapter: already sent header")

type ws2wc struct{ wserv pb.WatchServer }

func WatchServerToWatchClient(wserv pb.WatchServer) pb.WatchClient {
	return &ws2wc{wserv}
}

func (s *ws2wc) Watch(ctx context.Context, opts ...grpc.CallOption) (pb.Watch_WatchClient, error) {
	cs := newPipeStream(ctx, func(ss chanServerStream) error {
		return s.wserv.Watch(&ws2wcServerStream{ss})
	})
	return &ws2wcClientStream{cs}, nil
}

// ws2wcClientStream implements Watch_WatchClient
type ws2wcClientStream struct{ chanClientStream }

// ws2wcServerStream implements Watch_WatchServer
type ws2wcServerStream struct{ chanServerStream }

func (s *ws2wcClientStream) Send(wr *pb.WatchRequest) error {
	return s.SendMsg(wr)
}
func (s *ws2wcClientStream) Recv() (*pb.WatchResponse, error) {
	var v interface{}
	if err := s.RecvMsg(&v); err != nil {
		return nil, err
	}
	return v.(*pb.WatchResponse), nil
}

func (s *ws2wcServerStream) Send(wr *pb.WatchResponse) error {
	return s.SendMsg(wr)
}
func (s *ws2wcServerStream) Recv() (*pb.WatchRequest, error) {
	var v interface{}
	if err := s.RecvMsg(&v); err != nil {
		return nil, err
	}
	return v.(*pb.WatchRequest), nil
}
