package util

type CollectionError uint

const (
	_ CollectionError = iota

	// ErrEmptySliceParameter occurs when a function is called with an empty slice parameter where a non-empty slice is required.
	ErrEmptySliceParameter

	// ErrSelectorOffsetOutOfBounds occurs when the pointer returned by a Selector[T, U] does not point to a field inside the struct that was passed to it.
	ErrSelectorOffsetOutOfBounds
)

func (e CollectionError) Error() string {
	switch e {
	case ErrEmptySliceParameter:
		return "empty slice parameter where a non-empty slice is required"
	case ErrSelectorOffsetOutOfBounds:
		return "the pointer returned by the selector does not point to a field inside the struct that was passed to it"
	default:
		return "unknown error"
	}
}
