# create-post

Create a new post (note, activity) on behalf of a registered go-activitypub account and schedule it for delivery to all their followers.

```
$> ./bin/create-post -h
Create a new post (note, activity) on behalf of a registered go-activitypub account and schedule it for delivery to all their followers.
Usage:
	 ./bin/create-post [options]
Valid options are:
  -account-name string
    	The name of the go-activitypub account creating the post.
  -accounts-database-uri string
    	A registered sfomuseum/go-activitypub/database.AccountsDatabase URI. (default "null://")
  -activities-database-uri string
    	A registered sfomuseum/go-activitypub/database.ActivitiesDatabase URI. (default "null://")
  -deliveries-database-uri string
    	A registered sfomuseum/go-activitypub/database.DeliveriesDatabase URI. (default "null://")
  -delivery-queue-uri string
    	A registered sfomuseum/go-activitypub/queue/DeliveryQueue URI. (default "synchronous://")
  -followers-database-uri string
    	A registered sfomuseum/go-activitypub/database.FollowersDatabase URI. (default "null://")
  -hostname string
    	The hostname (domain) of the ActivityPub server delivering activities. (default "localhost:8080")
  -in-reply-to string
    	The URI of that the post is in reply to (optional).
  -insecure
    	A boolean flag indicating the ActivityPub server delivering activities is insecure (not using TLS).
  -lambda-function-uri string
    	A valid aaronland/go-aws-lambda.LambdaFunction URI in the form of "lambda://FUNCTION_NAME}?region={AWS_REGION}&credentials={CREDENTIALS}". This flag is required if the -mode flag is "invoke".
  -max-attempts int
    	The maximum number of attempts to deliver the activity. (default 5)
  -message string
    	The body (content) of the message to post.
  -mode string
    	The operating mode for creating new posts. Valid options are: cli, lambda and invoke, where "lambda" means to run as an AWS Lambda function and "invoke" means to invoke this tool as a specific Lambda function. (default "cli")
  -post-tags-database-uri string
    	A registered sfomuseum/go-activitypub/database.PostTagsDatabase URI. (default "null://")
  -posts-database-uri string
    	A registered sfomuseum/go-activitypub/database.PostsDatabase URI. (default "null://")
  -verbose
    	Enable verbose (debug) logging.
```

### Example

```
$> ./bin/create-post \
	-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
	-activities-database-uri '$(ACTIVITIES_DB_URI)' \
	-followers-database-uri '$(FOLLOWERS_DB_URI)' \
	-posts-database-uri '$(POSTS_DB_URI)' \
	-post-tags-database-uri '$(POST_TAGS_DB_URI)' \
	-deliveries-database-uri '$(DELIVERIES_DB_URI)' \
	-delivery-queue-uri '$(DELIVERY_QUEUE_URI)' \
	-account-name alice \
	-message "$(MESSAGE)" \
	-hostname localhost:8080 \
	-insecure \
	-verbose
```

_Note: It is acknowledged that it's kind of annoying to have to pass all those `*-database-uri` flags. There is not an immediate solution for this inconvenience but I am thinking about it._

## AWS

### Running as a Lambda function

First build the `create-post` tool as a Lambda function. The easiest way to do this is to run the `lambda-create-post` Makefile target:

```
$> make lambda-create-post
if test -f bootstrap; then rm -f bootstrap; fi
if test -f create-post.zip; then rm -f create-post.zip; fi
GOARCH=arm64 GOOS=linux go build -mod vendor -ldflags="-s -w" -tags lambda.norpc -o bootstrap cmd/create-post/main.go
zip create-post.zip bootstrap
  adding: bootstrap (deflated 74%)
rm -f bootstrap
```

Installing and configuring the Lambda function is outside the scope of this documentation. At a minimum you will need to assign, at a minimu, the following environment variables:

| Key | Value |
| --- | --- |
| ACTIVITYPUB_ACCOUNTS_DATABASE_URI   | awsdynamodb://{ACCOUNTS_TABLE}?partition_key=Id&allow_scans=true&region={REGION}&credentials=iam: |
| ACTIVITYPUB_ACTIVITIES_DATABASE_URI | awsdynamodb://{ACTIVITIES_TABLES}?partition_key=Id&allow_scans=true&region={REGION}&credentials=iam: |
| ACTIVITYPUB_DELIVERIES_DATABASE_URI | awsdynamodb://{DELIVERIES_TABLES}?partition_key=Id&allow_scans=true&region={REGION}&credentials=iam: |
| ACTIVITYPUB_DELIVERY_QUEUE_URI | awssqs-creds://?region={REGION}&credentials=iam:&queue-url={SQS_QUEUE_URI} |
| ACTIVITYPUB_HOSTNAME                | {YOUR_DOMAIN} |
| ACTIVITYPUB_MODE                    | lambda |
| ACTIVITYPUB_POST_TAGS_DATABASE_URI  | awsdynamodb://{POST_TAGS_TABLE}?partition_key=Id&allow_scans=true&region={REGION}&credentials=iam: |
| ACTIVITYPUB_POSTS_DATABASE_URI      | awsdynamodb://{POSTS_TABLE}?partition_key=Id&allow_scans=true&region={REGION}&credentials=iam: |
| ACTIVITYPUB_VERBOSE                 | {BOOLEAN} |

Enviromment variables map to regular command line flags using the following rules:

* Replace all instances of `-` with `_`
* Upper-case the flag string
* Prepend the string with `SFOMUSEUM_`

For example the flag `-accounts-database-uri` becomes the `SFOMUSEUM_ACCOUNTS_DATABASE_URI` environment variable.

_Note: This example assumes an Amazon DynamoDB (Docstore) database implementation and an Amazon SQS backed delivery queue since they are both easy to integrate with Lambda functions but they are just examples. Adjust according to your needs and circumstances._

### Invoking (running from) a Lambda function

To invoke the `create-post` tool running as a Lambda function (like the one installed above) you do the following:

```
$> ./bin/create-post \
	-mode invoke \
	-account-name testbot \
	-message 'This is a second test' \
	-lambda-function-uri 'lambda://{FUNCTION_NAME}?region={AWS_REGION}&credentials={CREDENTIALS}&type=REQUEST'
```

### Credentials

All of the AWS-related URI declarations use the [aaronland/go-aws-auth](https://github.com/aaronland/go-aws-auth?tab=readme-ov-file#credentials) package to define and load AWS credentials. This package differs from the default AWS credential-ing options in that it provides more authentication options defined as string labels. They are:

| Label | Description |
| --- | --- |
| `anon:` | Empty or anonymous credentials. |
| `env:` | Read credentials from AWS defined environment variables. |
| `iam:` | Assume AWS IAM credentials are in effect. |
| `sts:{ARN}` | Assume the role defined by `{ARN}` using STS credentials. |
| `{AWS_PROFILE_NAME}` | This this profile from the default AWS credentials location. |
| `{AWS_CREDENTIALS_PATH}:{AWS_PROFILE_NAME}` | This this profile from a user-defined AWS credentials location. |
