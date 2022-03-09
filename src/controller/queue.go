package main

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

var (
	// 队列满载.
	ErrEmptyQueue = errors.New("queue is empty")
	// 队列为空.
	ErrFullQueue = errors.New("queue is full")
)

type waiter chan interface{}

func newWaiter() waiter {
	w := make(chan interface{}, 1)
	return w
}

type Queue struct {
	maxSize     int
	mute        bool
	mutex       sync.Mutex
	items       *list.List // store items
	putters     *list.List // store blocked Put operators
	getters     *list.List // store blocked Get operators
	update_time time.Time
}

// maxSize参数可设置队列大小
// 若maxSize为0，那么队列是无限大，并且向队列中加元素无需等待
func NewQueue(maxSize int) *Queue {
	q := new(Queue)
	q.mutex = sync.Mutex{}
	q.maxSize = maxSize
	q.items = list.New()
	q.putters = list.New()
	q.getters = list.New()
	q.update_time = time.Now()
	return q
}

func (q *Queue) newPutter() *list.Element {
	w := newWaiter()
	return q.putters.PushBack(w)
}

func (q *Queue) newGetter() *list.Element {
	w := newWaiter()
	return q.getters.PushBack(w)
}

func (q *Queue) notifyPutter(getter *list.Element) bool {
	if getter != nil {
		q.getters.Remove(getter)
	}
	if q.putters.Len() == 0 {
		return false
	}
	e := q.putters.Front()
	q.putters.Remove(e)
	w := e.Value.(waiter)
	w <- true
	return true
}

func (q *Queue) notifyGetter(putter *list.Element, val interface{}) bool {
	if putter != nil {
		q.putters.Remove(putter)
	}
	if q.getters.Len() == 0 {
		return false
	}
	e := q.getters.Front()
	q.getters.Remove(e)
	w := e.Value.(waiter)
	w <- val
	return true
}

func (q *Queue) clearPending() {
	for !q.isfull() && q.putters.Len() != 0 {
		q.notifyPutter(nil)
	}
	for !q.isempty() && q.getters.Len() != 0 {
		v := q.get()
		q.notifyGetter(nil, v)
	}
}

func (q *Queue) get() interface{} {
	e := q.items.Front()
	q.items.Remove(e)
	return e.Value
}

func (q *Queue) put(val interface{}) {
	q.items.PushBack(val)
}

// Same as Get(-1).
func (q *Queue) GetNoWait() (interface{}, error) {
	return q.Get(-1)
}

// * If timeout less than 0, If Queue is empty, return (nil, ErrEmptyQueue).
//
// * If timeout equals to 0, block until get a value from Queue.
//
// * If timeout greater than 0, wait timeout seconds until get a value from Queue,
// if timeout passed, return (nil, ErrEmptyQueue).
func (q *Queue) Get(timeout float64) (interface{}, error) {
	q.mutex.Lock()
	q.clearPending()
	isempty := q.isempty()
	if timeout < 0.0 && isempty {
		q.mutex.Unlock()
		return nil, ErrEmptyQueue
	}

	if !isempty {
		v := q.get()
		q.notifyPutter(nil)
		q.mutex.Unlock()
		return v, nil
	}

	e := q.newGetter()
	q.mutex.Unlock()
	w := e.Value.(waiter)

	var v interface{}
	if timeout == 0.0 {
		v = <-w
	} else {
		select {
		case v = <-w:
		case <-time.After(time.Duration(timeout) * time.Second):
			return nil, ErrEmptyQueue
		}
	}
	q.mutex.Lock()
	q.notifyPutter(e)
	q.mutex.Unlock()
	return v, nil
}

// Same as Put(-1).
func (q *Queue) PutNoWait(val interface{}) error {
	return q.Put(val, -1)
}

// * If timeout less than 0, If Queue is full, return (nil, ErrFullQueue).
//
// * If timeout equals to 0, block until put a value into Queue.
//
// * If timeout greater than 0, wait timeout seconds until put a value into Queue,
// if timeout passed, return (nil, ErrFullQueue).
func (q *Queue) Put(val interface{}, timeout float64) error {
	q.mutex.Lock()
	q.clearPending()
	isfull := q.isfull()
	if timeout < 0.0 && isfull {
		q.mutex.Unlock()
		return ErrFullQueue
	}

	if !isfull {
		if !q.notifyGetter(nil, val) {
			q.put(val)
		}
		q.mutex.Unlock()
		return nil
	}

	e := q.newPutter()
	q.mutex.Unlock()
	w := e.Value.(waiter)
	if timeout == 0.0 {
		<-w
	} else {
		select {
		case <-w:
		case <-time.After(time.Duration(timeout) * time.Second):
			return ErrFullQueue
		}
	}

	q.mutex.Lock()
	if !q.notifyGetter(e, val) {
		q.put(e)
	}
	q.mutex.Unlock()
	return nil
}

func (q *Queue) size() int {
	return q.items.Len()
}

// Return size of Queue.
func (q *Queue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.size()
}

func (q *Queue) isempty() bool {
	return (q.size() == 0)
}

// Return true if Queue is empty.
func (q *Queue) IsEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.isempty()
}

func (q *Queue) isfull() bool {
	return (q.maxSize > 0 && q.maxSize <= q.size())
}

// Return true if Queue is full.
func (q *Queue) IsFull() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.isfull()
}

func (q *Queue) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.items = q.items.Init()
}
