package main

import (
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Unsigned constraints.Unsigned

// Find returns a poiner to the first element from `s` for which `f` evaluates to `true`.
// If no such element exists, the second return value is `false`.
func Find[E any](s []E, f func(E) bool) (*E, bool) {
	for i := range s {
		if f(s[i]) {
			return &s[i], true
		}
	}
	return nil, false
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
func Pow[T Number, U Unsigned](a T, b U) T {
	var c T = 1
	for b > 0 {
		if b&1 != 0 {
			c *= a
		}
		a *= a
		b >>= 1
	}
	return c
}
