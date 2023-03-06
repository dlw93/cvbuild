package util

type CollectionError string

const (
	ErrEmptySliceParameter = CollectionError("slice is empty")
)

func (e CollectionError) Error() string {
	return string(e)
}
