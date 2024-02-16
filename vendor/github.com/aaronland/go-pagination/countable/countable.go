// package countable provides implementions of the pagintion.Options and pagination.Results interfaces for use with page-based or numbered pagination.
package countable

import (
	"github.com/aaronland/go-pagination"
)

func NextPage(r pagination.Results) int64 {

	if r.Method() != pagination.Countable {
		return 0
	}

	return r.Next().(int64)
}

func PreviousPage(r pagination.Results) int64 {

	if r.Method() != pagination.Countable {
		return 0
	}

	return r.Previous().(int64)
}

func PageFromOptions(opts pagination.Options) int64 {

	if opts.Method() != pagination.Countable {
		return 0
	}

	return opts.Pointer().(int64)
}
