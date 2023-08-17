package recipe

import (
	"context"

	v3 "oldnicke/etcd/clientv3"
	"oldnicke/etcd/mvcc/mvccpb"
)

// Barrier creates a key in etcd to block processes, then deletes the key to
// release all blocked processes.
type Barrier struct {
	client *v3.Client
	ctx    context.Context

	key string
}

func NewBarrier(client *v3.Client, key string) *Barrier {
	return &Barrier{client, context.TODO(), key}
}

// Hold creates the barrier key causing processes to block on Wait.
func (b *Barrier) Hold() error {
	_, err := newKey(b.client, b.key, v3.NoLease)
	return err
}

// Release deletes the barrier key to unblock all waiting processes.
func (b *Barrier) Release() error {
	_, err := b.client.Delete(b.ctx, b.key)
	return err
}

// Wait blocks on the barrier key until it is deleted. If there is no key, Wait
// assumes Release has already been called and returns immediately.
func (b *Barrier) Wait() error {
	resp, err := b.client.Get(b.ctx, b.key, v3.WithFirstKey()...)
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		// key already removed
		return nil
	}
	_, err = WaitEvents(
		b.client,
		b.key,
		resp.Header.Revision,
		[]mvccpb.Event_EventType{mvccpb.PUT, mvccpb.DELETE})
	return err
}
