package container

import (
	"github.com/hyicode/utils/assert"
	"testing"
)

func TestStack(t *testing.T) {
	s := NewStack[int]()
	count := 100
	for i := 0; i < count; i++ {
		s.Push(i)
		assert.EqualFatalf(t, i+1, s.Len(), "push")
	}
	expectV := 0
	s.Range(func(v int) (stop bool) {
		assert.EqualFatalf(t, expectV, v, "range")
		expectV++
		return false
	})
	assert.EqualFatalf(t, expectV, count, "expectV")
	for i := count; i > 0; i-- {
		v := s.Pop()
		assert.EqualFatalf(t, i-1, v, "pop")
	}
}
