accounts:
	go run cmd/add-actor/main.go -accounts-database-uri 'sql://sqlite3?dsn=test.db' -account-id bob
	go run cmd/add-actor/main.go -accounts-database-uri 'sql://sqlite3?dsn=test.db' -account-id alice

# Bob wants to follow Alice

follow:
	go run cmd/follow/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=test.db' \
		-account-id bob \
		-follow alice@localhost:8080 

# Bob wants to unfollow Alice

unfollow:
	go run cmd/follow/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=test.db' \
		-account-id bob \
		-follow alice@localhost:8080 
		-undo

# Alice wants to post something (to Bob, if Bob is following Alice)

post:
	go run cmd/post/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=test.db' \
		-followers-database-uri 'sql://sqlite3?dsn=test.db' \
		-posts-database-uri 'sql://sqlite3?dsn=test.db' \
		-account-id alice

server:
	go run cmd/server/main.go \
		-accounts-database-uri 'sql://sqlite3?dsn=test.db' \
		-followers-database-uri 'sql://sqlite3?dsn=test.db' \
		-hostname localhost:8080
