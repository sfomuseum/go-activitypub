# go-pubsub

Go package to provide a common interface for abstract publish and subscribe operations.

## Documentation

Documentation is incomplete at this time.

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/publish cmd/publish/main.go
go build -mod vendor -ldflags="-s -w" -o bin/subscribe cmd/subscribe/main.go
```

## Examples

```
$> ./bin/publish \
	-publisher-uri 'awssqs-creds://?region={REGION}&credentials={CREDENTIALS}&queue-url=https://sqs.{REGION}.amazonaws.com/{ACCOUNT}/{QUEUE}' \
	'hello world'
```

```
$> ./bin/subscribe \
	-subscriber-uri 'awssqs-creds://?region={REGION}&credentials={CREDENTIALS}&queue-url=https://sqs.{REGION}.amazonaws.com/{ACCOUNT}/{QUEUE}'
2024/09/04 10:59:57 INFO Listening for messages
hello world
```

## See also

* https://gocloud.dev/howto/pubsub/
* https://github.com/google/go-cloud/issues/1368
* https://github.com/go-redis/redis/v8