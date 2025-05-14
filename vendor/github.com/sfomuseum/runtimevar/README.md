# runtimevar

Simple wrapper around the Go Cloud runtimevar package

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/runtimevar.svg)](https://pkg.go.dev/github.com/sfomuseum/runtimevar)

## Example

```
package main

import (
	"context"
	"flag"
	"fmt"
	
	"github.com/sfomuseum/runtimevar"
)

func main() {

	flag.Parse()

	ctx := context.Background()

	for _, uri := range flag.Args() {
		str_var, _ := runtimevar.StringVar(ctx, uri)
		fmt.Printf(str_var)
	}
}
```

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/runtimevar cmd/runtimevar/main.go
```

### runtimevar

```
$> ./bin/runtimevar -h
Usage of ./bin/runtimevar:
  -timeout int
    	The maximum number of second in which a variable can be resolved. If 0 no timeout is applied.
```

#### Example

```
$> go run cmd/runtimevar/main.go 'constant://?val=hello+world'
hello world
```

## Supported services

The following Go Cloud `runtimevar` services are supported by the runtimevar tool by default:

* [AWS Parameter Store](https://gocloud.dev/howto/runtimevar/#awsps)
* [Blobvar](https://gocloud.dev/howto/runtimevar/#blob)
* [Local](https://gocloud.dev/howto/runtimevar/#local)

### AWS Parameter Store

It is possible to load runtime variables from AWS Parameter Store using [aaronland/go-aws-auth](https://github.com/aaronland/go-aws-auth) credential strings. For example:

```
$> go run cmd/runtimevar/main.go 'awsparamstore://hello-world?region=us-west-2&credentials=session'
hello world
```

Valid `aaronland/go-aws-auth` credential strings are:

Credentials for AWS sessions are defined as string labels. They are:

| Label | Description |
| --- | --- |
| `env:` | Read credentials from AWS defined environment variables. |
| `iam:` | Assume AWS IAM credentials are in effect. |
| `{AWS_PROFILE_NAME}` | Use the profile from the default AWS credentials location. |
| `{AWS_CREDENTIALS_PATH}:{AWS_PROFILE_NAME}` | Use the profile from a user-defined AWS credentials location. |

### Blob

The following [GoCloud Blob](https://gocloud.dev/howto/blob/) providers are supported by default:

* [S3](https://gocloud.dev/howto/blob/#s3)
* [Local](https://gocloud.dev/howto/blob/#local)

#### s3blob

In addition to the default `s3://` Blob source it is possible to load runtime variables from S3 buckets using [aaronland/go-aws-auth](https://github.com/aaronland/go-aws-auth) credential strings. For example:

```
$> go run cmd/runtimevar/main.go 'blobvar://hello-world?bucket-uri={BUCKET_URI}'
hello world
```

Where `{BUCKET_URI}` is a URL-escaped value like this:

```
s3blob://your-bucket?region=us-east-1&credentials=session
```

Valid `aaronland/go-aws-auth` credential strings are:

Credentials for AWS sessions are defined as string labels. They are:

| Label | Description |
| --- | --- |
| `env:` | Read credentials from AWS defined environment variables. |
| `iam:` | Assume AWS IAM credentials are in effect. |
| `{AWS_PROFILE_NAME}` | Use the profile from the default AWS credentials location. |
| `{AWS_CREDENTIALS_PATH}:{AWS_PROFILE_NAME}` | Use the profile from a user-defined AWS credentials location. |


## See also

* https://gocloud.dev/howto/runtimevar
* https://github.com/aaronland/go-aws-auth