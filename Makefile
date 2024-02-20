GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

LDFLAGS=-s -w

SQLITE3=sqlite3

ACCOUNTS_DB=accounts.db
FOLLOWERS_DB=followers.db
FOLLOWING_DB=following.db
POSTS_DB=posts.db
NOTES_DB=notes.db
MESSAGES_DB=messages.db
BLOCKS_DB=blocks.db

ACCOUNTS_DB_URI=sql://sqlite3?dsn=$(ACCOUNTS_DB)
FOLLOWERS_DB_URI=sql://sqlite3?dsn=$(FOLLOWERS_DB)
FOLLOWING_DB_URI=sql://sqlite3?dsn=$(FOLLOWING_DB)
BLOCKS_DB_URI=sql://sqlite3?dsn=$(BLOCKS_DB)
POSTS_DB_URI=sql://sqlite3?dsn=$(POSTS_DB)
NOTES_DB_URI=sql://sqlite3?dsn=$(NOTES_DB)
MESSAGES_DB_URI=sql://sqlite3?dsn=$(MESSAGES_DB)

ACCOUNTS_DB_URI=awsdynamodb://accounts?partition_key=Id&allow_scans=true&local=true
FOLLOWING_DB_URI=awsdynamodb://following?partition_key=Id&allow_scans=true&local=true
FOLLOWERS_DB_URI=awsdynamodb://followers?partition_key=Id&allow_scans=true&local=true
BLOCKS_DB_URI=awsdynamodb://blocks?partition_key=Id&allow_scans=true&local=true
NOTES_DB_URI=awsdynamodb://notes?partition_key=Id&allow_scans=true&local=true
POSTS_DB_URI=awsdynamodb://posts?partition_key=Id&allow_scans=true&local=true
MESSAGES_DB_URI=awsdynamodb://messages?partition_key=Id&allow_scans=true&local=true

db-sqlite:
	rm -f *.db
	$(SQLITE3) $(ACCOUNTS_DB) < schema/sqlite/accounts.schema
	$(SQLITE3) $(FOLLOWERS_DB) < schema/sqlite/followers.schema
	$(SQLITE3) $(FOLLOWING_DB) < schema/sqlite/following.schema
	$(SQLITE3) $(POSTS_DB) < schema/sqlite/posts.schema
	$(SQLITE3) $(NOTES_DB) < schema/sqlite/notes.schema
	$(SQLITE3) $(MESSAGES_DB) < schema/sqlite/messages.schema
	$(SQLITE3) $(BLOCKS_DB) < schema/sqlite/blocks.schema

accounts:
	go run cmd/add-account/main.go -accounts-database-uri '$(ACCOUNTS_DB_URI)' -account-name bob
	go run cmd/add-account/main.go -accounts-database-uri '$(ACCOUNTS_DB_URI)' -account-name alice

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
		-account-name alice \
		-message "$(MESSAGE)" \
		-hostname localhost:8080 \
		-insecure

inbox:
	go run cmd/inbox/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-messages-database-uri '$(MESSAGES_DB_URI)' \
		-notes-database-uri '$(NOTES_DB_URI)' \
		-account-name $(ACCOUNT)

server:
	go run cmd/server/main.go \
		-accounts-database-uri '$(ACCOUNTS_DB_URI)' \
		-followers-database-uri '$(FOLLOWERS_DB_URI)' \
		-following-database-uri '$(FOLLOWING_DB_URI)' \
		-notes-database-uri '$(NOTES_DB_URI)' \
		-messages-database-uri '$(MESSAGES_DB_URI)' \
		-blocks-database-uri '$(BLOCKS_DB_URI)' \
		-allow-create \
		-verbose \
		-hostname localhost:8080 \
		-insecure

# https://aws.amazon.com/about-aws/whats-new/2018/08/use-amazon-dynamodb-local-more-easily-with-the-new-docker-image/
# https://hub.docker.com/r/amazon/dynamodb-local/

dynamo-local:
	docker run --rm -it -p 8000:8000 amazon/dynamodb-local

dynamo-tables-local:
	go run -mod vendor cmd/create-dynamodb-tables/main.go \
		-refresh \
		-dynamodb-client-uri 'awsdynamodb://?local=true'
