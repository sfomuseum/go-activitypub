package runtimevar

import (
	"context"
	"fmt"
	"net/url"

	_ "github.com/aaronland/gocloud-blob/s3"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	_ "gocloud.dev/blob/s3blob"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"

	"github.com/aaronland/go-aws-auth"
	"github.com/aaronland/gocloud-blob/bucket"
	gc "gocloud.dev/runtimevar"
	"gocloud.dev/runtimevar/awsparamstore"
	"gocloud.dev/runtimevar/blobvar"
)

// StringVar returns the latest string value contained by 'uri', which is expected
// to be a valid `gocloud.dev/runtimevar` URI.
func StringVar(ctx context.Context, uri string) (string, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return "", fmt.Errorf("Failed to parse URI, %w", err)
	}

	if u.Scheme == "" {
		return u.Path, nil
	}

	q := u.Query()

	if q.Get("decoder") == "" {
		q.Set("decoder", "string")
		u.RawQuery = q.Encode()
	}

	var v *gc.Variable
	var v_err error

	switch u.Scheme {
	case "awsparamstore":

		// https://gocloud.dev/howto/runtimevar/#awsps-ctor

		creds := q.Get("credentials")
		region := q.Get("region")

		if creds != "" {

			aws_uri := fmt.Sprintf("aws://%s?credentials=%s", region, creds)
			aws_auth, err := auth.NewSSMClient(ctx, aws_uri)

			if err != nil {
				return "", fmt.Errorf("Failed to create AWS session credentials, %w", err)
			}

			v, v_err = awsparamstore.OpenVariableV2(aws_auth, u.Host, gc.StringDecoder, nil)
		}

	case "blobvar":

		if !q.Has("bucket-uri") {
			return "", fmt.Errorf("Missing ?bucket-uri parameter")
		}

		b_uri, err := url.QueryUnescape(q.Get("bucket-uri"))

		if err != nil {
			return "", fmt.Errorf("Failed to unescape bucket URI, %w", err)
		}

		b, err := bucket.OpenBucket(ctx, b_uri)

		if err != nil {
			return "", fmt.Errorf("Failed to open bucket, %w", err)
		}

		defer b.Close()

		v, v_err = blobvar.OpenVariable(b, u.Host, gc.StringDecoder, nil)

	default:
		// pass
	}

	if v == nil {

		uri = u.String()

		v, v_err = gc.OpenVariable(ctx, uri)
	}

	if v_err != nil {
		return "", fmt.Errorf("Failed to open variable, %w", v_err)
	}

	defer v.Close()

	snapshot, err := v.Latest(ctx)

	if err != nil {
		return "", fmt.Errorf("Failed to derive latest snapshot for variable, %w", err)
	}

	return snapshot.Value.(string), nil
}
