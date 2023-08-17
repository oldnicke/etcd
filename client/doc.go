/*
Package client provides bindings for the etcd APIs.

Create a Config and exchange it for a Client:

	import (
		"net/http"
		"context"

		"github.com/oldnicke/etcd/client"
	)

	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: DefaultTransport,
	}

	c, err := client.New(cfg)
	if err != nil {
		// handle error
	}

Clients are safe for concurrent use by multiple goroutines.

Create a KeysAPI using the Client, then use it to interact with etcd:

	kAPI := client.NewKeysAPI(c)

	// create a new key /foo with the value "bar"
	_, err = kAPI.Create(context.Background(), "/foo", "bar")
	if err != nil {
		// handle error
	}

	// delete the newly created key only if the value is still "bar"
	_, err = kAPI.Delete(context.Background(), "/foo", &DeleteOptions{PrevValue: "bar"})
	if err != nil {
		// handle error
	}

Use a custom context to set timeouts on your operations:

	import "time"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// set a new key, ignoring its previous state
	_, err := kAPI.Set(ctx, "/ping", "pong", nil)
	if err != nil {
		if err == context.DeadlineExceeded {
			// request took longer than 5s
		} else {
			// handle error
		}
	}
*/
package client
