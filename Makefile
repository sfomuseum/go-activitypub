GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

MYSQL=mysql
SQLITE3=sqlite3

DYNAMODB_TABLE_PREFIX=

TAGS=null

SUPPORTED_DATABASES=dynamodb sqlite mysql

DATABASE=mysql
DATABASE_LOWER=$(shell echo $(DATABASE) | tr '[:upper:]' '[:lower:]')
DATABASE_UPPER=$(shell echo $(DATABASE) | tr '[:lower:]' '[:upper:]')

ifeq ($(filter $(DATABASE_LOWER),$(SUPPORTED_DATABASES)),)
$(error "DATABASE is undefined or not one of $(SUPPORTED_DATABASES).  Set DATABASE=…")
endif

DELIVERY_QUEUE_URI=synchronous://

migrate:
	GOARCH=amd64 GOOS=linux go build -tags mysql -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/migrate cmd/migrate/main.go
	GOARCH=amd64 GOOS=linux go build -tags mysql -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/list-accounts cmd/list-accounts/main.go

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/add-account cmd/add-account/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/add-aliases cmd/add-aliases/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/block cmd/block/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/boost-note cmd/boost-note/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/counts-for-date cmd/counts-for-date/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/create-dynamodb-tables cmd/create-dynamodb-tables/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/create-post cmd/create-post/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/deliver-activity cmd/deliver-activity/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/get-account cmd/get-account/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/follow cmd/follow/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/list-boosts cmd/list-boosts/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/list-followers cmd/list-followers/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/list-activities cmd/list-activities/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/list-addresses cmd/list-addresses/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/list-aliases cmd/list-aliases/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/list-deliveries cmd/list-deliveries/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/inbox cmd/inbox/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/retrieve-actor cmd/retrieve-actor/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/retrieve-delivery cmd/retrieve-delivery/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/retrieve-note cmd/retrieve-note/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/server cmd/server/main.go

lambda:
	@make lambda-server
	@make lambda-create-post
	@make lambda-deliver-activity

lambda-server:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags lambda.norpc -o bootstrap cmd/server/main.go
	zip server.zip bootstrap
	rm -f bootstrap

lambda-create-post:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f create-post.zip; then rm -f create-post.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags lambda.norpc -o bootstrap cmd/create-post/main.go
	zip create-post.zip bootstrap
	rm -f bootstrap

lambda-deliver-activity:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f deliver.zip; then rm -f deliver.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags lambda.norpc -o bootstrap cmd/deliver-activity/main.go
	zip deliver.zip bootstrap
	rm -f bootstrap

# The rest of these Makefile targets are for local testing

ACCOUNTS_SQLITE_DB=work/accounts.db
ALIASES_SQLITE_DB=work/aliases.db
ACTIVITIES_SQLITE_DB=work/activities.db
FOLLOWERS_SQLITE_DB=work/followers.db
FOLLOWING_SQLITE_DB=work/following.db
POSTS_SQLITE_DB=work/posts.db
POST_TAGS_SQLITE_DB=work/posts_tags.db
NOTES_SQLITE_DB=work/notes.db
MESSAGES_SQLITE_DB=work/messages.db
BLOCKS_SQLITE_DB=work/blocks.db
DELIVERIES_SQLITE_DB=work/deliveries.db
BOOSTS_SQLITE_DB=work/boosts.db
LIKES_SQLITE_DB=work/liks.db
PROPERTIES_SQLITE_DB=work/properties.db

ACCOUNTS_SQLITE_URI=sql://sqlite3?dsn=file:$(ACCOUNTS_SQLITE_DB)%3Fcache%3Dshared
ALIASES_SQLITE_URI=sql://sqlite3?dsn=file:$(ALIASES_SQLITE_DB)%3Fcache%3Dshared
ACTIVITIES_SQLITE_URI=sql://sqlite3?dsn=file:$(ACTIVITIES_SQLITE_DB)%3Fcache%3Dshared
FOLLOWERS_SQLITE_URI=sql://sqlite3?dsn=file:$(FOLLOWERS_SQLITE_DB)%3Fcache%3Dshared
FOLLOWING_SQLITE_URI=sql://sqlite3?dsn=file:$(FOLLOWING_SQLITE_DB)%3Fcache%3Dshared
BLOCKS_SQLITE_URI=sql://sqlite3?dsn=file:$(BLOCKS_SQLITE_DB)%3Fcache%3Dshared
POSTS_SQLITE_URI=sql://sqlite3?dsn=file:$(POSTS_SQLITE_DB)%3Fcache%3Dshared
POST_TAGS_SQLITE_URI=sql://sqlite3?dsn=file:$(POST_TAGS_SQLITE_DB)%3Fcache%3Dshared
NOTES_SQLITE_URI=sql://sqlite3?dsn=file:$(NOTES_SQLITE_DB)%3Fcache%3Dshared
MESSAGES_SQLITE_URI=sql://sqlite3?dsn=file:$(MESSAGES_SQLITE_DB)%3Fcache%3Dshared
DELIVERIES_SQLITE_URI=sql://sqlite3?dsn=file:$(DELIVERIES_SQLITE_DB)%3Fcache%3Dshared
BOOSTS_SQLITE_URI=sql://sqlite3?dsn=file:$(BOOSTS_SQLITE_DB)%3Fcache%3Dshared
LIKES_SQLITE_URI=sql://sqlite3?dsn=file:$(LIKES_SQLITE_DB)%3Fcache%3Dshared
PROPERTIES_SQLITE_URI=sql://sqlite3?dsn=file:$(PROPERTIES_SQLITE_DB)cache%3Dshared

