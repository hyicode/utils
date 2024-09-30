package sync

import "sync"

type Map[K comparable, V any] struct {
	_map sync.Map
}

func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	v, exist := m._map.Load(key)
	if !exist {
		var val V
		return val, false
	}
	return v.(V), true
}

func (m *Map[K, V]) Store(key K, value V) {
	m._map.Store(key, value)
}

func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	v, ok := m._map.Swap(key, value)
	if !ok {
		var val V
		return val, false
	}
	return v.(V), true
}

func (m *Map[K, V]) CompareAndSwap(key K, old, new V) bool {
	return m._map.CompareAndSwap(key, old, new)
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m._map.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

func (m *Map[K, V]) Delete(key K) {
	m._map.Delete(key)
}

func (m *Map[K, V]) CompareAndDelete(key K, old, new V) bool {
	return m._map.CompareAndSwap(key, old, new)
}

func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, exist := m._map.LoadAndDelete(key)
	if !exist {
		var val V
		return val, false
	}
	return v.(V), true
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, exist := m._map.LoadOrStore(key, value)
	if !exist {
		return value, false
	}
	return v.(V), true
}
