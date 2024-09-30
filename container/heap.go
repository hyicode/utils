package container

import (
	"container/heap"
)

type HeapElement interface {
	Less(HeapElement) bool
}

type _heap[T HeapElement] struct {
	sliceStack[T]
}

func (h *_heap[T]) Len() int {
	return h.len()
}

func (h *_heap[T]) Less(i, j int) bool {
	return (h).sliceStack[i].Less((h).sliceStack[j])
}

func (h *_heap[T]) Swap(i, j int) {
	(h.sliceStack)[i], (h.sliceStack)[j] = (h.sliceStack)[j], (h.sliceStack)[i]
}

func (h *_heap[T]) Push(x any) {
	v := x.(T)
	h.push(v)
}

func (h *_heap[T]) Pop() any {
	return h.pop()
}

type Heap[T HeapElement] struct {
	_h _heap[T]
}

func NewHeap[T HeapElement]() *Heap[T] {
	return new(Heap[T]).Init()
}

func (h *Heap[T]) Init() *Heap[T] {
	heap.Init(&h._h)
	return h
}

func (h *Heap[T]) Push(x T) {
	heap.Push(&h._h, x)
}

func (h *Heap[T]) Pop() T {
	return heap.Pop(&h._h).(T)
}

func (h *Heap[T]) Remove(i int) T {
	return heap.Remove(&h._h, i).(T)
}

func (h *Heap[T]) Len() int {
	return h._h.Len()
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[T]) Fix(i int) {
	heap.Fix(&h._h, i)
}