# constant://?val=user:password
MYSQL_CREDENTIALS=constant%3A%2F%2F%3Fval%3Duser%3Apassword
MYSQL_DSN={credentials}@/activitypub
MYSQL_URI=sql://mysql?dsn=$(MYSQL_DSN)&credentials-uri=$(MYSQL_CREDENTIALS)

ACCOUNTS_MYSQL_URI=$(MYSQL_URI)
ALIASES_MYSQL_URI=$(MYSQL_URI)
ACTIVITIES_MYSQL_URI=$(MYSQL_URI)
FOLLOWERS_MYSQL_URI=$(MYSQL_URI)
FOLLOWING_MYSQL_URI=$(MYSQL_URI)
BLOCKS_MYSQL_URI=$(MYSQL_URI)
POSTS_MYSQL_URI=$(MYSQL_URI)
POST_TAGS_MYSQL_URI=$(MYSQL_URI)
NOTES_MYSQL_URI=$(MYSQL_URI)
MESSAGES_MYSQL_URI=$(MYSQL_URI)
DELIVERIES_MYSQL_URI=$(MYSQL_URI)
BOOSTS_MYSQL_URI=$(MYSQL_URI)
LIKES_MYSQL_URI=$(MYSQL_URI)
PROPERTIES_MYSQL_URI=$(MYSQL_URI)

ACCOUNTS_DYNAMODB_URI=awsdynamodb://$(DYNAMODB_TABLE_PREFIX)accounts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
ACTIVITIES_DYNAMODB_URI=awsdynamodb://$(DYNAMODB_TABLE_PREFIX)activities?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
ALIASES_DYNAMODB_URI=awsdynamodb://$(DYNAMODB_TABLE_PREFIX)aliases?partition_key=Name&allow_scans=true&local=true&region=localhost&credentials=anon:
BLOCKS_DYNAMODB_URI=awsdynamodb://$(DYNAMODB_TABLE_PREFIX)blocks?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
BOOSTS_DYNAMODB_URI=awsdynamodb://$(DYNAMODB_TABLE_PREFIX)boosts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
DELIVERIES_DYNAMODB_URI=awsdynamodb://$(DYNAMODB_TABLE_PREFIX)deliveries?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
FOLLOWING_DYNAMODB_URI=awsdynamodb://$(DYNAMODB_TABLE_PREFIX)following?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
FOLLOWERS_DYNAMODB_URI=awsdynamodb://$(DYNAMODB_TABLE_PREFIX)followers?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
LIKES_DYNAMODB_URI=awsdynamodb://$(DYNAMODB_TABLE_PREFIX)likes?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
NOTES_DYNAMODB_URI=awsdynamodb://$(DYNAMODB_TABLE_PREFIX)notes?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
MESSAGES_DYNAMODB_URI=awsdynamodb://$(DYNAMODB_TABLE_PREFIX)messages?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
POST_TAGS_DYNAMODB_URI=awsdynamodb://$(DYNAMODB_TABLE_PREFIX)post_tags?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
POSTS_DYNAMODB_URI=awsdynamodb://$(DYNAMODB_TABLE_PREFIX)posts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:
PROPERTIES_DYNAMODB_URI=awsdynamodb://$(DYNAMODB_TABLE_PREFIX)properties?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:

