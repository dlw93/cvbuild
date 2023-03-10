package util

import "unsafe"

// A Selector[T, U any] is a function that returns a pointer to a field of type U in a struct of type T.
type Selector[T, U any] func(*T) *U

// A Offset[T, U any] is the offset of a field of type U in a struct of type T.
type Offset[T, U any] uint

// Get returns a pointer to the field of type U in struct t of type T that is referred to by this Offset.
//
// That is, a call to o.Get(&t) on an Offset[T, U] o created by a call to newOffset[T, U](func(*T) *U { return &t.u }) returns &t.u.
//
//	type Type struct {
//	    // ...
//		id uint
//	    // ...
//	}
//	id := newOffset(func(*Type) *uint { return &t.id })
//
//	t := Type{id: 42}
//	_ = *id.Get(&t) == t.id // true
//
// Note that the Offset is not validated prior to accessing the referred field and thus may cause a panic if it does not point to a field of type U inside the given struct t of type T.
func (o Offset[T, U]) Get(t *T) *U {
	return (*U)(unsafe.Add(unsafe.Pointer(t), o))
}

// newOffset returns the offset of the field of type U in a struct of type T that is referred to by the given Selector.
//
// If the address returned by the Selector is not within the bounds of the struct, newOffset panics with ErrSelectorOffsetOutOfBounds.
func newOffset[T, U any](s Selector[T, U]) Offset[T, U] {
	var x T
	y := s(&x)
	xp := uintptr(unsafe.Pointer(&x))
	yp := uintptr(unsafe.Pointer(y))
	if xp <= yp && yp < xp+unsafe.Sizeof(x) {
		return Offset[T, U](yp - xp)
	}
	panic(ErrSelectorOffsetOutOfBounds)
}

type JoinEqualityPredicate[R, S any, T comparable] struct {
	Left  Offset[R, T]
	Right Offset[S, T]
}

func (p *JoinEqualityPredicate[R, S, T]) Evaluate(r *R, s *S) bool {
	return *p.Left.Get(r) == *p.Right.Get(s)
}

func newJoinEqualityPredicate[R, S any, T comparable](s Selector[R, T], t Selector[S, T]) *JoinEqualityPredicate[R, S, T] {
	return &JoinEqualityPredicate[R, S, T]{newOffset(s), newOffset(t)}
}

func Join2[R, S any, T comparable](r []R, s []S, p Selector[R, T], q Selector[S, T]) []Pair[R, S] {
	result := make([]Pair[R, S], 0, len(r)*len(s))
	predicate := newJoinEqualityPredicate(p, q)
	for _, r := range r {
		for _, s := range s {
			if predicate.Evaluate(&r, &s) {
				return append(result, Pair[R, S]{r, s})
			}
		}
	}
	return result
}
