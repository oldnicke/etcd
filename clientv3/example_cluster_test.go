package clientv3_test

import (
	"context"
	"fmt"
	"log"

	"github.com/oldnicke/etcd/clientv3"
)

func ExampleCluster_memberList() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	resp, err := cli.MemberList(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("members:", len(resp.Members))
	// Output: members: 3
}

func ExampleCluster_memberAdd() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints[:2],
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	peerURLs := endpoints[2:]
	mresp, err := cli.MemberAdd(context.Background(), peerURLs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("added member.PeerURLs:", mresp.Member.PeerURLs)
	// added member.PeerURLs: [http://localhost:32380]
}

func ExampleCluster_memberAddAsLearner() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints[:2],
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	peerURLs := endpoints[2:]
	mresp, err := cli.MemberAddAsLearner(context.Background(), peerURLs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("added member.PeerURLs:", mresp.Member.PeerURLs)
	fmt.Println("added member.IsLearner:", mresp.Member.IsLearner)
	// added member.PeerURLs: [http://localhost:32380]
	// added member.IsLearner: true
}

func ExampleCluster_memberRemove() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints[1:],
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	resp, err := cli.MemberList(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	_, err = cli.MemberRemove(context.Background(), resp.Members[0].ID)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleCluster_memberUpdate() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	resp, err := cli.MemberList(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	peerURLs := []string{"http://localhost:12380"}
	_, err = cli.MemberUpdate(context.Background(), resp.Members[0].ID, peerURLs)
	if err != nil {
		log.Fatal(err)
	}
}
