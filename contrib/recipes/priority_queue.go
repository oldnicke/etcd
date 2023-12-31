package recipe

import (
	"context"
	"fmt"

	v3 "github.com/oldnicke/etcd/clientv3"
	"github.com/oldnicke/etcd/mvcc/mvccpb"
)

// PriorityQueue implements a multi-reader, multi-writer distributed queue.
type PriorityQueue struct {
	client *v3.Client
	ctx    context.Context
	key    string
}

// NewPriorityQueue creates an etcd priority queue.
func NewPriorityQueue(client *v3.Client, key string) *PriorityQueue {
	return &PriorityQueue{client, context.TODO(), key + "/"}
}

// Enqueue puts a value into a queue with a given priority.
func (q *PriorityQueue) Enqueue(val string, pr uint16) error {
	prefix := fmt.Sprintf("%s%05d", q.key, pr)
	_, err := newSequentialKV(q.client, prefix, val)
	return err
}

// Dequeue returns Enqueue()'d items in FIFO order. If the
// queue is empty, Dequeue blocks until items are available.
func (q *PriorityQueue) Dequeue() (string, error) {
	// TODO: fewer round trips by fetching more than one key
	resp, err := q.client.Get(q.ctx, q.key, v3.WithFirstKey()...)
	if err != nil {
		return "", err
	}

	kv, err := claimFirstKey(q.client, resp.Kvs)
	if err != nil {
		return "", err
	} else if kv != nil {
		return string(kv.Value), nil
	} else if resp.More {
		// missed some items, retry to read in more
		return q.Dequeue()
	}

	// nothing to dequeue; wait on items
	ev, err := WaitPrefixEvents(
		q.client,
		q.key,
		resp.Header.Revision,
		[]mvccpb.Event_EventType{mvccpb.PUT})
	if err != nil {
		return "", err
	}

	ok, err := deleteRevKey(q.client, string(ev.Kv.Key), ev.Kv.ModRevision)
	if err != nil {
		return "", err
	} else if !ok {
		return q.Dequeue()
	}
	return string(ev.Kv.Value), err
}
