package mockwait

import (
	"go.etcd.io/etcd/pkg/testutil"
	"go.etcd.io/etcd/pkg/wait"
)

type WaitRecorder struct {
	wait.Wait
	testutil.Recorder
}

type waitRecorder struct {
	testutil.RecorderBuffered
}

func NewRecorder() *WaitRecorder {
	wr := &waitRecorder{}
	return &WaitRecorder{Wait: wr, Recorder: wr}
}
func NewNop() wait.Wait { return NewRecorder() }

func (w *waitRecorder) Register(id uint64) <-chan interface{} {
	w.Record(testutil.Action{Name: "Register"})
	return nil
}
func (w *waitRecorder) Trigger(id uint64, x interface{}) {
	w.Record(testutil.Action{Name: "Trigger"})
}

func (w *waitRecorder) IsRegistered(id uint64) bool {
	panic("waitRecorder.IsRegistered() shouldn't be called")
}
