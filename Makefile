# $> urlescape 'test.db?cache=shared'
DSN=test.db%3Fcache%3Dshared

accounts:
	go run cmd/add-actor/main.go -accounts-database-uri 'sql://sqlite3?dsn=$(DSN)' -account-id bob
	go run cmd/add-actor/main.go -accounts-database-uri 'sql://sqlite3?dsn=$(DSN)' -account-id alice

# Bob wants to follow Alice

follow:
	go run cmd/follow/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=$(DSN)' \
		-following-database-uri 'sql://sqlite3?dsn=$(DSN)' \
		-account-id bob \
		-follow alice@localhost:8080 

# Bob wants to unfollow Alice

unfollow:
	go run cmd/follow/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=$(DSN)' \
		-following-database-uri 'sql://sqlite3?dsn=$(DSN)' \
		-account-id bob \
		-follow alice@localhost:8080 \
		-undo

# Alice wants to post something (to Bob, if Bob is following Alice)

post:
	go run cmd/post/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=$(DSN)' \
		-followers-database-uri 'sql://sqlite3?dsn=$(DSN)' \
		-posts-database-uri 'sql://sqlite3?dsn=$(DSN)' \
		-account-id alice

server:
	go run cmd/server/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=$(DSN)' \
		-followers-database-uri 'sql://sqlite3?dsn=$(DSN)' \
		-following-database-uri 'sql://sqlite3?dsn=$(DSN)' \
		-notes-database-uri 'sql://sqlite3?dsn=$(DSN)' \
		-hostname localhost:8080
