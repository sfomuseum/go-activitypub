package dynamodb

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/aaronland/go-aws-auth"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	aws_dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewClientWithURI(ctx context.Context, uri string) (*aws_dynamodb.Client, error) {
	return NewClient(ctx, uri)
}

func NewClient(ctx context.Context, uri string) (*aws_dynamodb.Client, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	cfg, err := auth.NewConfig(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create config, %w", err)
	}

	client_opts := make([]func(*aws_dynamodb.Options), 0)

	is_local := false

	if q.Has("local") {

		v, err := strconv.ParseBool(q.Get("local"))

		if err != nil {
			return nil, fmt.Errorf("Invalid ?local= parameter, %w", err)
		}

		is_local = v
	}

	// https://dave.dev/blog/2021/07/14-07-2021-awsddb/
	if is_local {

		cfg.Region = "localhost"

		cfg.EndpointResolver = aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000", SigningRegion: "localhost"}, nil
			})

		creds_opts := func(o *aws_dynamodb.Options) {
			o.Credentials = credentials.NewStaticCredentialsProvider("local", "host", "")
		}

		client_opts = append(client_opts, creds_opts)
	}

	client := aws_dynamodb.NewFromConfig(cfg, client_opts...)
	return client, nil
}
