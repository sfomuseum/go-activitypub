module github.com/sfomuseum/go-activitypub

go 1.24.2

toolchain go1.24.3

replace github.com/hpcloud/tail v1.0.0 => github.com/sfomuseum/tail v1.0.2

require (
	github.com/aaronland/go-aws-dynamodb v0.4.2
	github.com/aaronland/go-http-sanitize v0.0.8
	github.com/aaronland/go-http-server v1.5.0
	github.com/aaronland/go-pagination v0.3.0
	github.com/aaronland/go-pagination-sql v0.2.0
	github.com/aaronland/go-roster v1.0.0
	github.com/aaronland/gocloud-blob v0.5.0
	github.com/aaronland/gocloud-docstore v0.0.9
	github.com/aws/aws-lambda-go v1.48.0
	github.com/aws/aws-sdk-go-v2 v1.36.3
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.43.1
	github.com/bwmarrin/snowflake v0.3.0
	github.com/fogleman/gg v1.3.0
	github.com/go-fed/httpsig v1.1.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/google/uuid v1.6.0
	github.com/mattn/go-sqlite3 v1.14.28
	github.com/mitchellh/copystructure v1.2.0
	github.com/rs/cors v1.11.1
	github.com/sfomuseum/go-flags v0.10.0
	github.com/sfomuseum/go-pubsub v0.0.20
	github.com/sfomuseum/go-template v1.10.1
	github.com/sfomuseum/iso8601duration v1.1.0
	github.com/sfomuseum/runtimevar v1.3.0
	github.com/tidwall/gjson v1.18.0
	gocloud.dev v0.41.0
	golang.org/x/image v0.27.0
	golang.org/x/net v0.40.0
)

require (
	github.com/aaronland/go-aws-auth v1.7.0 // indirect
	github.com/aaronland/go-aws-session v0.2.1 // indirect
	github.com/aaronland/go-string v1.0.0 // indirect
	github.com/akrylysov/algnhsa v1.1.0 // indirect
	github.com/aws/aws-sdk-go v1.55.6 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.10 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.29.12 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.65 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.17.69 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/service/cognitoidentity v1.27.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/iam v1.37.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.79.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sns v1.34.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.38.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssm v1.58.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.17 // indirect
	github.com/aws/smithy-go v1.22.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/google/wire v0.6.0 // indirect
	github.com/googleapis/gax-go/v2 v2.14.1 // indirect
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jtacoma/uritemplates v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/redis/go-redis/v9 v9.7.3 // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/whosonfirst/go-ioutil v1.0.2 // indirect
	github.com/whosonfirst/go-sanitize v0.1.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	google.golang.org/api v0.228.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/grpc v1.71.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)
