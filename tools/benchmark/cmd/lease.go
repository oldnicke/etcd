package cmd

import (
	"context"
	"fmt"
	"time"

	v3 "oldnicke/etcd/clientv3"
	"oldnicke/etcd/pkg/report"

	"github.com/spf13/cobra"
	"gopkg.in/cheggaaa/pb.v1"
)

var leaseKeepaliveCmd = &cobra.Command{
	Use:   "lease-keepalive",
	Short: "Benchmark lease keepalive",

	Run: leaseKeepaliveFunc,
}

var (
	leaseKeepaliveTotal int
)

func init() {
	RootCmd.AddCommand(leaseKeepaliveCmd)
	leaseKeepaliveCmd.Flags().IntVar(&leaseKeepaliveTotal, "total", 10000, "Total number of lease keepalive requests")
}

func leaseKeepaliveFunc(cmd *cobra.Command, args []string) {
	requests := make(chan struct{})
	clients := mustCreateClients(totalClients, totalConns)

	bar = pb.New(leaseKeepaliveTotal)
	bar.Format("Bom !")
	bar.Start()

	r := newReport()
	for i := range clients {
		wg.Add(1)
		go func(c v3.Lease) {
			defer wg.Done()
			resp, err := c.Grant(context.Background(), 100)
			if err != nil {
				panic(err)
			}
			for range requests {
				st := time.Now()
				_, err := c.KeepAliveOnce(context.TODO(), resp.ID)
				r.Results() <- report.Result{Err: err, Start: st, End: time.Now()}
				bar.Increment()
			}
		}(clients[i])
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < leaseKeepaliveTotal; i++ {
			requests <- struct{}{}
		}
		close(requests)
	}()

	rc := r.Run()
	wg.Wait()
	close(r.Results())
	bar.Finish()
	fmt.Printf("%s", <-rc)
}
