package util

type CollectionError string

func (e CollectionError) Error() string {
	return string(e)
}

const (
	ErrEmptySliceParameter = CollectionError("slice is empty")
)
