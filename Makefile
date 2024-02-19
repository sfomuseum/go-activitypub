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

db:
	rm -f *.db
	$(SQLITE3) $(ACCOUNTS_DB) < schema/sqlite/accounts.schema
	$(SQLITE3) $(FOLLOWERS_DB) < schema/sqlite/followers.schema
	$(SQLITE3) $(FOLLOWING_DB) < schema/sqlite/following.schema
	$(SQLITE3) $(POSTS_DB) < schema/sqlite/posts.schema
	$(SQLITE3) $(NOTES_DB) < schema/sqlite/notes.schema
	$(SQLITE3) $(MESSAGES_DB) < schema/sqlite/messages.schema
	$(SQLITE3) $(BLOCKS_DB) < schema/sqlite/blocks.schema

accounts:
	go run cmd/add-account/main.go -accounts-database-uri 'sql://sqlite3?dsn=$(ACCOUNTS_DB)' -account-name bob
	go run cmd/add-account/main.go -accounts-database-uri 'sql://sqlite3?dsn=$(ACCOUNTS_DB)' -account-name alice

# Bob wants to follow Alice

follow:
	go run cmd/follow/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=$(ACCOUNTS_DB)' \
		-following-database-uri 'sql://sqlite3?dsn=$(FOLLOWING_DB)' \
		-account-name bob \
		-follow alice@localhost:8080 

# Bob wants to unfollow Alice

unfollow:
	go run cmd/follow/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=$(ACCOUNTS_DB)' \
		-following-database-uri 'sql://sqlite3?dsn=$(FOLLOWING_DB)' \
		-account-name bob \
		-follow alice@localhost:8080 \
		-undo

block:
	go run cmd/block/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=$(ACCOUNTS_DB)' \
		-blocks-database-uri 'sql://sqlite3?dsn=$(BLOCKS_DB)' \
		-account-name bob \
		-block-host block.club

unblock:
	go run cmd/block/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=$(ACCOUNTS_DB)' \
		-blocks-database-uri 'sql://sqlite3?dsn=$(BLOCKS_DB)' \
		-account-name bob \
		-block-host block.club \
		-undo

# Alice wants to post something (to Bob, if Bob is following Alice)

post:
	go run cmd/post/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=$(ACCOUNTS_DB)' \
		-followers-database-uri 'sql://sqlite3?dsn=$(FOLLOWERS_DB)' \
		-posts-database-uri 'sql://sqlite3?dsn=$(POSTS_DB)' \
		-account-name alice \
		-message "$(MESSAGE)"

inbox:
	go run cmd/inbox/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=$(ACCOUNTS_DB)' \
		-messages-database-uri 'sql://sqlite3?dsn=$(MESSAGES_DB)' \
		-notes-database-uri 'sql://sqlite3?dsn=$(NOTES_DB)' \
		-account-name $(ACCOUNT)

server:
	go run cmd/server/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=$(ACCOUNTS_DB)' \
		-followers-database-uri 'sql://sqlite3?dsn=$(FOLLOWERS_DB)' \
		-following-database-uri 'sql://sqlite3?dsn=$(FOLLOWING_DB)' \
		-notes-database-uri 'sql://sqlite3?dsn=$(NOTES_DB)' \
		-messages-database-uri 'sql://sqlite3?dsn=$(MESSAGES_DB)' \
		-blocks-database-uri 'sql://sqlite3?dsn=$(BLOCKS_DB)' \
		-allow-create \
		-hostname localhost:8080
