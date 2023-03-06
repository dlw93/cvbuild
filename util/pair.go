package util

import (
	"fmt"
)

type Pair[T, U any] struct {
	First  T
	Second U
}

// NewPair returns a pointer to a new pair.
func NewPair[T, U any](first T, second U) *Pair[T, U] {
	return &Pair[T, U]{first, second}
}

// String returns a string representation of the pair.
func (p *Pair[T, U]) String() string {
	return fmt.Sprintf("(%v, %v)", p.First, p.Second)
}

// Unpack returns the elements of the pair.
func (p *Pair[T, U]) Unpack() (T, U) {
	return p.First, p.Second
}
