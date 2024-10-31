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

Add a new ActivityPub account.

```
$> ./bin/add-account -h

Add a new ActivityPub account.
Usage:
	 ./bin/add-account [options]
Valid options are:
  -account-icon-uri gocloud.dev/blob
    	A valid gocloud.dev/blob URI (as in the bucket URI + filename) referencing the icon URI for the account.
  -account-id int
    	An optional unique identifier to assign to the account being created. If 0 then an ID will be generated automatically.
  -account-name string
    	The user (preferred) name for the account being created.
  -account-type string
    	The type of account being created. Valid options are: Person, Service. (default "Person")
  -accounts-database-uri string
    	A valid sfomuseum/go-activitypub/database.AccountsDatabase URI.
  -alias value
    	Zero or more aliases for the account being created.
  -aliases-database-uri string
    	A valid sfomuseum/go-activitypub/database.AliasesDatabase URI.
  -allow-remote-icon-uri
    	Allow the -account-icon-uri flag to specify a remote URI.
  -blurb string
    	The descriptive blurb (caption) for the account being created.
  -discoverable
    	Boolean flag indicating whether the account should be discoverable. (default true)
  -display-name string
    	The display name for the account being created.
  -embed-icon-uri
    	If true then assume the -account-icon-uri flag references a local file and read its body in to a base64-encoded value to be stored with the account record.
  -private-key-uri gocloud.dev/runtimevar
    	A valid gocloud.dev/runtimevar referencing the PEM-encoded private key for the account.
  -properties-database-uri string
    	A valid sfomuseum/go-activitypub/database.PropertiesDatabase URI.
  -property value
    	Zero or more {KEY}={VALUE} properties to be assigned to the new account.
  -public-key-uri gocloud.dev/runtimevar
    	A valid gocloud.dev/runtimevar referencing the PEM-encoded public key for the account.
  -url string
    	The URL for the account being created.
```

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