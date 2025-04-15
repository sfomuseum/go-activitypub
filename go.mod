module github.com/sfomuseum/go-activitypub

go 1.23
toolchain go1.24.1

replace github.com/hpcloud/tail v1.0.0 => github.com/sfomuseum/tail v1.0.2

require (
	github.com/aaronland/go-aws-dynamodb v0.4.1
	github.com/aaronland/go-http-sanitize v0.0.8
	github.com/aaronland/go-http-server v1.5.0
	github.com/aaronland/go-pagination v0.3.0
	github.com/aaronland/go-pagination-sql v0.2.0
	github.com/aaronland/go-roster v1.0.0
	github.com/aaronland/gocloud-blob v0.3.1
	github.com/aaronland/gocloud-docstore v0.0.8
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go-v2 v1.32.3
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.36.3
	github.com/bwmarrin/snowflake v0.3.0
	github.com/fogleman/gg v1.3.0
	github.com/go-fed/httpsig v1.1.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/google/uuid v1.6.0
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/mitchellh/copystructure v1.2.0
	github.com/rs/cors v1.11.1
	github.com/sfomuseum/go-flags v0.10.0
	github.com/sfomuseum/go-pubsub v0.0.19
	github.com/sfomuseum/go-template v1.10.0
	github.com/sfomuseum/iso8601duration v1.1.0
	github.com/sfomuseum/runtimevar v1.2.1
	github.com/tidwall/gjson v1.18.0
	gocloud.dev v0.40.0
	golang.org/x/image v0.21.0
	golang.org/x/net v0.28.0
)

require (
	github.com/aaronland/go-aws-auth v1.6.5 // indirect
	github.com/aaronland/go-aws-session v0.2.1 // indirect
	github.com/aaronland/go-string v1.0.0 // indirect
	github.com/akrylysov/algnhsa v1.1.0 // indirect
	github.com/aws/aws-sdk-go v1.55.5 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.4 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.27.39 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.41 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.17 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.17.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/cognitoidentity v1.26.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/iam v1.36.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.59.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sns v1.31.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.36.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssm v1.54.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.32.2 // indirect
	github.com/aws/smithy-go v1.22.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/wire v0.6.0 // indirect
	github.com/googleapis/gax-go/v2 v2.13.0 // indirect
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jtacoma/uritemplates v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/redis/go-redis/v9 v9.7.0 // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/whosonfirst/go-ioutil v1.0.2 // indirect
	github.com/whosonfirst/go-sanitize v0.1.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/xerrors v0.0.0-20240716161551-93cc26a95ae9 // indirect
	google.golang.org/api v0.191.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240812133136-8ffd90a71988 // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)
