package container

import (
	"github.com/hyicode/utils/assert"
	"testing"
)

type Int int

func (i Int) Less(j HeapElement) bool {
	return i < j.(Int)
}

func TestHeap(t *testing.T) {
	h := NewHeap[Int]()
	count := 100
	has := func(i Int) bool {
		for _, v := range h._h.sliceStack {
			if v == i {
				return true
			}
		}
		return false
	}
	for i := 0; i < count; i++ {
		h.Push(Int(i))
		assert.EqualFatalf(t, i+1, h.Len(), "push")
		assert.EqualFatalf(t, true, has(Int(i)), "has")
	}
	assert.EqualFatalf(t, count, h.Len(), "count")

	for i := 0; i < count; i++ {
		h.Pop()
		assert.EqualFatalf(t, count-i-1, h.Len(), "pop")
	}
}