ACCOUNTS_DB_URI=$(ACCOUNTS_$(DATABASE_UPPER)_URI)
ACTIVITIES_DB_URI=$(ACTIVITIES_$(DATABASE_UPPER)_URI)
ALIASES_DB_URI=$(ALIASES_$(DATABASE_UPPER)_URI)
BLOCKS_DB_URI=$(BLOCKS_$(DATABASE_UPPER)_URI)
BOOSTS_DB_URI=$(BOOSTS_$(DATABASE_UPPER)_URI)
DELIVERIES_DB_URI=$(DELIVERIES_$(DATABASE_UPPER)_URI)
FOLLOWING_DB_URI=$(FOLLOWING_$(DATABASE_UPPER)_URI)
FOLLOWERS_DB_URI=$(FOLLOWERS_$(DATABASE_UPPER)_URI)
LIKES_DB_URI=$(LIKES_$(DATABASE_UPPER)_URI)
NOTES_DB_URI=$(NOTES_$(DATABASE_UPPER)_URI)
MESSAGES_DB_URI=$(MESSAGES_$(DATABASE_UPPER)_URI)
POST_TAGS_DB_URI=$(POST_TAGS_$(DATABASE_UPPER)_URI)
POSTS_DB_URI=$(POSTS_$(DATABASE_UPPER)_URI)
PROPERTIES_DB_URI=$(PROPERTIES_$(DATABASE_UPPER)_URI)

local-tables-dynamodb:
	go run -mod $(GOMOD) cmd/create-dynamodb-tables/main.go \
		-refresh \
		-table-prefix '$(DYNAMODB_TABLE_PREFIX)' \
		-dynamodb-client-uri 'awsdynamodb://?region=localhost&credentials=anon:&local=true'

local-tables-mysql:
	$(MYSQL) -u$(MYSQL_USER) -p activitypub < schema/mysql/activitypub.schema

local-tables-sqlite:
	rm -f *.db
	$(SQLITE3) $(ACCOUNTS_DB) < schema/sqlite/accounts.schema
	$(SQLITE3) $(ALIASES_DB) < schema/sqlite/aliases.schema
	$(SQLITE3) $(FOLLOWERS_DB) < schema/sqlite/followers.schema
	$(SQLITE3) $(FOLLOWING_DB) < schema/sqlite/following.schema
	$(SQLITE3) $(POSTS_DB) < schema/sqlite/posts.schema
	$(SQLITE3) $(POST_TAGS_DB) < schema/sqlite/post_tags.schema
	$(SQLITE3) $(NOTES_DB) < schema/sqlite/notes.schema
	$(SQLITE3) $(MESSAGES_DB) < schema/sqlite/messages.schema
	$(SQLITE3) $(BLOCKS_DB) < schema/sqlite/blocks.schema
	$(SQLITE3) $(BOOSTS_DB) < schema/sqlite/boosts.schema
	$(SQLITE3) $(LIKES_DB) < schema/sqlite/likes.schema
	$(SQLITE3) $(PROPERTIES_DB) < schema/sqlite/properties.schema
	$(SQLITE3) $(DELIVERIES_DB) < schema/sqlite/deliveries.schema

local-tables:
	@make local-tables-$(DATABASE)

local-setup:
	@make local-tables
	@make local-accounts
	@make local-server

deliver-pubsub:
	go run cmd/deliver-activity/main.go \
		-mode pubsub \
		-subscriber-uri 'redis://?channel=activitypub' \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-activities-database-uri '$(ACTIVITIES_DB_URI)' \
		-deliveries-database-uri '$(DELIVERIES_DB_URI)' \
		-followers-database-uri '$(FOLLOWERS_DB_URI)' \
		-posts-database-uri '$(POSTS_DB_URI)' \
		-post-tags-database-uri '$(POST_TAGS_DB_URI)' \
		-insecure \
		-verbose

local-accounts:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/add-account/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-aliases-database-uri '$(ALIASES_DB_URI)' \
		-properties-database-uri '$(PROPERTIES_DB_URI)' \
		-account-name bob \
		-alias robert \
		-account-type Service \
		-account-icon-uri fixtures/icons/bob.jpg \
		-property 'url:www=https://bob.com' \
		-embed-icon-uri
	go run -mod $(GOMOD) -tags $(TAGS) cmd/add-account/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-aliases-database-uri '$(ALIASES_DB_URI)' \
		-properties-database-uri '$(PROPERTIES_DB_URI)' \
		-account-name doug \
		-alias doug \
		-property 'url:www=https://bob.com/doug' \
		-account-type Service
	go run -mod $(GOMOD) -tags $(TAGS) cmd/add-account/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-aliases-database-uri '$(ALIASES_DB_URI)' \
		-properties-database-uri '$(PROPERTIES_DB_URI)' \
		-account-name alice \
		-account-type Person \
		-property 'url:www=https://www.alice.info' \
		-account-icon-uri 's3blob://sfomuseum-media/ap/icons/sfo.jpg?region=us-west-2&credentials=session'
		# -allow-remote-icon-uri \
		# -account-icon-uri https://static.sfomuseum.org/media/172/956/659/5/1729566595_kjcAQKRw176gxIieIWZySjhlNzgKNxoA_s.jpg

