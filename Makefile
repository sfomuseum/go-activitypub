# $> urlescape 'test.db?cache=shared'
DSN=test.db%3Fcache%3Dshared%26mode%3Drwc

SQLITE3=sqlite3

ACCOUNTS_DB=accounts.db
FOLLOWERS_DB=followers.db
FOLLOWING_DB=following.db
POSTS_DB=posts.db
NOTES_DB=notes.db
MESSAGES_DB=messages.db

db:
	rm -f *.db
	$(SQLITE3) $(ACCOUNTS_DB) < schema/accounts.sqlite.schema
	$(SQLITE3) $(FOLLOWERS_DB) < schema/followers.sqlite.schema
	$(SQLITE3) $(FOLLOWING_DB) < schema/following.sqlite.schema
	$(SQLITE3) $(POSTS_DB) < schema/posts.sqlite.schema
	$(SQLITE3) $(NOTES_DB) < schema/notes.sqlite.schema
	$(SQLITE3) $(MESSAGES_DB) < schema/messages.sqlite.schema

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

# Alice wants to post something (to Bob, if Bob is following Alice)

post:
	go run cmd/post/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=$(ACCOUNTS_DB)' \
		-followers-database-uri 'sql://sqlite3?dsn=$(FOLLOWERS_DB)' \
		-posts-database-uri 'sql://sqlite3?dsn=$(POSTS_DB)' \
		-account-name alice \
		-message "$(MESSAGE)"

server:
	go run cmd/server/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=$(ACCOUNTS_DB)' \
		-followers-database-uri 'sql://sqlite3?dsn=$(FOLLOWERS_DB)' \
		-following-database-uri 'sql://sqlite3?dsn=$(FOLLOWING_DB)' \
		-notes-database-uri 'sql://sqlite3?dsn=$(NOTES_DB)' \
		-messages-database-uri 'sql://sqlite3?dsn=$(MESSAGES_DB)' \
		-hostname localhost:8080
