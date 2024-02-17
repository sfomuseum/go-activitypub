# go-activitypub

```
$> sqlite3 test.db < schema/accounts.sqlite.schema
```

```
$> go run cmd/add-actor/main.go -account-database-uri 'sql://sqlite3?dsn=test.db' -account-id bob@localhost:8080
$> go run cmd/add-actor/main.go -account-database-uri 'sql://sqlite3?dsn=test.db' -account-id alice@localhost:8080
```

```
$> go run cmd/server/main.go -accounts-database-uri 'sql://sqlite3?dsn=test.db' -hostname localhost:8080
```

```
$> go run cmd/follow/main.go \
	-accounts-database-uri 'sql://sqlite3?dsn=test.db' \
	-account-id bob@localhost:8080 \
	-follow http://localhost:8080/profile/alice \
	-inbox http://localhost:8080/inbox/alice
```

## See also

* https://github.com/w3c/activitypub/blob/gh-pages/activitypub-tutorial.txt