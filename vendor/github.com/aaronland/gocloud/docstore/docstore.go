package docstore

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	aa_dynamodb "github.com/aaronland/go-aws/v3/dynamodb"
	"gocloud.dev/docstore"
	"gocloud.dev/docstore/awsdynamodb/v2"
)

const DYNAMODB_FALLBACK_FUNC_KEY string = "aaronland-dynamodb-fallback-func"

func OpenCollection(ctx context.Context, uri string) (*docstore.Collection, error) {

	var col *docstore.Collection

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse '%s', %w", uri, err)
	}

	if u.Scheme == "awsdynamodb" {

		table := u.Host

		q := u.Query()

		partition_key := q.Get("partition_key")
		q_allow_scans := q.Get("allow_scans")

		cl_q := url.Values{}
		cl_q.Set("region", q.Get("region"))
		cl_q.Set("credentials", q.Get("credentials"))

		if q.Has("local") {
			cl_q.Set("local", q.Get("local"))
		}

		cl_uri := url.URL{}
		cl_uri.Scheme = "dynamodb"
		cl_uri.RawQuery = cl_q.Encode()

		cl, err := aa_dynamodb.NewClient(ctx, cl_uri.String())

		if err != nil {
			return nil, fmt.Errorf("Failed to create DynamoDB client (%s), %w", cl_uri.String(), err)
		}

		col_opts := &awsdynamodb.Options{}

		if q_allow_scans != "" {

			allow, err := strconv.ParseBool(q_allow_scans)

			if err != nil {
				return nil, fmt.Errorf("Failed to parse ?allow_scans= parameter, %w", err)
			}

			col_opts.AllowScans = allow

			v := ctx.Value(DYNAMODB_FALLBACK_FUNC_KEY)

			if v != nil {

				switch v.(type) {
				case func() interface{}:
					// pass
				default:
					return nil, fmt.Errorf("Invalid fallback func %T", v)
				}

				fn := v.(func() interface{})

				fallback_func := awsdynamodb.InMemorySortFallback(fn)
				col_opts.RunQueryFallback = fallback_func
			}

		}

		c, err := awsdynamodb.OpenCollection(cl, table, partition_key, "", col_opts)

		if err != nil {
			return nil, fmt.Errorf("Failed to open collection, %w", err)
		}

		col = c

	} else {

		c, err := docstore.OpenCollection(ctx, uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to open collection, %w", err)
		}

		col = c
	}

	return col, nil
}
