package etcdserver

import (
	"io"

	"github.com/oldnicke/etcd/etcdserver/api/snap"
	pb "github.com/oldnicke/etcd/etcdserver/etcdserverpb"
	"github.com/oldnicke/etcd/pkg/pbutil"
	"github.com/oldnicke/etcd/pkg/types"
	"github.com/oldnicke/etcd/raft/raftpb"
	"github.com/oldnicke/etcd/wal"
	"github.com/oldnicke/etcd/wal/walpb"

	"go.uber.org/zap"
)

type Storage interface {
	// Save function saves ents and state to the underlying stable storage.
	// Save MUST block until st and ents are on stable storage.
	Save(st raftpb.HardState, ents []raftpb.Entry) error
	// SaveSnap function saves snapshot to the underlying stable storage.
	SaveSnap(snap raftpb.Snapshot) error
	// Close closes the Storage and performs finalization.
	Close() error
	// Release releases the locked wal files older than the provided snapshot.
	Release(snap raftpb.Snapshot) error
	// Sync WAL
	Sync() error
}

type storage struct {
	*wal.WAL
	*snap.Snapshotter
}

func NewStorage(w *wal.WAL, s *snap.Snapshotter) Storage {
	return &storage{w, s}
}

// SaveSnap saves the snapshot file to disk and writes the WAL snapshot entry.
func (st *storage) SaveSnap(snap raftpb.Snapshot) error {
	walsnap := walpb.Snapshot{
		Index: snap.Metadata.Index,
		Term:  snap.Metadata.Term,
	}
	// save the snapshot file before writing the snapshot to the wal.
	// This makes it possible for the snapshot file to become orphaned, but prevents
	// a WAL snapshot entry from having no corresponding snapshot file.
	err := st.Snapshotter.SaveSnap(snap)
	if err != nil {
		return err
	}
	// gofail: var raftBeforeWALSaveSnaphot struct{}

	return st.WAL.SaveSnapshot(walsnap)
}

// Release releases resources older than the given snap and are no longer needed:
// - releases the locks to the wal files that are older than the provided wal for the given snap.
// - deletes any .snap.db files that are older than the given snap.
func (st *storage) Release(snap raftpb.Snapshot) error {
	if err := st.WAL.ReleaseLockTo(snap.Metadata.Index); err != nil {
		return err
	}
	return st.Snapshotter.ReleaseSnapDBs(snap)
}

// readWAL reads the WAL at the given snap and returns the wal, its latest HardState and cluster ID, and all entries that appear
// after the position of the given snap in the WAL.
// The snap must have been previously saved to the WAL, or this call will panic.
func readWAL(lg *zap.Logger, waldir string, snap walpb.Snapshot, unsafeNoFsync bool) (w *wal.WAL, id, cid types.ID, st raftpb.HardState, ents []raftpb.Entry) {
	var (
		err       error
		wmetadata []byte
	)

	repaired := false
	for {
		if w, err = wal.Open(lg, waldir, snap); err != nil {
			if lg != nil {
				lg.Fatal("failed to open WAL", zap.Error(err))
			} else {
				plog.Fatalf("open wal error: %v", err)
			}
		}
		if unsafeNoFsync {
			w.SetUnsafeNoFsync()
		}
		if wmetadata, st, ents, err = w.ReadAll(); err != nil {
			w.Close()
			// we can only repair ErrUnexpectedEOF and we never repair twice.
			if repaired || err != io.ErrUnexpectedEOF {
				if lg != nil {
					lg.Fatal("failed to read WAL, cannot be repaired", zap.Error(err))
				} else {
					plog.Fatalf("read wal error (%v) and cannot be repaired", err)
				}
			}
			if !wal.Repair(lg, waldir) {
				if lg != nil {
					lg.Fatal("failed to repair WAL", zap.Error(err))
				} else {
					plog.Fatalf("WAL error (%v) cannot be repaired", err)
				}
			} else {
				if lg != nil {
					lg.Info("repaired WAL", zap.Error(err))
				} else {
					plog.Infof("repaired WAL error (%v)", err)
				}
				repaired = true
			}
			continue
		}
		break
	}
	var metadata pb.Metadata
	pbutil.MustUnmarshal(&metadata, wmetadata)
	id = types.ID(metadata.NodeID)
	cid = types.ID(metadata.ClusterID)
	return w, id, cid, st, ents
}
