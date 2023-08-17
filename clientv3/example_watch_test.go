package clientv3_test

import (
	"context"
	"fmt"
	"log"

	"oldnicke/etcd/clientv3"
)

func ExampleWatcher_watch() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	rch := cli.Watch(context.Background(), "foo")
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
	// PUT "foo" : "bar"
}

func ExampleWatcher_watchWithPrefix() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	rch := cli.Watch(context.Background(), "foo", clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
	// PUT "foo1" : "bar"
}

func ExampleWatcher_watchWithRange() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// watches within ['foo1', 'foo4'), in lexicographical order
	rch := cli.Watch(context.Background(), "foo1", clientv3.WithRange("foo4"))
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
	// PUT "foo1" : "bar"
	// PUT "foo2" : "bar"
	// PUT "foo3" : "bar"
}

func ExampleWatcher_watchWithProgressNotify() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}

	rch := cli.Watch(context.Background(), "foo", clientv3.WithProgressNotify())
	wresp := <-rch
	fmt.Printf("wresp.Header.Revision: %d\n", wresp.Header.Revision)
	fmt.Println("wresp.IsProgressNotify:", wresp.IsProgressNotify())
	// wresp.Header.Revision: 0
	// wresp.IsProgressNotify: true
}
