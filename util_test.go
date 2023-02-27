package main

import (
	"fmt"
	"strconv"
	"testing"

	"golang.org/x/exp/slices"
)

func TestFind(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	v := []int{0, 1, 5, 9, 10}
	want := []*int{nil, &s[0], &s[4], &s[8], nil}
	for i, e := range v {
		t.Run(fmt.Sprintf("Find %d", e), func(t *testing.T) {
			if got := Find(s, NewEqualityPredicate(e)); got != want[i] {
				t.Errorf("got %v, want %v", got, want[i])
			}
		})
	}
}

func TestMap(t *testing.T) {
	s := []Pair[int, int]{{1, 2}, {3, 4}, {5, 6}}
	f := func(p Pair[int, int]) string { return strconv.Itoa(p.First + p.Second) }
	want := []string{"3", "7", "11"}
	if got := Map(s, f); !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestPow(t *testing.T) {
	s := []Pair[int, uint]{{3, 0}, {3, 1}, {3, 2}, {3, 3}, {3, 4}}
	want := []int{1, 3, 9, 27, 81}
	for i, p := range s {
		t.Run(fmt.Sprintf("Pow(%d, %d) == %d", p.First, p.Second, want[i]), func(t *testing.T) {
			if got := Pow(p.First, p.Second); got != want[i] {
				t.Errorf("got %v, want %v", got, want[i])
			}
		})
	}
}
