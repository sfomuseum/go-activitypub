package pagination

import (
	"github.com/jtacoma/uritemplates"
	"math"
)

// Method defines the type of pagination Options or Results.
type Method uint8

const (
	// Countable is a page-based or numbered pagination implementation
	Countable Method = iota
	// Cursor is a token or cursor-based pagination implementation where the total number of results is unknown
	Cursor
)

// Options provides a common interface for criteria used to paginate a query.
type Options interface {
	// Get or set the number of results to return in a Results
	PerPage(...int64) int64
	// Get or set current page to return query results for
	Spill(...int64) int64
	// Get or set the name of the column to use when paginating results
	Column(...string) string
	// Get or set the pointer (page number or cursor) for the next set of query results
	Pointer(...interface{}) interface{}
	// Return the Method type associated with Options
	Method() Method
}

// Results provides an interface for pagination information for a query response.
type Results interface {
	// The total number of results for a query.
	Total() int64
	// The number of results per page for a query.
	PerPage() int64
	// The current page number (offset) for a paginated query response.
	Page() int64
	// The total number of pages for a paginated query response.
	Pages() int64
	// The value to use to advance to the next set of results in a query response.
	Next() interface{}
	// The value to use to rewind to the previous set of results in a query response.
	Previous() interface{}
	// The URL to the next set of results in a query response
	NextURL(t *uritemplates.UriTemplate) (string, error)
	// The URL to the previous set of results in a query response
	PreviousURL(*uritemplates.UriTemplate) (string, error)
	// Return the Method type associated with Results
	Method() Method
}

// PagesForCount returns the number of pages that total_count will span using criteria defined in opts.
func PagesForCount(opts Options, total_count int64) int64 {

	per_page := int64(math.Max(1.0, float64(opts.PerPage())))
	spill := int64(math.Max(1.0, float64(opts.Spill())))

	if spill >= per_page {
		spill = per_page - 1
	}

	pages := int64(math.Ceil(float64(total_count) / float64(per_page)))
	return pages
}
