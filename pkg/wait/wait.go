// Package wait provides utility functions for polling, listening using Go
// channel.
package wait

import (
	"log"
	"sync"
)

// Wait is an interface that provides the ability to wait and trigger events that
// are associated with IDs.
type Wait interface {
	// Register waits returns a chan that waits on the given ID.
	// The chan will be triggered when Trigger is called with
	// the same ID.
	Register(id uint64) <-chan interface{}
	// Trigger triggers the waiting chans with the given ID.
	Trigger(id uint64, x interface{})
	IsRegistered(id uint64) bool
}

type list struct {
	l sync.RWMutex
	m map[uint64]chan interface{}
}

// New creates a Wait.
func New() Wait {
	return &list{m: make(map[uint64]chan interface{})}
}

func (w *list) Register(id uint64) <-chan interface{} {
	w.l.Lock()
	defer w.l.Unlock()
	ch := w.m[id]
	if ch == nil {
		ch = make(chan interface{}, 1)
		w.m[id] = ch
	} else {
		log.Panicf("dup id %x", id)
	}
	return ch
}

func (w *list) Trigger(id uint64, x interface{}) {
	w.l.Lock()
	ch := w.m[id]
	delete(w.m, id)
	w.l.Unlock()
	if ch != nil {
		ch <- x
		close(ch)
	}
}

func (w *list) IsRegistered(id uint64) bool {
	w.l.RLock()
	defer w.l.RUnlock()
	_, ok := w.m[id]
	return ok
}

type waitWithResponse struct {
	ch <-chan interface{}
}

func NewWithResponse(ch <-chan interface{}) Wait {
	return &waitWithResponse{ch: ch}
}

func (w *waitWithResponse) Register(id uint64) <-chan interface{} {
	return w.ch
}
func (w *waitWithResponse) Trigger(id uint64, x interface{}) {}
func (w *waitWithResponse) IsRegistered(id uint64) bool {
	panic("waitWithResponse.IsRegistered() shouldn't be called")
}
