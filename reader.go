package activitypub

import (
	"io"
)

func DefaultLimitedReader(r io.Reader) io.Reader {

	n := int64(1024 * 1024)
	return NewLimitedReader(r, n)
}

func NewLimitedReader(r io.Reader, n int64) io.Reader {

	return &io.LimitedReader{
		R: r,
		N: n,
	}
}
