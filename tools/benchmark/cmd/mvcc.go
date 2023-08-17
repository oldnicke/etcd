package cmd

import (
	"os"
	"time"

	"go.uber.org/zap"

	"go.etcd.io/etcd/lease"
	"go.etcd.io/etcd/mvcc"
	"go.etcd.io/etcd/mvcc/backend"

	"github.com/spf13/cobra"
)

var (
	batchInterval int
	batchLimit    int

	s mvcc.KV
)

func initMVCC() {
	bcfg := backend.DefaultBackendConfig()
	bcfg.Path, bcfg.BatchInterval, bcfg.BatchLimit = "mvcc-bench", time.Duration(batchInterval)*time.Millisecond, batchLimit
	be := backend.New(bcfg)
	s = mvcc.NewStore(zap.NewExample(), be, &lease.FakeLessor{}, nil, mvcc.StoreConfig{})
	os.Remove("mvcc-bench") // boltDB has an opened fd, so removing the file is ok
}

// mvccCmd represents the MVCC storage benchmarking tools
var mvccCmd = &cobra.Command{
	Use:   "mvcc",
	Short: "Benchmark mvcc",
	Long: `storage subcommand is a set of various benchmark tools for MVCC storage subsystem of etcd.
Actual benchmarks are implemented as its subcommands.`,

	PersistentPreRun: mvccPreRun,
}

func init() {
	RootCmd.AddCommand(mvccCmd)

	mvccCmd.PersistentFlags().IntVar(&batchInterval, "batch-interval", 100, "Interval of batching (milliseconds)")
	mvccCmd.PersistentFlags().IntVar(&batchLimit, "batch-limit", 10000, "A limit of batched transaction")
}

func mvccPreRun(cmd *cobra.Command, args []string) {
	initMVCC()
}
