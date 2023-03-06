package util

type cartesianProduct[T, U any] Pair[[]T, []U]

type cartesianProductIterator[T, U any] struct {
	*cartesianProduct[T, U]
	i, j int
	v    Pair[*T, *U] // points to s[i] and t[j]
}

func (it *cartesianProductIterator[T, U]) Next() bool {
	if it.i >= len(it.First) {
		return false
	}
	it.v.First = &it.First[it.i]
	it.v.Second = &it.Second[it.j]
	it.j++
	if it.j == len(it.Second) {
		it.i++
		it.j = 0
	}
	return true
}

func (cp cartesianProduct[T, U]) Iterator() (Iterator[Pair[*T, *U]], *Pair[*T, *U]) {
	it := &cartesianProductIterator[T, U]{cartesianProduct: &cp}
	return it, &it.v
}

type join[T, U any] struct {
	s []T
	t []U
	f JoinPredicate[*T, *U]
}

type joinIterator[T, U any] struct {
	*join[T, U]
	i, j int
	v    Pair[*T, *U] // points to join.s[i] and join.t[j]
}

func (j join[T, U]) Iterator() (Iterator[Pair[*T, *U]], *Pair[*T, *U]) {
	it := &joinIterator[T, U]{join: &j}
	return it, &it.v
}

func (it *joinIterator[T, U]) Next() bool {
	if it.i >= len(it.s) {
		return false
	}
	s, t := &it.s[it.i], &it.t[it.j]
	for !it.f(s, t) {
		it.j++
		if it.j == len(it.t) {
			it.i++
			it.j = 0
			if it.i == len(it.s) {
				return false
			}
		}
	}
	it.v.First = s
	it.v.Second = t
	return true
}
