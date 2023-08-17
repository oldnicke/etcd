package cmd

import (
	"crypto/rand"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/oldnicke/etcd/lease"
	"github.com/oldnicke/etcd/pkg/report"
	"github.com/oldnicke/etcd/pkg/traceutil"

	"github.com/spf13/cobra"
)

// mvccPutCmd represents a storage put performance benchmarking tool
var mvccPutCmd = &cobra.Command{
	Use:   "put",
	Short: "Benchmark put performance of storage",

	Run: mvccPutFunc,
}

var (
	mvccTotalRequests int
	storageKeySize    int
	valueSize         int
	txn               bool
	nrTxnOps          int
)

func init() {
	mvccCmd.AddCommand(mvccPutCmd)

	mvccPutCmd.Flags().IntVar(&mvccTotalRequests, "total", 100, "a total number of keys to put")
	mvccPutCmd.Flags().IntVar(&storageKeySize, "key-size", 64, "a size of key (Byte)")
	mvccPutCmd.Flags().IntVar(&valueSize, "value-size", 64, "a size of value (Byte)")
	mvccPutCmd.Flags().BoolVar(&txn, "txn", false, "put a key in transaction or not")
	mvccPutCmd.Flags().IntVar(&nrTxnOps, "txn-ops", 1, "a number of keys to put per transaction")

	// TODO: after the PR https://github.com/spf13/cobra/pull/220 is merged, the below pprof related flags should be moved to RootCmd
	mvccPutCmd.Flags().StringVar(&cpuProfPath, "cpuprofile", "", "the path of file for storing cpu profile result")
	mvccPutCmd.Flags().StringVar(&memProfPath, "memprofile", "", "the path of file for storing heap profile result")

}

func createBytesSlice(bytesN, sliceN int) [][]byte {
	rs := make([][]byte, sliceN)
	for i := range rs {
		rs[i] = make([]byte, bytesN)
		if _, err := rand.Read(rs[i]); err != nil {
			panic(err)
		}
	}
	return rs
}

func mvccPutFunc(cmd *cobra.Command, args []string) {
	if cpuProfPath != "" {
		f, err := os.Create(cpuProfPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to create a file for storing cpu profile result: ", err)
			os.Exit(1)
		}

		err = pprof.StartCPUProfile(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to start cpu profile: ", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}

	if memProfPath != "" {
		f, err := os.Create(memProfPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to create a file for storing heap profile result: ", err)
			os.Exit(1)
		}

		defer func() {
			err := pprof.WriteHeapProfile(f)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to write heap profile result: ", err)
				// can do nothing for handling the error
			}
		}()
	}

	keys := createBytesSlice(storageKeySize, mvccTotalRequests*nrTxnOps)
	vals := createBytesSlice(valueSize, mvccTotalRequests*nrTxnOps)

	weight := float64(nrTxnOps)
	r := newWeightedReport()
	rrc := r.Results()

	rc := r.Run()

	if txn {
		for i := 0; i < mvccTotalRequests; i++ {
			st := time.Now()

			tw := s.Write(traceutil.TODO())
			for j := i; j < i+nrTxnOps; j++ {
				tw.Put(keys[j], vals[j], lease.NoLease)
			}
			tw.End()

			rrc <- report.Result{Start: st, End: time.Now(), Weight: weight}
		}
	} else {
		for i := 0; i < mvccTotalRequests; i++ {
			st := time.Now()
			s.Put(keys[i], vals[i], lease.NoLease)
			rrc <- report.Result{Start: st, End: time.Now()}
		}
	}

	close(r.Results())
	fmt.Printf("%s", <-rc)
}
