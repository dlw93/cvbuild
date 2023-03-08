package util

type Reference[T any] interface {
	Value() T
}

type Iterator[T any] interface {
	Next() bool
}

type Collectable[T any] interface {
	Collect() []T
}

type CollectableReferenceIterator[T any] interface {
	Iterator[Reference[T]]
	Collectable[T]
}

type CollectableValueIterator[T any] interface {
	Iterator[T]
	Collectable[T]
}

type Iterable[T any] interface {
	Iterator() (Iterator[T], *T)
}

type collectableReferenceIterator[T any] struct {
	Iterator[Reference[T]]
	v Reference[T]
}

func (it collectableReferenceIterator[T]) Collect() []T {
	var result []T
	for it.Next() {
		result = append(result, it.v.Value())
	}
	return result
}

type collectablePointerIterator[T any] struct {
	Iterator[T]
	v *T
}

func (it collectablePointerIterator[T]) Collect() []T {
	var result []T
	for it.Next() {
		result = append(result, *it.v)
	}
	return result
}

type collectableValueIterator[T any] struct {
	Iterator[T]
	v T
}

func (it collectableValueIterator[T]) Collect() []T {
	var result []T
	for it.Next() {
		result = append(result, it.v)
	}
	return result
}