# Bob wants to follow Alice

follow:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/follow/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-following-database-uri '$(FOLLOWING_DB_URI)' \
		-messages-database-uri '$(MESSAGES_DB_URI)' \
		-account-name bob \
		-follow alice@localhost:8080 \
		-hostname localhost:8080 \
		-verbose \
		-insecure

# Bob wants to unfollow Alice

unfollow:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/follow/main.go  \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-following-database-uri '$(FOLLOWING_DB_URI)' \
		-messages-database-uri '$(MESSAGES_DB_URI)' \
		-account-name bob \
		-follow alice@localhost:8080 \
		-hostname localhost:8080 \
		-insecure \
		-verbose \
		-undo

block:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/block/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-blocks-database-uri '$(BLOCKS_DB_URI)' \
		-account-name bob \
		-block-host block.club \
		-verbose

unblock:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/block/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-blocks-database-uri '$(BLOCKS_DB_URI)' \
		-account-name bob \
		-block-host block.club \
		-undo \
		-verbose

# Alice wants to post something (to Bob, if Bob is following Alice)

post:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/create-post/main.go \
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

boost-note:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/boost-note/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-activities-database-uri '$(ACTIVITIES_DB_URI)' \
		-followers-database-uri '$(FOLLOWERS_DB_URI)' \
		-deliveries-database-uri '$(DELIVERIES_DB_URI)' \
		-account-name doug \
		-note "$(NOTE)" \
		-hostname localhost:8080 \
		-delivery-queue-uri '$(DELIVERY_QUEUE_URI)' \
		-insecure \
		-verbose

list-boosts:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/list-boosts/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-boosts-database-uri '$(BOOSTS_DB_URI)' \
		-account-name $(ACCOUNT) \
		-hostname localhost:8080 \
		-insecure \
		-verbose

# -mention $(MENTION) \

reply:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/create-post/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-followers-database-uri '$(FOLLOWERS_DB_URI)' \
		-posts-database-uri '$(POSTS_DB_URI)' \
		-post-tags-database-uri '$(POST_TAGS_DB_URI)' \
		-deliveries-database-uri '$(DELIVERIES_DB_URI)' \
		-account-name bob \
		-message "$(MESSAGE)" \
		-in-reply-to $(INREPLYTO) \
		-hostname localhost:8080 \
		-insecure \
		-verbose

delivery:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/retrieve-delivery/main.go \
		-deliveries-database-uri '$(DELIVERIES_DB_URI)' \
		-delivery-id $(ID) \
		-verbose

list-inbox:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/inbox/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-messages-database-uri '$(MESSAGES_DB_URI)' \
		-notes-database-uri '$(NOTES_DB_URI)' \
		-account-name $(ACCOUNT) \
		-verbose

SERVER_DISABLED=false
SERVER_VERBOSE=true

local-server:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/server/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-aliases-database-uri '$(ALIASES_DB_URI)' \
		-followers-database-uri '$(FOLLOWERS_DB_URI)' \
		-following-database-uri '$(FOLLOWING_DB_URI)' \
		-notes-database-uri '$(NOTES_DB_URI)' \
		-messages-database-uri '$(MESSAGES_DB_URI)' \
		-blocks-database-uri '$(BLOCKS_DB_URI)' \
		-posts-database-uri '$(POSTS_DB_URI)' \
		-post-tags-database-uri '$(POST_TAGS_DB_URI)' \
		-boosts-database-uri '$(BOOSTS_DB_URI)' \
		-likes-database-uri '$(LIKES_DB_URI)' \
		-properties-database-uri '$(PROPERTIES_DB_URI)' \
		-process-message-queue-uri 'stdout://' \
		-allow-remote-icon-uri \
		-allow-create \
		-verbose=$(SERVER_VERBOSE) \
		-disabled=$(SERVER_DISABLED) \
		-hostname localhost:8080 \
		-insecure

list-activities:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/list-activities/main.go \
		-activities-database-uri '$(ACTIVITIES_DB_URI)' \
		-verbose

list-deliveries:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/list-deliveries/main.go \
		-deliveries-database-uri '$(DELIVERIES_DB_URI)' \
		-verbose

retrieve:
	go run -mod $(GOMOD) -tags $(TAGS) cmd/retrieve-actor/main.go \
		-address $(ADDRESS) \
		-verbose \
		-insecure
