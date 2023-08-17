package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	rootCommand = &cobra.Command{
		Use:   "etcd-dump-db",
		Short: "etcd-dump-db inspects etcd db files.",
	}
	listBucketCommand = &cobra.Command{
		Use:   "list-bucket [data dir or db file path]",
		Short: "bucket lists all buckets.",
		Run:   listBucketCommandFunc,
	}
	iterateBucketCommand = &cobra.Command{
		Use:   "iterate-bucket [data dir or db file path] [bucket name]",
		Short: "iterate-bucket lists key-value pairs in reverse order.",
		Run:   iterateBucketCommandFunc,
	}
	getHashCommand = &cobra.Command{
		Use:   "hash [data dir or db file path]",
		Short: "hash computes the hash of db file.",
		Run:   getHashCommandFunc,
	}
)

var flockTimeout time.Duration
var iterateBucketLimit uint64
var iterateBucketDecode bool

func init() {
	rootCommand.PersistentFlags().DurationVar(&flockTimeout, "timeout", 10*time.Second, "time to wait to obtain a file lock on db file, 0 to block indefinitely")
	iterateBucketCommand.PersistentFlags().Uint64Var(&iterateBucketLimit, "limit", 0, "max number of key-value pairs to iterate (0< to iterate all)")
	iterateBucketCommand.PersistentFlags().BoolVar(&iterateBucketDecode, "decode", false, "true to decode Protocol Buffer encoded data")

	rootCommand.AddCommand(listBucketCommand)
	rootCommand.AddCommand(iterateBucketCommand)
	rootCommand.AddCommand(getHashCommand)
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func listBucketCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalf("Must provide at least 1 argument (got %v)", args)
	}
	dp := args[0]
	if !strings.HasSuffix(dp, "db") {
		dp = filepath.Join(snapDir(dp), "db")
	}
	if !existFileOrDir(dp) {
		log.Fatalf("%q does not exist", dp)
	}

	bts, err := getBuckets(dp)
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range bts {
		fmt.Println(b)
	}
}

func iterateBucketCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.Fatalf("Must provide 2 arguments (got %v)", args)
	}
	dp := args[0]
	if !strings.HasSuffix(dp, "db") {
		dp = filepath.Join(snapDir(dp), "db")
	}
	if !existFileOrDir(dp) {
		log.Fatalf("%q does not exist", dp)
	}
	bucket := args[1]
	err := iterateBucket(dp, bucket, iterateBucketLimit, iterateBucketDecode)
	if err != nil {
		log.Fatal(err)
	}
}

func getHashCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalf("Must provide at least 1 argument (got %v)", args)
	}
	dp := args[0]
	if !strings.HasSuffix(dp, "db") {
		dp = filepath.Join(snapDir(dp), "db")
	}
	if !existFileOrDir(dp) {
		log.Fatalf("%q does not exist", dp)
	}

	hash, err := getHash(dp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("db path: %s\nHash: %d\n", dp, hash)
}
