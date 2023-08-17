package clientv3

import (
	"testing"
	"time"

	"github.com/oldnicke/etcd/pkg/testutil"
)

func TestTxnPanics(t *testing.T) {
	defer testutil.AfterTest(t)

	kv := &kv{}

	errc := make(chan string)
	df := func() {
		if s := recover(); s != nil {
			errc <- s.(string)
		}
	}

	cmp := Compare(CreateRevision("foo"), "=", 0)
	op := OpPut("foo", "bar")

	tests := []struct {
		f func()

		err string
	}{
		{
			f: func() {
				defer df()
				kv.Txn(nil).If(cmp).If(cmp)
			},

			err: "cannot call If twice!",
		},
		{
			f: func() {
				defer df()
				kv.Txn(nil).Then(op).If(cmp)
			},

			err: "cannot call If after Then!",
		},
		{
			f: func() {
				defer df()
				kv.Txn(nil).Else(op).If(cmp)
			},

			err: "cannot call If after Else!",
		},
		{
			f: func() {
				defer df()
				kv.Txn(nil).Then(op).Then(op)
			},

			err: "cannot call Then twice!",
		},
		{
			f: func() {
				defer df()
				kv.Txn(nil).Else(op).Then(op)
			},

			err: "cannot call Then after Else!",
		},
		{
			f: func() {
				defer df()
				kv.Txn(nil).Else(op).Else(op)
			},

			err: "cannot call Else twice!",
		},
	}

	for i, tt := range tests {
		go tt.f()
		select {
		case err := <-errc:
			if err != tt.err {
				t.Errorf("#%d: got %s, wanted %s", i, err, tt.err)
			}
		case <-time.After(time.Second):
			t.Errorf("#%d: did not panic, wanted panic %s", i, tt.err)
		}
	}
}
