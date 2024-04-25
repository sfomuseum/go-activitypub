GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/server cmd/server/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/add-account cmd/add-account/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/create-post cmd/create-post/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/list-followers cmd/list-followers/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/list-addresses cmd/list-addresses/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/counts-for-date cmd/counts-for-date/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/inbox cmd/inbox/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/create-dynamodb-tables cmd/create-dynamodb-tables/main.go

lambda:
	@make lambda-server
	@make lambda-deliver-post

lambda-server:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags lambda.norpc -o bootstrap cmd/server/main.go
	zip server.zip bootstrap
	rm -f bootstrap

lambda-deliver-post:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f deliver.zip; then rm -f deliver.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags lambda.norpc -o bootstrap cmd/deliver-post/main.go
	zip deliver.zip bootstrap
	rm -f bootstrap

# The rest of these Makefile targets are for local testing

SQLITE3=sqlite3
TABLE_PREFIX=

ACCOUNTS_DB=accounts.db
FOLLOWERS_DB=followers.db
FOLLOWING_DB=following.db
POSTS_DB=posts.db
POST_TAGS_DB=posts.db
NOTES_DB=notes.db
MESSAGES_DB=messages.db
BLOCKS_DB=blocks.db
DELIVERIES_DB=deliveries.db
BOOSTS_DB=boosts.db
LIKES_DB=likes.db
PROPERTIES_DB=properties.db

ACCOUNTS_DB_URI=sql://sqlite3?dsn=$(ACCOUNTS_DB)
FOLLOWERS_DB_URI=sql://sqlite3?dsn=$(FOLLOWERS_DB)
FOLLOWING_DB_URI=sql://sqlite3?dsn=$(FOLLOWING_DB)
BLOCKS_DB_URI=sql://sqlite3?dsn=$(BLOCKS_DB)
POSTS_DB_URI=sql://sqlite3?dsn=$(POSTS_DB)
POST_TAGS_DB_URI=sql://sqlite3?dsn=$(POST_TAGS_DB)
NOTES_DB_URI=sql://sqlite3?dsn=$(NOTES_DB)
MESSAGES_DB_URI=sql://sqlite3?dsn=$(MESSAGES_DB)
DELIVERIES_DB_URI=sql://sqlite3?dsn=$(DELIVERIES_DB)
BOOSTS_DB_URI=sql://sqlite3?dsn=$(BOOSTS_DB)
LIKES_DB_URI=sql://sqlite3?dsn=$(LIKES_DB)
PROPERTIES_DB_URI=sql://sqlite3?dsn=$(PROPERTIES_DB)


ACCOUNTS_DB_URI=awsdynamodb://$(TABLE_PREFIX)accounts?partition_key=Id&allow_scans=true&local=true
ALIASES_DB_URI=awsdynamodb://$(TABLE_PREFIX)aliases?partition_key=Name&allow_scans=true&local=true
BLOCKS_DB_URI=awsdynamodb://$(TABLE_PREFIX)blocks?partition_key=Id&allow_scans=true&local=true
BOOSTS_DB_URI=awsdynamodb://$(TABLE_PREFIX)boosts?partition_key=Id&allow_scans=true&local=true
DELIVERIES_DB_URI=awsdynamodb://$(TABLE_PREFIX)deliveries?partition_key=Id&allow_scans=true&local=true
FOLLOWING_DB_URI=awsdynamodb://$(TABLE_PREFIX)following?partition_key=Id&allow_scans=true&local=true
FOLLOWERS_DB_URI=awsdynamodb://$(TABLE_PREFIX)followers?partition_key=Id&allow_scans=true&local=true
LIKES_DB_URI=awsdynamodb://$(TABLE_PREFIX)likes?partition_key=Id&allow_scans=true&local=true
NOTES_DB_URI=awsdynamodb://$(TABLE_PREFIX)notes?partition_key=Id&allow_scans=true&local=true
MESSAGES_DB_URI=awsdynamodb://$(TABLE_PREFIX)messages?partition_key=Id&allow_scans=true&local=true
POST_TAGS_DB_URI=awsdynamodb://$(TABLE_PREFIX)post_tags?partition_key=Id&allow_scans=true&local=true
POSTS_DB_URI=awsdynamodb://$(TABLE_PREFIX)posts?partition_key=Id&allow_scans=true&local=true
PROPERTIES_DB_URI=awsdynamodb://$(TABLE_PREFIX)properties?partition_key=Id&allow_scans=true&local=true

