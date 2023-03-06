package util

type Iterator[T any] interface {
	Next() bool
}

type Iterable[T any] interface {
	Iterator() (Iterator[T], *T)
}

type Collector[T any] interface {
	Collect() []T
}

type CollectableIterator[T any] interface {
	Iterator[T]
	Collector[T]
}

type CollectableIterable[T any] interface {
	Iterator() (CollectableIterator[T], *T)
}

type defaultCollectableIterator[T any] struct {
	Iterator[T]
	v *T
}

func (it defaultCollectableIterator[T]) Collect() []T {
	var result []T
	for it.Iterator.Next() {
		result = append(result, *it.v)
	}
	return result
}

type defaultCollectableIterable[T any] struct {
	Iterable[T]
}

func (c defaultCollectableIterable[T]) Iterator() (CollectableIterator[T], *T) {
	it, v := c.Iterable.Iterator()
	return &defaultCollectableIterator[T]{it, v}, v
}

func withDefaultCollector[T any](c Iterable[T]) CollectableIterable[T] {
	return &defaultCollectableIterable[T]{c}
}
