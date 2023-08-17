package command

import (
	"fmt"
	"os"

	v3 "oldnicke/etcd/clientv3"
	pb "oldnicke/etcd/etcdserver/etcdserverpb"
	mvccpb "oldnicke/etcd/mvcc/mvccpb"
)

type pbPrinter struct{ printer }

type pbMarshal interface {
	Marshal() ([]byte, error)
}

func newPBPrinter() printer {
	return &pbPrinter{
		&printerRPC{newPrinterUnsupported("protobuf"), printPB},
	}
}

func (p *pbPrinter) Watch(r v3.WatchResponse) {
	evs := make([]*mvccpb.Event, len(r.Events))
	for i, ev := range r.Events {
		evs[i] = (*mvccpb.Event)(ev)
	}
	wr := pb.WatchResponse{
		Header:          &r.Header,
		Events:          evs,
		CompactRevision: r.CompactRevision,
		Canceled:        r.Canceled,
		Created:         r.Created,
	}
	printPB(&wr)
}

func printPB(v interface{}) {
	m, ok := v.(pbMarshal)
	if !ok {
		ExitWithError(ExitBadFeature, fmt.Errorf("marshal unsupported for type %T (%v)", v, v))
	}
	b, err := m.Marshal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	fmt.Print(string(b))
}
