# Command line tools

The easiest way to build the tools provided by this package is to run the `cli` Makefile target. For example

```
$> make cli
cd ../ && make cli && cd -
go build -mod vendor -ldflags="-s -w" -o bin/server cmd/server/main.go
go build -mod vendor -ldflags="-s -w" -o bin/add-account cmd/add-account/main.go
go build -mod vendor -ldflags="-s -w" -o bin/get-account cmd/get-account/main.go
go build -mod vendor -ldflags="-s -w" -o bin/create-post cmd/create-post/main.go
go build -mod vendor -ldflags="-s -w" -o bin/deliver-activity cmd/deliver-activity/main.go
go build -mod vendor -ldflags="-s -w" -o bin/list-followers cmd/list-followers/main.go
go build -mod vendor -ldflags="-s -w" -o bin/list-addresses cmd/list-addresses/main.go
go build -mod vendor -ldflags="-s -w" -o bin/counts-for-date cmd/counts-for-date/main.go
go build -mod vendor -ldflags="-s -w" -o bin/inbox cmd/inbox/main.go
go build -mod vendor -ldflags="-s -w" -o bin/create-dynamodb-tables cmd/create-dynamodb-tables/main.go
```

## Tools

### add-account

### add-aliases

### block

### boost-note

### counts-for-date

### create-dynamodb-tables

### create-icon

### create-post

### deliver-activity

### follow

### get-account

### get-note

### inbox

### list-activities

### list-addresses

### list-aliases

### list-boosts

### list-deliveries

### list-followers

### parse-activity

### post-from-uri

### retrieve-actor

### retrieve-delivery

### retrieve-note

Retrieve a given note from the Notes database.

```
$> ./bin/retrieve-note -h
  -body
    	Display the (ActivityPub) body of the note.
  -note-id int
    	The unique 64-bit note ID to retrieve.
  -notes-database-uri string
    	A valid sfomuseum/go-activitypub/database.NotesDatabase URI.
  -verbose
    	Enable verbose (debug) logging.
```

### server