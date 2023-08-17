// Package main is a simple wrapper of the real etcd entrypoint package
// (located at github.com/oldnicke/etcd/etcdmain) to ensure that etcd is still
// "go getable"; e.g. `go get github.com/oldnicke/etcd` works as expected and
// builds a binary in $GOBIN/etcd
//
// This package should NOT be extended or modified in any way; to modify the
// etcd binary, work in the `github.com/oldnicke/etcd/etcdmain` package.
package main

import "github.com/oldnicke/etcd/etcdmain"

func main() {
	etcdmain.Main()
}
