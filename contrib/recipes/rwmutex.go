package recipe

import (
	"context"

	v3 "oldnicke/etcd/clientv3"
	"oldnicke/etcd/clientv3/concurrency"
	"oldnicke/etcd/mvcc/mvccpb"
)

type RWMutex struct {
	s   *concurrency.Session
	ctx context.Context

	pfx   string
	myKey *EphemeralKV
}

func NewRWMutex(s *concurrency.Session, prefix string) *RWMutex {
	return &RWMutex{s, context.TODO(), prefix + "/", nil}
}

func (rwm *RWMutex) RLock() error {
	rk, err := newUniqueEphemeralKey(rwm.s, rwm.pfx+"read")
	if err != nil {
		return err
	}
	rwm.myKey = rk
	// wait until nodes with "write-" and a lower revision number than myKey are gone
	for {
		if done, werr := rwm.waitOnLastRev(rwm.pfx + "write"); done || werr != nil {
			return werr
		}
	}
}

func (rwm *RWMutex) Lock() error {
	rk, err := newUniqueEphemeralKey(rwm.s, rwm.pfx+"write")
	if err != nil {
		return err
	}
	rwm.myKey = rk
	// wait until all keys of lower revision than myKey are gone
	for {
		if done, werr := rwm.waitOnLastRev(rwm.pfx); done || werr != nil {
			return werr
		}
		//  get the new lowest key until this is the only one left
	}
}

// waitOnLowest will wait on the last key with a revision < rwm.myKey.Revision with a
// given prefix. If there are no keys left to wait on, return true.
func (rwm *RWMutex) waitOnLastRev(pfx string) (bool, error) {
	client := rwm.s.Client()
	// get key that's blocking myKey
	opts := append(v3.WithLastRev(), v3.WithMaxModRev(rwm.myKey.Revision()-1))
	lastKey, err := client.Get(rwm.ctx, pfx, opts...)
	if err != nil {
		return false, err
	}
	if len(lastKey.Kvs) == 0 {
		return true, nil
	}
	// wait for release on blocking key
	_, err = WaitEvents(
		client,
		string(lastKey.Kvs[0].Key),
		rwm.myKey.Revision(),
		[]mvccpb.Event_EventType{mvccpb.DELETE})
	return false, err
}

func (rwm *RWMutex) RUnlock() error { return rwm.myKey.Delete() }
func (rwm *RWMutex) Unlock() error  { return rwm.myKey.Delete() }
