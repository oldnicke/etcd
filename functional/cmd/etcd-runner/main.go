// etcd-runner is a program for testing etcd clientv3 features
// against a fault injected cluster.
package main

import "oldnicke/etcd/functional/runner"

func main() {
	runner.Start()
}
