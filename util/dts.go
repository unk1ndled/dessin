package util

type Deque[T any] struct {
	items []T
}

func NewDeque[T any]() *Deque[T] {
	return &Deque[T]{items: []T{}}
}

// PushFront adds an item to the front of the deque.
func (d *Deque[T]) PushFront(item T) {
	d.items = append([]T{item}, d.items...)
}

// PushBack adds an item to the back of the deque.
func (d *Deque[T]) PushBack(item T) {
	d.items = append(d.items, item)
}

// PopFront removes and returns the item from the front of the deque.
func (d *Deque[T]) PopFront() (T, bool) {
	if len(d.items) == 0 {
		var zero T
		return zero, false
	}
	item := d.items[0]
	d.items[0] = *new(T) // Clear the reference ie setting the discarded value to a zero of type T
	d.items = d.items[1:]
	return item, true
}

// PopBack removes and returns the item from the back of the deque.
func (d *Deque[T]) PopBack() (T, bool) {
	if len(d.items) == 0 {
		var zero T
		return zero, false
	}
	index := len(d.items) - 1
	item := d.items[index]
	d.items[index] = *new(T) // Clear the reference
	d.items = d.items[:index]
	return item, true
}

// PeekFront returns the item at the front of the deque without removing it.
func (d *Deque[T]) PeekFront() (T, bool) {
	if len(d.items) == 0 {
		var zero T
		return zero, false
	}
	return d.items[0], true
}

// PeekBack returns the item at the back of the deque without removing it.
func (d *Deque[T]) PeekBack() (T, bool) {
	if len(d.items) == 0 {
		var zero T
		return zero, false
	}
	return d.items[len(d.items)-1], true
}

// Size returns the number of items in the deque.
func (d *Deque[T]) Size() int {
	return len(d.items)
}

// IsEmpty checks if the deque is empty.
func (d *Deque[T]) IsEmpty() bool {
	return len(d.items) == 0
}
