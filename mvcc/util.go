package mvcc

import (
	"encoding/binary"
	"fmt"

	"oldnicke/etcd/mvcc/backend"
	"oldnicke/etcd/mvcc/mvccpb"
)

func UpdateConsistentIndex(be backend.Backend, index uint64) {
	tx := be.BatchTx()
	tx.Lock()
	defer tx.Unlock()

	var oldi uint64
	_, vs := tx.UnsafeRange(metaBucketName, consistentIndexKeyName, nil, 0)
	if len(vs) != 0 {
		oldi = binary.BigEndian.Uint64(vs[0])
	}

	if index <= oldi {
		return
	}

	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, index)
	tx.UnsafePut(metaBucketName, consistentIndexKeyName, bs)
}

func WriteKV(be backend.Backend, kv mvccpb.KeyValue) {
	ibytes := newRevBytes()
	revToBytes(revision{main: kv.ModRevision}, ibytes)

	d, err := kv.Marshal()
	if err != nil {
		panic(fmt.Errorf("cannot marshal event: %v", err))
	}

	be.BatchTx().Lock()
	be.BatchTx().UnsafePut(keyBucketName, ibytes, d)
	be.BatchTx().Unlock()
}
