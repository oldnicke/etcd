package cmd

import (
	"context"
	"fmt"
	"sync"
	"time"

	v3 "oldnicke/etcd/clientv3"
	"oldnicke/etcd/pkg/report"

	"github.com/spf13/cobra"
	"gopkg.in/cheggaaa/pb.v1"
)

// watchGetCmd represents the watch command
var watchGetCmd = &cobra.Command{
	Use:   "watch-get",
	Short: "Benchmark watch with get",
	Long:  `Benchmark for serialized key gets with many unsynced watchers`,
	Run:   watchGetFunc,
}

var (
	watchGetTotalWatchers int
	watchGetTotalStreams  int
	watchEvents           int
	firstWatch            sync.Once
)

func init() {
	RootCmd.AddCommand(watchGetCmd)
	watchGetCmd.Flags().IntVar(&watchGetTotalWatchers, "watchers", 10000, "Total number of watchers")
	watchGetCmd.Flags().IntVar(&watchGetTotalStreams, "streams", 1, "Total number of watcher streams")
	watchGetCmd.Flags().IntVar(&watchEvents, "events", 8, "Number of events per watcher")
}

func watchGetFunc(cmd *cobra.Command, args []string) {
	clients := mustCreateClients(totalClients, totalConns)
	getClient := mustCreateClients(1, 1)

	// setup keys for watchers
	watchRev := int64(0)
	for i := 0; i < watchEvents; i++ {
		v := fmt.Sprintf("%d", i)
		resp, err := clients[0].Put(context.TODO(), "watchkey", v)
		if err != nil {
			panic(err)
		}
		if i == 0 {
			watchRev = resp.Header.Revision
		}
	}

	streams := make([]v3.Watcher, watchGetTotalStreams)
	for i := range streams {
		streams[i] = v3.NewWatcher(clients[i%len(clients)])
	}

	bar = pb.New(watchGetTotalWatchers * watchEvents)
	bar.Format("Bom !")
	bar.Start()

	// report from trying to do serialized gets with concurrent watchers
	r := newReport()
	ctx, cancel := context.WithCancel(context.TODO())
	f := func() {
		defer close(r.Results())
		for {
			st := time.Now()
			_, err := getClient[0].Get(ctx, "abc", v3.WithSerializable())
			if ctx.Err() != nil {
				break
			}
			r.Results() <- report.Result{Err: err, Start: st, End: time.Now()}
		}
	}

	wg.Add(watchGetTotalWatchers)
	for i := 0; i < watchGetTotalWatchers; i++ {
		go doUnsyncWatch(streams[i%len(streams)], watchRev, f)
	}

	rc := r.Run()
	wg.Wait()
	cancel()
	bar.Finish()
	fmt.Printf("Get during watch summary:\n%s", <-rc)
}

func doUnsyncWatch(stream v3.Watcher, rev int64, f func()) {
	defer wg.Done()
	wch := stream.Watch(context.TODO(), "watchkey", v3.WithRev(rev))
	if wch == nil {
		panic("could not open watch channel")
	}
	firstWatch.Do(func() { go f() })
	i := 0
	for i < watchEvents {
		wev := <-wch
		i += len(wev.Events)
		bar.Add(len(wev.Events))
	}
}
