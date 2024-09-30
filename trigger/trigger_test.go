package trigger

import (
	"fmt"
	"github.com/hyicode/utils/assert"
	"runtime"
	"strconv"
	"testing"
)

func TestNewEventTable(t *testing.T) {
	const (
		TestKey0 EventName[string]  = "test_1"
		TestKey1 EventName[int]     = "test_2"
		TestKey2 EventName[float64] = "test_3"
	)
	table := NewEventTable()
	TestKey0.On(table, func(event string) {
		assert.EqualErrorf(t, "hello", event, "")
	})
	TestKey1.On(table, func(event int) {
		assert.EqualErrorf(t, 1, event, "")
	})
	TestKey2.On(table, func(event float64) {
		assert.EqualErrorf(t, 2.0, event, "")
	})

	TestKey0.Trigger(table, "hello")
	TestKey1.Trigger(table, 1)
	TestKey2.Trigger(table, 2)
}

func BenchmarkNewEventTable(b *testing.B) {
	const (
		KeyNum = 100000
		CbNum  = 100
	)

	table := NewEventTable()
	keys := make([]EventName[int], 0, KeyNum)
	counter := 0
	for i := 0; i < KeyNum; i++ {
		key := EventName[int]("test_" + strconv.Itoa(i))
		for j := 0; j < CbNum; j++ {
			key.On(table, func(event int) {
				counter++
			})
		}
		keys = append(keys, key)
	}

	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := keys[i%KeyNum]
		key.Trigger(table, i)
	}
	fmt.Println(counter)
}
