package dynamodb

// Move this in to aaronland/go-aws-dynamodb

import (
	"context"
	"fmt"
	aa_session "github.com/aaronland/go-aws-session"
	"github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	aws_dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"net/url"
	"os"
	"strconv"
)

func NewSessionWithURI(ctx context.Context, uri string) (*aws_session.Session, error) {

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
		os.Setenv("AWS_ACCESS_KEY_ID", "DUMMYIDEXAMPLE")
		os.Setenv("AWS_SECRET_ACCESS_KEY", "DUMMYEXAMPLEKEY")
		credentials = "env:"
		region = "us-east-1"
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

func NewClientWithURI(ctx context.Context, uri string) (*aws_dynamodb.DynamoDB, error) {

	sess, err := NewSessionWithURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create session, %w", err)
	}

	client := aws_dynamodb.New(sess)
	return client, nil
}
