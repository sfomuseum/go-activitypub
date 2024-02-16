package countable

import (
	"github.com/aaronland/go-pagination"
)

const PER_PAGE int64 = 10
const PAGE int64 = 1
const SPILL int64 = 2
const COUNTABLE string = "*"

type CountableOptions struct {
	pagination.Options
	perpage int64
	page    int64
	spill   int64
	column  string
}

func NewCountableOptions() (pagination.Options, error) {

	opts := &CountableOptions{
		perpage: PER_PAGE,
		page:    PAGE,
		spill:   SPILL,
		column:  COUNTABLE,
	}

	return opts, nil
}

func (p *CountableOptions) Method() pagination.Method {
	return pagination.Countable
}

func (opts *CountableOptions) PerPage(args ...int64) int64 {

	if len(args) >= 1 {
		opts.perpage = args[0]
	}

	return opts.perpage
}

func (opts *CountableOptions) Pointer(args ...interface{}) interface{} {

	if len(args) >= 1 {
		opts.page = args[0].(int64)
	}

	return opts.page
}

func (opts *CountableOptions) Spill(args ...int64) int64 {

	if len(args) >= 1 {
		opts.spill = args[0]
	}

	return opts.spill
}

func (opts *CountableOptions) Column(args ...string) string {

	if len(args) >= 1 {
		opts.column = args[0]
	}

	return opts.column
}
