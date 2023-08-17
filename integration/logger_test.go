package integration

import (
	"io/ioutil"

	"github.com/oldnicke/etcd/clientv3"

	"google.golang.org/grpc/grpclog"
)

func init() {
	clientv3.SetLogger(grpclog.NewLoggerV2(ioutil.Discard, ioutil.Discard, ioutil.Discard))
}
