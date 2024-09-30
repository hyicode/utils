package container

import (
	"github.com/hyicode/utils/assert"
	"testing"
)

func TestSet(t *testing.T) {
	s := NewSet[int]()
	count := 100
	for i := 0; i < count; i++ {
		for j := 0; j < count; j++ {
			s.Add(i)
		}
		assert.EqualFatalf(t, i+1, s.Len(), "add")
		assert.EqualFatalf(t, true, s.Has(i), "add")
	}

	assert.EqualErrorf(t, count, s.Len(), "len")

	for i := 0; i < count*2; i++ {
		if i < count {
			assert.EqualFatalf(t, true, s.Has(i), "has")
		} else {
			assert.EqualFatalf(t, false, s.Has(i), "not has")
		}
	}

	for i := 0; i < count; i++ {
		for j := 0; j < count; j++ {
			s.Remove(i)
		}
		assert.EqualFatalf(t, count-i-1, s.Len(), "remove")
		assert.EqualFatalf(t, false, s.Has(i), "remove")
	}

}
