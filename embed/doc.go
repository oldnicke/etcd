/*
Package embed provides bindings for embedding an etcd server in a program.

Launch an embedded etcd server using the configuration defaults:

	import (
		"log"
		"time"

		"github.com/oldnicke/etcd/embed"
	)

	func main() {
		cfg := embed.NewConfig()
		cfg.Dir = "default.etcd"
		e, err := embed.StartEtcd(cfg)
		if err != nil {
			log.Fatal(err)
		}
		defer e.Close()
		select {
		case <-e.Server.ReadyNotify():
			log.Printf("Server is ready!")
		case <-time.After(60 * time.Second):
			e.Server.Stop() // trigger a shutdown
			log.Printf("Server took too long to start!")
		}
		log.Fatal(<-e.Err())
	}
*/
package embed
