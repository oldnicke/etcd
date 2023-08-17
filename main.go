// Package main is a simple wrapper of the real etcd entrypoint package
// (located at oldnicke/etcd/etcdmain) to ensure that etcd is still
// "go getable"; e.g. `go get oldnicke/etcd` works as expected and
// builds a binary in $GOBIN/etcd
//
// This package should NOT be extended or modified in any way; to modify the
// etcd binary, work in the `oldnicke/etcd/etcdmain` package.
package main

import "oldnicke/etcd/etcdmain"

func main() {
	etcdmain.Main()
}
