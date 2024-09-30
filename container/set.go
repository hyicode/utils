package container

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Len() int { return len(s) }

func (s Set[T]) Add(v T) { s[v] = struct{}{} }

func (s Set[T]) Remove(v T) {
	delete(s, v)
}

func (s Set[T]) Has(v T) bool {
	_, exist := s[v]
	return exist
}

func (s Set[T]) Range(f func(v T) (stop bool)) {
	for k := range s {
		if f(k) {
			return
		}
	}
}
