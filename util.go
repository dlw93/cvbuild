package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type Unsigned constraints.Unsigned

type Number interface {
	constraints.Integer | constraints.Float
}

type Pair[T, U comparable] struct {
	First  T
	Second U
}

// NewPair returns a pointer to a new pair.
func NewPair[T, U comparable](first T, second U) *Pair[T, U] {
	return &Pair[T, U]{first, second}
}

// String returns a string representation of the pair.
func (p *Pair[T, U]) String() string {
	return fmt.Sprintf("(%v, %v)", p.First, p.Second)
}

// Unwrap returns the elements of the pair.
func (p *Pair[T, U]) Unwrap() (T, U) {
	return p.First, p.Second
}

// A Predicate is a function that checks whether a value satisfies some condition.
type Predicate[T any] func(T) bool

// A Reducer is a function that combines two values into one.
type Reducer[T, U any] func(U, T) U

// NewEqualityPredicate returns a predicate that checks whether a value is equal to `c`.
func NewEqualityPredicate[T comparable](c T) Predicate[T] {
	return func(v T) bool { return c == v }
}

// Find returns a poiner to the first element from `s` for which `f` evaluates to `true` or `nil` if no such element exists.
func Find[E any](s []E, f Predicate[E]) *E {
	for i := range s {
		if f(s[i]) {
			return &s[i]
		}
	}
	return nil
}

// Map applies `f` to each element of slice `s` and returns a new slice containing the results.
func Map[A, B any](s []A, f func(A) B) []B {
	t := make([]B, len(s))
	for i := range s {
		t[i] = f(s[i])
	}
	return t
}

// Reduce applies `f` to each element of slice `s` and the current value of `a` and returns the final value of `a`.
//
// If `a` is not provided and types `A` and `B` are equal, the first element of `s` is used as the initial accumulator value.
// In this case, if `s` is empty, Reduce panics.
func Reduce[A, B any](s []A, f Reducer[A, B], a ...B) B {
	var r B
	if len(a) > 0 {
		r = a[0]
	} else if len(s) > 0 {
		if p, ok := any(&s[0]).(*B); ok {
			r = *p
			s = s[1:]
		} else {
			err := fmt.Errorf("Reduce: initial value for accumulator `a` unspecified, but type parameters `A` and `B` have different types %T and %T", s[0], r)
			panic(err)
		}
	} else {
		err := fmt.Errorf("Reduce: initial value for accumulator `a` unspecified, but slice `s` is empty")
		panic(err)
	}
	for i := range s {
		r = f(r, s[i])
	}
	return r
}

// Pow computes `a` to the power of `b` in O(log2 b) time.
func Pow[B Number, E Unsigned](a B, b E) B {
	var c B = 1
	for b > 0 {
		if b&1 != 0 {
			c *= a
		}
		a *= a
		b >>= 1
	}
	return c
}
