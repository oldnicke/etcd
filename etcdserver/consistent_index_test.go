package etcdserver

import "testing"

func TestConsistentIndex(t *testing.T) {
	var i consistentIndex
	i.setConsistentIndex(10)
	if g := i.ConsistentIndex(); g != 10 {
		t.Errorf("value = %d, want 10", g)
	}
}
