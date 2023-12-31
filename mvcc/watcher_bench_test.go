package mvcc

import (
	"fmt"
	"testing"

	"github.com/oldnicke/etcd/lease"
	"github.com/oldnicke/etcd/mvcc/backend"

	"go.uber.org/zap"
)

func BenchmarkKVWatcherMemoryUsage(b *testing.B) {
	be, tmpPath := backend.NewDefaultTmpBackend()
	watchable := newWatchableStore(zap.NewExample(), be, &lease.FakeLessor{}, nil, nil, StoreConfig{})

	defer cleanup(watchable, be, tmpPath)

	w := watchable.NewWatchStream()

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w.Watch(0, []byte(fmt.Sprint("foo", i)), nil, 0)
	}
}
