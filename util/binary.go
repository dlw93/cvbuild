package util

// A JoinPredicate is a function that checks whether two values satisfy some condition.
type JoinPredicate[T, U any] func(T, U) bool

// Join returns an Iterable that iterates over all pairs of elements from `s` and `t` for which `f` evaluates to `true`.
func Join[T, U any](s []T, t []U, f JoinPredicate[*T, *U]) Iterable[Pair[*T, *U]] {
	return &join[T, U]{s, t, f}
}

// CartesianProduct returns an Iterable that iterates over all pairs of elements from `s` and `t`.
// It works like Join with a condition that always evaluates to `true`, but is faster because the condition is not actually evaluated.
func CartesianProduct[T, U any](s []T, t []U) Iterable[Pair[*T, *U]] {
	return &cartesianProduct[T, U]{s, t}
}
