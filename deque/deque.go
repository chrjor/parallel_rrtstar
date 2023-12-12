package deque

import "sync/atomic"

// DEQue interface
type DEQue[T any] interface {
	PushBottom(task *T)
	IsEmpty() bool
	PopTop() *T
	PopBottom() *T
}

// unboundedDEQue struct
type unboundedDEQue[T any] struct {
	tasks  *CircularArray[T] // Circular array of tasks
	bottom int32             // Bottom index
	top    atomic.Int32      // Top index
}

// Create a new unboundedDEQue
func NewUnboundedDEQue[T any]() DEQue[T] {
	return &unboundedDEQue[T]{
		tasks:  NewCircularArray[T](32),
		bottom: 0,
		top:    atomic.Int32{},
	}
}

// Check if deque empty
func (d *unboundedDEQue[T]) IsEmpty() bool {
	return d.bottom <= d.top.Load()
}

// Push to bottom (non-concurrent)
func (d *unboundedDEQue[T]) PushBottom(task *T) {
	oldBottom := d.bottom
	oldTop := d.top.Load()
	size := oldBottom - oldTop
	if size >= d.tasks.Capacity()-1 {
		d.tasks = d.tasks.Resize(oldTop, oldBottom)
	}
	d.tasks.Put(oldBottom, task)
	d.bottom++
}

// Pop from bottom (concurrent only with PopTop)
func (d *unboundedDEQue[T]) PopBottom() *T {
	d.bottom--
	oldTop := d.top.Load()
	newTop := oldTop + 1
	size := d.bottom - oldTop
	if size < 0 {
		d.bottom = oldTop
		return nil
	}
	task := d.tasks.Get(d.bottom)
	if size > 0 {
		return task
	}
	// Only use compare and sway when task is last
	if !d.top.CompareAndSwap(oldTop, newTop) {
		task = nil
	}
	d.bottom = newTop
	return task
}

// Pop from top (concurrent)
func (d *unboundedDEQue[T]) PopTop() *T {
	oldTop := d.top.Load()
	newTop := oldTop + 1
	size := d.bottom - oldTop
	if size <= 0 {
		return nil
	}
	task := d.tasks.Get(oldTop)
	if d.top.CompareAndSwap(oldTop, newTop) {
		return task
	}
	return nil
}
