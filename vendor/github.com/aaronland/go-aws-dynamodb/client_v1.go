package dynamodb

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"

	aa_session "github.com/aaronland/go-aws-session"
	"github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
)

// NewClientV1 returns an aws-sdk-go (v1) compatible client which is still necessary
// for some packages (like gocloud.dev/docstore)
func NewClientV1(ctx context.Context, uri string) (*aws_dynamodb.DynamoDB, error) {

	sess, err := newSessionWithURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create session, %w", err)
	}

	client := aws_dynamodb.New(sess)
	return client, nil
}

func newSessionWithURI(ctx context.Context, uri string) (*aws_session.Session, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()
	region := q.Get("region")
	credentials := q.Get("credentials")
	local := q.Get("local")

	is_local := false

	if local != "" {

		l, err := strconv.ParseBool(local)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?local parameter, %w", err)
		}

		is_local = l
	}

	if is_local {
		os.Setenv("AWS_ACCESS_KEY_ID", "local")
		os.Setenv("AWS_SECRET_ACCESS_KEY", "host")
		credentials = "env:"
		region = "localhost"
	}

	dsn := fmt.Sprintf("credentials=%s region=%s", credentials, region)
	sess, err := aa_session.NewSessionWithDSN(dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new session, %w", err)
	}

	if is_local {
		endpoint := "http://localhost:8000"
		sess.Config.Endpoint = aws.String(endpoint)
	}

	return sess, nil
}
