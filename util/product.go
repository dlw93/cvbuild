package util

type cartesianProduct[T, U any] struct {
	s []T
	t []U
}

type cartesianProductIterator[T, U any] struct {
	*cartesianProduct[T, U]
	i, j int
	v    ReferencePair[T, U] // points to s[i] and t[j]
}

func (it *cartesianProductIterator[T, U]) Next() bool {
	if it.i >= len(it.s) {
		return false
	}
	it.v.First = &it.s[it.i]
	it.v.Second = &it.t[it.j]
	it.j++
	if it.j == len(it.t) {
		it.i++
		it.j = 0
	}
	return true
}

func (cp cartesianProduct[T, U]) Iterator() (CollectableReferenceIterator[Pair[T, U]], *Pair[*T, *U]) {
	it := &cartesianProductIterator[T, U]{cartesianProduct: &cp}
	return &collectableReferenceIterator[Pair[T, U]]{it, &it.v}, &it.v.Pair
}
