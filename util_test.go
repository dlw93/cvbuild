package main

import (
	"fmt"
	"strconv"
	"testing"

	"golang.org/x/exp/slices"
)

type pair[A, B comparable] struct {
	a A
	b B
}

func wrap[A, B comparable](a A, b B) pair[A, B] {
	return pair[A, B]{a, b}
}

type predicate[T any] func(T) bool

func newEqPredicate[T comparable](c T) predicate[T] {
	return func(v T) bool { return c == v }
}

func TestFind(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	v := []int{0, 1, 5, 9, 10}
	want := []pair[*int, bool]{{nil, false}, {&s[0], true}, {&s[4], true}, {&s[8], true}, {nil, false}}
	for i, e := range v {
		t.Run(fmt.Sprintf("Find %d", e), func(t *testing.T) {
			if got, ok := Find(s, newEqPredicate(e)); wrap(got, ok) != want[i] {
				t.Errorf("got %v, want %v", wrap(got, ok), want[i])
			}
		})
	}
}

func TestMap(t *testing.T) {
	s := []pair[int, int]{{1, 2}, {3, 4}, {5, 6}}
	f := func(p pair[int, int]) string { return strconv.Itoa(p.a + p.b) }
	want := []string{"3", "7", "11"}
	if got := Map(s, f); !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestPow(t *testing.T) {
	s := []pair[int, uint]{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {3, 4}}
	want := []int{1, 3, 9, 27, 81}
	for i, p := range s {
		t.Run(fmt.Sprintf("Pow(%d, %d) == %d", p.a, p.b, want[i]), func(t *testing.T) {
			if got := Pow(p.a, p.b); got != want[i] {
				t.Errorf("got %v, want %v", got, want[i])
			}
		})
	}
}
