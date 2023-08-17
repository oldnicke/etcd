package cmd

import (
	"sync"
	"time"

	"oldnicke/etcd/pkg/transport"

	"github.com/spf13/cobra"
	"gopkg.in/cheggaaa/pb.v1"
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "A low-level benchmark tool for etcd3",
	Long: `benchmark is a low-level benchmark tool for etcd3.
It uses gRPC client directly and does not depend on
etcd client library.
	`,
}

var (
	endpoints    []string
	totalConns   uint
	totalClients uint
	precise      bool
	sample       bool

	bar *pb.ProgressBar
	wg  sync.WaitGroup

	tls transport.TLSInfo

	cpuProfPath string
	memProfPath string

	user string

	dialTimeout time.Duration

	targetLeader bool
)

func init() {
	RootCmd.PersistentFlags().StringSliceVar(&endpoints, "endpoints", []string{"127.0.0.1:2379"}, "gRPC endpoints")
	RootCmd.PersistentFlags().UintVar(&totalConns, "conns", 1, "Total number of gRPC connections")
	RootCmd.PersistentFlags().UintVar(&totalClients, "clients", 1, "Total number of gRPC clients")

	RootCmd.PersistentFlags().BoolVar(&precise, "precise", false, "use full floating point precision")
	RootCmd.PersistentFlags().BoolVar(&sample, "sample", false, "'true' to sample requests for every second")
	RootCmd.PersistentFlags().StringVar(&tls.CertFile, "cert", "", "identify HTTPS client using this SSL certificate file")
	RootCmd.PersistentFlags().StringVar(&tls.KeyFile, "key", "", "identify HTTPS client using this SSL key file")
	RootCmd.PersistentFlags().StringVar(&tls.TrustedCAFile, "cacert", "", "verify certificates of HTTPS-enabled servers using this CA bundle")

	RootCmd.PersistentFlags().StringVar(&user, "user", "", "provide username[:password] and prompt if password is not supplied.")
	RootCmd.PersistentFlags().DurationVar(&dialTimeout, "dial-timeout", 0, "dial timeout for client connections")

	RootCmd.PersistentFlags().BoolVar(&targetLeader, "target-leader", false, "connect only to the leader node")
}
