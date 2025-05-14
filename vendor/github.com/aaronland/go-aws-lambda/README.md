# go-aws-lambda

Opinionated Go package for doing things with AWS Lambda functions.

## Documentation

Documentation is incomplete at this time.

## Tools

$> make cli
go build -mod vendor -o bin/invoke cmd/invoke/main.go

### invoke

```
$> ./bin/invoke \
	-lambda-uri 'lambda://{FUNCTION_NAME}?region={AWS_REGION}&credentials={CREDENTIALS}' \
	-json '{JSON_ENCODED_ARGS}'
```

Where `{CREDENTIALS}` is expected to be a [aaronland/go-aws-session](https://github.com/aaronland/go-aws-session) credentials string:


| Label | Description |
| --- | --- |
| `anon:` | Empty or anonymous credentials. |
| `env:` | Read credentials from AWS defined environment variables. |
| `iam:` | Assume AWS IAM credentials are in effect. |
| `sts:{ARN}` | Assume the role defined by `{ARN}` using STS credentials. |
| `{AWS_PROFILE_NAME}` | This this profile from the default AWS credentials location. |
| `{AWS_CREDENTIALS_PATH}:{AWS_PROFILE_NAME}` | This this profile from a user-defined AWS credentials location. |
