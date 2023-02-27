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

func (p *Pair[T, U]) Unwrap() (T, U) {
	return p.First, p.Second
}

// A Predicate is a function that checks whether a value satisfies some condition.
type Predicate[T any] func(T) bool

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

// Map applies `f` to each element of `s` and returns a new slice containing the results.
func Map[A, B any](s []A, f func(A) B) []B {
	t := make([]B, len(s))
	for i := range s {
		t[i] = f(s[i])
	}
	return t
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
