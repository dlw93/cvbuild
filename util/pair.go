package util

import (
	"fmt"
)

// Pair is just what its name suggests
type Pair[T, U any] struct {
	First  T
	Second U
}

// ValuePair is an alias for Pair[T, U] that exists for the sake of naming consistency.
type ValuePair[T, U any] struct {
	Pair[T, U]
}

// ReferencePair acts like a Pair[*T, *U] with an additional method to dereference its elements.
type ReferencePair[T, U any] struct {
	Pair[*T, *U]
}

// NewPair returns a pointer to a new pair.
func NewPair[T, U any](first T, second U) *Pair[T, U] {
	return &Pair[T, U]{first, second}
}

func NewValuePair[T, U any](first T, second U) *ValuePair[T, U] {
	p := NewPair(first, second)
	return &ValuePair[T, U]{*p}
}

func NewReferencePair[T, U any](first *T, second *U) *ReferencePair[T, U] {
	p := NewPair(first, second)
	return &ReferencePair[T, U]{*p}
}

// String returns a string representation of the pair.
func (p *Pair[T, U]) String() string {
	return fmt.Sprintf("(%v, %v)", p.First, p.Second)
}

// String returns a string representation of the pair.
func (p *ReferencePair[T, U]) String() string {
	return fmt.Sprintf("(%v, %v)", *p.First, *p.Second)
}

// Unpack returns the elements of the pair.
func (p *Pair[T, U]) Unpack() (T, U) {
	return p.First, p.Second
}

// Unpack returns the dereferenced elements of the pair.
func (p *ReferencePair[T, U]) Unpack() (T, U) {
	return *p.First, *p.Second
}

// Value returns a copy of the pair.
func (p *Pair[T, U]) Value() Pair[T, U] {
	return *p
}

// Value returns a copy of the pair with its elements dereferenced.
func (p *ReferencePair[T, U]) Value() Pair[T, U] {
	return Pair[T, U]{*p.First, *p.Second}
}
