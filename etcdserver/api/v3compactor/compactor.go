package v3compactor

import (
	"context"
	"fmt"
	"time"

	pb "oldnicke/etcd/etcdserver/etcdserverpb"

	"github.com/coreos/pkg/capnslog"
	"github.com/jonboulle/clockwork"
	"go.uber.org/zap"
)

var (
	plog = capnslog.NewPackageLogger("oldnicke/etcd", "compactor")
)

const (
	ModePeriodic = "periodic"
	ModeRevision = "revision"
)

// Compactor purges old log from the storage periodically.
type Compactor interface {
	// Run starts the main loop of the compactor in background.
	// Use Stop() to halt the loop and release the resource.
	Run()
	// Stop halts the main loop of the compactor.
	Stop()
	// Pause temporally suspend the compactor not to run compaction. Resume() to unpose.
	Pause()
	// Resume restarts the compactor suspended by Pause().
	Resume()
}

type Compactable interface {
	Compact(ctx context.Context, r *pb.CompactionRequest) (*pb.CompactionResponse, error)
}

type RevGetter interface {
	Rev() int64
}

// New returns a new Compactor based on given "mode".
func New(
	lg *zap.Logger,
	mode string,
	retention time.Duration,
	rg RevGetter,
	c Compactable,
) (Compactor, error) {
	switch mode {
	case ModePeriodic:
		return newPeriodic(lg, clockwork.NewRealClock(), retention, rg, c), nil
	case ModeRevision:
		return newRevision(lg, clockwork.NewRealClock(), int64(retention), rg, c), nil
	default:
		return nil, fmt.Errorf("unsupported compaction mode %s", mode)
	}
}
