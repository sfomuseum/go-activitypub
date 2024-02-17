follow:
	go run cmd/follow/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=test.db' \
		-account-id bob@localhost:8080 \
		-follow http://localhost:8080/profile/alice \
		-inbox http://localhost:8080/inbox/alice

unfollow:
	go run cmd/follow/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=test.db' \
		-account-id bob@localhost:8080 \
		-follow http://localhost:8080/profile/alice \
		-inbox http://localhost:8080/inbox/alice \
		-undo

post:
	go run cmd/post/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=test.db' \
		-followers-database-uri 'sql://sqlite3?dsn=test.db' \
		-posts-database-uri 'sql://sqlite3?dsn=test.db' \
		-account-id alice@localhost:8080

server:
	go run cmd/server/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=test.db' \
		-followers-database-uri 'sql://sqlite3?dsn=test.db' \
		-hostname localhost:8080