db-sqlite:
	rm -f *.db
	$(SQLITE3) $(ACCOUNTS_DB) < schema/sqlite/accounts.schema
	$(SQLITE3) $(FOLLOWERS_DB) < schema/sqlite/followers.schema
	$(SQLITE3) $(FOLLOWING_DB) < schema/sqlite/following.schema
	$(SQLITE3) $(POSTS_DB) < schema/sqlite/posts.schema
	$(SQLITE3) $(POST_TAGS_DB) < schema/sqlite/post_tags.schema
	$(SQLITE3) $(NOTES_DB) < schema/sqlite/notes.schema
	$(SQLITE3) $(MESSAGES_DB) < schema/sqlite/messages.schema
	$(SQLITE3) $(BLOCKS_DB) < schema/sqlite/blocks.schema
	$(SQLITE3) $(DELIVERIES_DB) < schema/sqlite/deliveries.schema
	$(SQLITE3) $(BOOSTS_DB) < schema/sqlite/boosts.schema
	$(SQLITE3) $(LIKES_DB) < schema/sqlite/likes.schema
	$(SQLITE3) $(PROPERTIES_DB) < schema/sqlite/properties.schema

accounts:
	go run cmd/add-account/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-aliases-database-uri '$(ALIASES_DB_URI)' \
		-properties-database-uri '$(PROPERTIES_DB_URI)' \
		-account-name bob \
		-alias robert \
		-account-type Service \
		-account-icon-uri fixtures/icons/bob.jpg \
		-property 'url:www=https://bob.com' \
		-embed-icon-uri
	go run cmd/add-account/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-aliases-database-uri '$(ALIASES_DB_URI)' \
		-properties-database-uri '$(PROPERTIES_DB_URI)' \
		-account-name doug \
		-alias doug \
		-property 'url:www=https://bob.com/doug' \
		-account-type Service
	go run cmd/add-account/main.go \
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
	go run cmd/follow/main.go \
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
	go run cmd/follow/main.go \
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
	go run cmd/block/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-blocks-database-uri '$(BLOCKS_DB_URI)' \
		-account-name bob \
		-block-host block.club

unblock:
	go run cmd/block/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-blocks-database-uri '$(BLOCKS_DB_URI)' \
		-account-name bob \
		-block-host block.club \
		-undo

# Alice wants to post something (to Bob, if Bob is following Alice)

post:
	go run cmd/create-post/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-followers-database-uri '$(FOLLOWERS_DB_URI)' \
		-posts-database-uri '$(POSTS_DB_URI)' \
		-post-tags-database-uri '$(POST_TAGS_DB_URI)' \
		-deliveries-database-uri '$(DELIVERIES_DB_URI)' \
		-account-name alice \
		-message "$(MESSAGE)" \
		-hostname localhost:8080 \
		-insecure \
		-verbose

# -mention $(MENTION) \

reply:
	go run cmd/create-post/main.go \
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
	go run cmd/retrieve-delivery/main.go \
		-deliveries-database-uri '$(DELIVERIES_DB_URI)' \
		-delivery-id $(ID) \
		-verbose

inbox:
	go run cmd/inbox/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-messages-database-uri '$(MESSAGES_DB_URI)' \
		-notes-database-uri '$(NOTES_DB_URI)' \
		-account-name $(ACCOUNT)

server:
	go run cmd/server/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-aliases-database-uri '$(ALIASES_DB_URI)' \
		-followers-database-uri '$(FOLLOWERS_DB_URI)' \
		-following-database-uri '$(FOLLOWING_DB_URI)' \
		-notes-database-uri '$(NOTES_DB_URI)' \
		-messages-database-uri '$(MESSAGES_DB_URI)' \
		-blocks-database-uri '$(BLOCKS_DB_URI)' \
		-posts-database-uri '$(POSTS_DB_URI)' \
		-boosts-database-uri '$(BOOSTS_DB_URI)' \
		-likes-database-uri '$(LIKES_DB_URI)' \
		-properties-database-uri '$(PROPERTIES_DB_URI)' \
		-allow-remote-icon-uri \
		-allow-create \
		-verbose \
		-hostname localhost:8080 \
		-insecure

retrieve:
	go run cmd/retrieve-actor/main.go \
		-address $(ADDRESS) \
		-verbose \
		-insecure

# https://aws.amazon.com/about-aws/whats-new/2018/08/use-amazon-dynamodb-local-more-easily-with-the-new-docker-image/
# https://hub.docker.com/r/amazon/dynamodb-local/

dynamo-local:
	docker run --rm -it -p 8000:8000 amazon/dynamodb-local

dynamo-tables-local:
	go run -mod vendor cmd/create-dynamodb-tables/main.go \
		-refresh \
		-table-prefix '$(TABLE_PREFIX)' \
		-dynamodb-client-uri 'awsdynamodb://?local=true'

# I haven't been able to get this to work yet...
# https://dev.mysql.com/doc/mysql-installation-excerpt/8.3/en/docker-mysql-getting-started.html#docker-starting-mysql-server

mysql-local:
	docker run --rm -it -p3306:3306 container-registry.oracle.com/mysql/community-server:latest
