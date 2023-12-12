package deque

// CircularArray struct used internally by the unboundedDEQue
type CircularArray[T any] struct {
	c     int32 // Current capacity
	tasks []*T  // Tasks array
}

// Create new circular array
func NewCircularArray[T any](capacity int32) *CircularArray[T] {
	return &CircularArray[T]{
		c:     capacity,
		tasks: make([]*T, capacity),
	}
}

// Check capacity
func (a *CircularArray[T]) Capacity() int32 {
	return a.c
}

// Get item from array
func (a *CircularArray[T]) Get(item_idx int32) *T {
	return a.tasks[item_idx%a.c]
}

// Put item onto array
func (a *CircularArray[T]) Put(item_idx int32, item *T) {
	a.tasks[item_idx%a.c] = item
}

// Resize array and copy task references
func (a *CircularArray[T]) Resize(top int32, bottom int32) *CircularArray[T] {
	newArray := NewCircularArray[T](a.c * 2)
	for item_idx := top; item_idx < bottom; item_idx++ {
		newArray.Put(item_idx, a.Get(item_idx))
	}
	return newArray
}
