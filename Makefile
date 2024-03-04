GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/server cmd/server/main.go

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
NOTES_DB=notes.db
MESSAGES_DB=messages.db
BLOCKS_DB=blocks.db
DELIVERIES_DB=deliveries.db

ACCOUNTS_DB_URI=sql://sqlite3?dsn=$(ACCOUNTS_DB)
FOLLOWERS_DB_URI=sql://sqlite3?dsn=$(FOLLOWERS_DB)
FOLLOWING_DB_URI=sql://sqlite3?dsn=$(FOLLOWING_DB)
BLOCKS_DB_URI=sql://sqlite3?dsn=$(BLOCKS_DB)
POSTS_DB_URI=sql://sqlite3?dsn=$(POSTS_DB)
NOTES_DB_URI=sql://sqlite3?dsn=$(NOTES_DB)
MESSAGES_DB_URI=sql://sqlite3?dsn=$(MESSAGES_DB)
DELIVERIES_DB_URI=sql://sqlite3?dsn=$(DELIVERIES_DB)

ACCOUNTS_DB_URI=awsdynamodb://$(TABLE_PREFIX)accounts?partition_key=Id&allow_scans=true&local=true
FOLLOWING_DB_URI=awsdynamodb://$(TABLE_PREFIX)following?partition_key=Id&allow_scans=true&local=true
FOLLOWERS_DB_URI=awsdynamodb://$(TABLE_PREFIX)followers?partition_key=Id&allow_scans=true&local=true
BLOCKS_DB_URI=awsdynamodb://$(TABLE_PREFIX)blocks?partition_key=Id&allow_scans=true&local=true
NOTES_DB_URI=awsdynamodb://$(TABLE_PREFIX)notes?partition_key=Id&allow_scans=true&local=true
POSTS_DB_URI=awsdynamodb://$(TABLE_PREFIX)posts?partition_key=Id&allow_scans=true&local=true
MESSAGES_DB_URI=awsdynamodb://$(TABLE_PREFIX)messages?partition_key=Id&allow_scans=true&local=true
DELIVERIES_DB_URI=awsdynamodb://$(TABLE_PREFIX)deliveries?partition_key=Id&allow_scans=true&local=true
ALIASES_DB_URI=awsdynamodb://$(TABLE_PREFIX)aliases?partition_key=Name&allow_scans=true&local=true

db-sqlite:
	rm -f *.db
	$(SQLITE3) $(ACCOUNTS_DB) < schema/sqlite/accounts.schema
	$(SQLITE3) $(FOLLOWERS_DB) < schema/sqlite/followers.schema
	$(SQLITE3) $(FOLLOWING_DB) < schema/sqlite/following.schema
	$(SQLITE3) $(POSTS_DB) < schema/sqlite/posts.schema
	$(SQLITE3) $(NOTES_DB) < schema/sqlite/notes.schema
	$(SQLITE3) $(MESSAGES_DB) < schema/sqlite/messages.schema
	$(SQLITE3) $(BLOCKS_DB) < schema/sqlite/blocks.schema
	$(SQLITE3) $(DELIVERIES_DB) < schema/sqlite/deliveries.schema

accounts:
	go run cmd/add-account/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-aliases-database-uri '$(ALIASES_DB_URI)' \
		-account-name bob \
		-alias robert \
		-account-type Service \
		-account-icon-uri fixtures/icons/bob.jpg \
		-embed-icon-uri
	go run cmd/add-account/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-aliases-database-uri '$(ALIASES_DB_URI)' \
		-account-name alice \
		-account-type Person \
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
	go run cmd/post/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-followers-database-uri '$(FOLLOWERS_DB_URI)' \
		-posts-database-uri '$(POSTS_DB_URI)' \
		-deliveries-database-uri '$(DELIVERIES_DB_URI)' \
		-account-name alice \
		-message "$(MESSAGE)" \
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
