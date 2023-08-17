package mockstorage

import (
	"go.etcd.io/etcd/pkg/testutil"
	"go.etcd.io/etcd/raft"
	"go.etcd.io/etcd/raft/raftpb"
)

type storageRecorder struct {
	testutil.Recorder
	dbPath string // must have '/' suffix if set
}

func NewStorageRecorder(db string) *storageRecorder {
	return &storageRecorder{&testutil.RecorderBuffered{}, db}
}

func NewStorageRecorderStream(db string) *storageRecorder {
	return &storageRecorder{testutil.NewRecorderStream(), db}
}

func (p *storageRecorder) Save(st raftpb.HardState, ents []raftpb.Entry) error {
	p.Record(testutil.Action{Name: "Save"})
	return nil
}

func (p *storageRecorder) SaveSnap(st raftpb.Snapshot) error {
	if !raft.IsEmptySnap(st) {
		p.Record(testutil.Action{Name: "SaveSnap"})
	}
	return nil
}

func (p *storageRecorder) Release(st raftpb.Snapshot) error {
	if !raft.IsEmptySnap(st) {
		p.Record(testutil.Action{Name: "Release"})
	}
	return nil
}

func (p *storageRecorder) Sync() error {
	p.Record(testutil.Action{Name: "Sync"})
	return nil
}

func (p *storageRecorder) Close() error { return nil }
