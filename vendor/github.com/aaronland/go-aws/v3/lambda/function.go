package lambda

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/aaronland/go-aws/v3/auth"
	"github.com/aws/aws-sdk-go-v2/aws"
	aws_lambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	aws_lambda_types "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type LambdaFunction struct {
	client    *aws_lambda.Client
	func_name string
	func_type string
}

func NewLambdaFunction(ctx context.Context, uri string) (*LambdaFunction, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	cfg_uri := fmt.Sprintf("aws://%s?credentials=%s", q.Get("region"), q.Get("credentials"))
	cfg, err := auth.NewConfig(ctx, cfg_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive AWS config, %w", err)
	}

	func_name := u.Host
	func_type := "Event"

	if q.Get("type") != "" {
		func_type = q.Get("type")
	}

	cl := aws_lambda.NewFromConfig(cfg)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new session, %w", err)
	}

	f := &LambdaFunction{
		client:    cl,
		func_name: func_name,
		func_type: func_type,
	}

	return f, nil
}

func (f *LambdaFunction) Invoke(ctx context.Context, payload any) (*aws_lambda.InvokeOutput, error) {

	enc_payload, err := json.Marshal(payload)

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal payload, %w", err)
	}

	return f.InvokeWithJSON(ctx, enc_payload)
}

func (f *LambdaFunction) InvokeWithJSON(ctx context.Context, payload []byte) (*aws_lambda.InvokeOutput, error) {

	var func_type aws_lambda_types.InvocationType
	var log_type aws_lambda_types.LogType

	switch strings.ToUpper(f.func_type) {
	case "EVENT":
		func_type = aws_lambda_types.InvocationTypeEvent
	case "REQUEST", "REQUESTRESPONSE":
		func_type = aws_lambda_types.InvocationTypeRequestResponse
	case "DRYRUN":
		func_type = aws_lambda_types.InvocationTypeDryRun
	default:
		return nil, fmt.Errorf("Invalid or unsupported invocation type")
	}

	switch func_type {
	case aws_lambda_types.InvocationTypeRequestResponse:
		log_type = aws_lambda_types.LogTypeTail
	default:
		log_type = aws_lambda_types.LogTypeNone
	}

	input := &aws_lambda.InvokeInput{
		FunctionName:   aws.String(f.func_name),
		InvocationType: func_type,
		LogType:        log_type,
		Payload:        payload,
	}

	rsp, err := f.client.Invoke(ctx, input)

	if err != nil {
		return nil, fmt.Errorf("Failed to invoke function %s (%s), %w", f.func_name, f.func_type, err)
	}

	if input.InvocationType != aws_lambda_types.InvocationTypeRequestResponse {
		return nil, nil
	}

	enc_result := *rsp.LogResult

	result, err := base64.StdEncoding.DecodeString(enc_result)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode result, %w", err)
	}

	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("Unexpected status code  %d (%s)", rsp.StatusCode, string(result))
	}

	return rsp, nil
}
