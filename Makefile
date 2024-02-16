server:
	go run cmd/server/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=test.db' \
		-followers-database-uri 'sql://sqlite3?dsn=test.db' \
		-hostname localhost:8080
