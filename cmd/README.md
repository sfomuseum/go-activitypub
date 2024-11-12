# Command line tools

The easiest way to build the tools provided by this package is to run the `cli` Makefile target. For example

```
$> make cli
cd ../ && make cli && cd -
go build -mod vendor -ldflags="-s -w" -o bin/add-account cmd/add-account/main.go
go build -mod vendor -ldflags="-s -w" -o bin/add-aliases cmd/add-aliases/main.go
go build -mod vendor -ldflags="-s -w" -o bin/block cmd/block/main.go
go build -mod vendor -ldflags="-s -w" -o bin/boost-note cmd/boost-note/main.go
go build -mod vendor -ldflags="-s -w" -o bin/counts-for-date cmd/counts-for-date/main.go
go build -mod vendor -ldflags="-s -w" -o bin/create-dynamodb-tables cmd/create-dynamodb-tables/main.go
go build -mod vendor -ldflags="-s -w" -o bin/create-post cmd/create-post/main.go
go build -mod vendor -ldflags="-s -w" -o bin/deliver-activity cmd/deliver-activity/main.go
go build -mod vendor -ldflags="-s -w" -o bin/get-account cmd/get-account/main.go
go build -mod vendor -ldflags="-s -w" -o bin/follow cmd/follow/main.go
go build -mod vendor -ldflags="-s -w" -o bin/list-boosts cmd/list-boosts/main.go
go build -mod vendor -ldflags="-s -w" -o bin/list-followers cmd/list-followers/main.go
go build -mod vendor -ldflags="-s -w" -o bin/list-activities cmd/list-activities/main.go
go build -mod vendor -ldflags="-s -w" -o bin/list-addresses cmd/list-addresses/main.go
go build -mod vendor -ldflags="-s -w" -o bin/list-aliases cmd/list-aliases/main.go
go build -mod vendor -ldflags="-s -w" -o bin/list-deliveries cmd/list-deliveries/main.go
go build -mod vendor -ldflags="-s -w" -o bin/inbox cmd/inbox/main.go
go build -mod vendor -ldflags="-s -w" -o bin/retrieve-actor cmd/retrieve-actor/main.go
go build -mod vendor -ldflags="-s -w" -o bin/retrieve-delivery cmd/retrieve-delivery/main.go
go build -mod vendor -ldflags="-s -w" -o bin/retrieve-note cmd/retrieve-note/main.go
go build -mod vendor -ldflags="-s -w" -o bin/server cmd/server/main.go
```

## Structure

Generally speaking the code in the `cmd` folder are thin wrappers to invoke application code located in the `app` package. This is to make it easier to individual tools to be supplemented with custom code, or additional flags, on an as-needed basis.

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

Add aliases for a registered sfomuseum/go-activity account.

```
$> ./bin/add-aliases -h
Add aliases for a registered sfomuseum/go-activity account.
Usage:
	 ./bin/add-aliases [options]
Valid options are:
  -account-name string
    	A valid sfomuseum/go-activitypub account name
  -accounts-database-uri string
    	A registered sfomuseum/go-activitypub/AccountsDatabase URI. (default "null://")
  -alias value
    	One or more aliases to add for an account. Each -alias flag may be a CSV-encoded string containing multiple aliases.
  -aliases-database-uri string
    	A registered sfomuseum/go-activitypub/AliasesDatabase URI. (default "null://")
```

### block

Manage the blocking of third-parties on behalf of a registered sfomuseum/go-activity account.

```
$> ./bin/block -h
Manage the blocking of third-parties on behalf of a registered sfomuseum/go-activity account.
Usage:
	 ./bin/block [options]
Valid options are:
  -account-name string
    	The name of the account doing the blocking.
  -accounts-database-uri string
    	A known sfomuseum/go-activitypub/AccountsDatabase URI.
  -block-host string
    	The name of the host associated with the account being blocked.
  -block-name string
    	The name of the account being blocked. If "*" then all the accounts associated with the blocked host will be blocked. (default "*")
  -blocks-database-uri string
    	A known sfomuseum/go-activitypub/BlocksDatabase URI.
  -undo
    	Undo an existing block.
  -verbose
    	Enable verbose (debug) logging.
```

### boost-note

Boost an ActivityPub note on behalf of a registered go-activity account.

```
$> ./bin/boost-note -h
Boost an ActivityPub note on behalf of a registered go-activity account.
Usage:
	 ./bin/boost-note [options]
Valid options are:
  -account-name string
    	The account doing the boosting.
  -accounts-database-uri string
    	A known sfomuseum/go-activitypub/AccountsDatabase URI.
  -activities-database-uri string
    	A known sfomuseum/go-activitypub/ActivitiesDatabase URI.
  -deliveries-database-uri string
    	A known sfomuseum/go-activitypub/DeliveriesDatabase URI.
  -delivery-queue-uri string
    	A known sfomuseum/go-activitypub/queue.DeliveryQueue URI. (default "synchronous://")
  -followers-database-uri string
    	A known sfomuseum/go-activitypub/FollowersDatabase URI.
  -hostname string
    	The hostname (domain) for the account doing the boosting. (default "localhost:8080")
  -insecure
    	A boolean flag indicating the ActivityPub server delivering activities is insecure (not using TLS).
  -note string
    	The URI of the note being boosted.
  -verbose
    	Enable verbose (debug) logging.
```

### counts-for-date

### create-dynamodb-tables

Create (or refresh) the DynamoDB tables necessary for use with the go-activitypub package.

```
$> ./bin/create-dynamodb-tables -h
Create (or refresh) the DynamoDB tables necessary for use with the go-activitypub package.
Usage:
	 ./bin/create-dynamodb-tables [options]
Valid options are:
  -dynamodb-client-uri string
    	A valid aaronland/gocloud-docstore URI (dynamodb:// or awsdynamodb://).
  -refresh
    	Refresh tables if already present.
  -table value
    	Zero or more table names to create. If zero then all the default tables will be created.
  -table-prefix string
    	A optional prefix to assign to each table name.
```

### create-post

```
$> ./bin/create-post -h
Create a new post (note, activity) on behalf of a registered go-activitypub account and schedule it for delivery to all their followers.
Usage:
	 ./bin/create-post [options]
Valid options are:
  -account-name string
    	The name of the go-activitypub account creating the post.
  -accounts-database-uri string
    	A registered sfomuseum/go-activitypub/database.AccountsDatabase URI.
  -activities-database-uri string
    	A registered sfomuseum/go-activitypub/database.ActivitiesDatabase URI.
  -deliveries-database-uri string
    	A registered sfomuseum/go-activitypub/database.DeliveriesDatabase URI.
  -delivery-queue-uri string
    	A registered sfomuseum/go-activitypub/queue/DeliveryQueue URI. (default "synchronous://")
  -followers-database-uri string
    	A registered sfomuseum/go-activitypub/database.FollowersDatabase URI.
  -hostname string
    	The hostname (domain) of the ActivityPub server delivering activities. (default "localhost:8080")
  -in-reply-to string
    	The URI of that the post is in reply to (optional).
  -insecure
    	A boolean flag indicating the ActivityPub server delivering activities is insecure (not using TLS).
  -max-attempts int
    	The maximum number of attempts to deliver the activity. (default 5)
  -message string
    	The body (content) of the message to post.
  -post-tags-database-uri string
    	A registered sfomuseum/go-activitypub/database.PostTagsDatabase URI.
  -posts-database-uri string
    	A registered sfomuseum/go-activitypub/database.PostsDatabase URI.
  -verbose
    	Enable verbose (debug) logging.
```

### deliver-activity

Deliver an ActivityPub activity to subscribers.

```
$> ./bin/deliver-activity -h
Deliver an ActivityPub activity to subscribers.
Usage:
	 ./bin/deliver-activity [options]
Valid options are:
  -accounts-database-uri string
    	A registered sfomuseum/go-activitypub/database.AccountsDatabase URI.
  -activities-database-uri string
    	A registered sfomuseum/go-activitypub/database.ActivitiesDatabase URI.
  -allow-mentions
    	Enable support for processing mentions in (post) activities. This enabled posts to accounts not followed by author but where account is mentioned in post. (default true)
  -deliveries-database-uri string
    	A registered sfomuseum/go-activitypub/database.DeliveriesDatabase URI.
  -delivery-queue-uri string
    	A registered sfomuseum/go-activitypub/queue.DeliveryQueue URI. (default "synchronous://")
  -followers-database-uri string
    	A registered sfomuseum/go-activitypub/database.FollowersDatabase URI.
  -hostname string
    	The hostname of the ActivityPub server delivering activities. (default "localhost:8080")
  -insecure
    	A boolean flag indicating the ActivityPub server delivering activities is insecure.
  -max-attempts int
    	The maximum number of attempts to deliver the activity. (default 5)
  -mode string
    	The operation mode for delivering activities. Valid options are: lambda, pubsub. "cli" mode is currently disabled.
  -post-tags-database-uri string
    	A registered sfomuseum/go-activitypub/database.PostTagsDatabase URI. (default "null://")
  -posts-database-uri string
    	A registered sfomuseum/go-activitypub/database.PostsDatabase URI.
  -subscriber-uri string
    	A valid sfomuseum/go-pubsub/subscriber URI. Required if -mode parameter is 'pubsub'.
  -verbose
    	Enable verbose logging
```

### follow

Follow a @user@host ActivityPub account on behalf of a registered go-activitypub account.

```
$> ./bin/follow -h
Follow a @user@host ActivityPub account on behalf of a registered go-activitypub account.
Usage:
	 ./bin/follow [options]
Valid options are:
  -account-name string
    	The name of the account doing the following.
  -accounts-database-uri string
    	A registered sfomuseum/go-activitypub/AccountsDatabase URI.
  -follow string
    	The ActivityPub @user@host address to follow.
  -following-database-uri string
    	A registered sfomuseum/go-activitypub/FollowingDatabase URI.
  -hostname string
    	The hostname (domain) of the ActivityPub server delivering activities. (default "localhost:8080")
  -insecure
    	A boolean flag indicating the ActivityPub server delivering activities is insecure (not using TLS).
  -messages-database-uri string
    	A registered sfomuseum/go-activitypub/MessagesDatabase URI.
  -undo
    	Stop following the account defined by the -follow flag.
  -verbose
    	Enable verbose (debug) logging.
```

### get-account

Retrieve an ActivityPub account and emit its details as a JSON-encoded string.

```
$> ./bin/get-account -h
Retrieve an ActivityPub account and emit its details as a JSON-encoded string.
Usage:
	 ./bin/get-account [options]
Valid options are:
  -account-name string
    	A valid sfomuseum/go-activitypub account name
  -accounts-database-uri string
    	A registered sfomuseum/go-activitypub/AccountsDatabase URI. (default "null://")
  -properties-database-uri string
    	A registered sfomuseum/go-activitypub/PropertiesDatabase URI (default "null://")
```

### inbox

```
$> ./bin/inbox -h
Display all the messges (notes) received for a registered go-activitypub account.
Usage:
	 ./bin/inbox [options]
Valid options are:
  -account-name string
    	The name of the account whose inbox you want to display.
  -accounts-database-uri string
    	A registered sfomuseum/go-activitypub/AccountsDatabase URI.
  -messages-database-uri string
    	A registered sfomuseum/go-activitypub/MessagesDatabase URI.
  -notes-database-uri string
    	A registered sfomuseum/go-activitypub/NotesDatabase URI.
  -verbose
    	Enable verbose (debug) logging.
```

### list-activities

List all the activities that have been created.

```
$> ./bin/list-activities -h
List all the activities that have been created.
Usage:
	 ./bin/list-activities [options]
Valid options are:
  -activities-database-uri string
    	A registered sfomuseum/go-activitypub/database.ActivitiesDatabase URI.
  -verbose
    	Enable verbose (debug) logging.
```

### list-aliases

List the aliases for a registered sfomuseum/go-activity account.

```
$> ./bin/list-aliases -h
List the aliases for a registered sfomuseum/go-activity account.
Usage:
	 ./bin/list-aliases [options]
Valid options are:
  -account-name string
    	A valid sfomuseum/go-activitypub account name
  -accounts-database-uri string
    	A known sfomuseum/go-activitypub/AccountsDatabase URI. (default "null://")
  -aliases-database-uri string
    	A known sfomuseum/go-activitypub/AliasesDatabase URI. (default "null://")
```

### list-boosts

List all the boosts received by a go-activitypub account.

```
$> ./bin/list-boosts -h
List all the boosts received by a go-activitypub account.
Usage:
	 ./bin/list-boosts [options]
Valid options are:
  -account-name string
    	The account whose posts have been boosted.
  -accounts-database-uri string
    	A known sfomuseum/go-activitypub/AccountsDatabase URI.
  -boosts-database-uri string
    	A known sfomuseum/go-activitypub/BlockDatabase URI.
  -hostname string
    	The hostname (domain) for the account doing the boosting. (default "localhost:8080")
  -insecure
    	A boolean flag indicating the ActivityPub server delivering activities is insecure (not using TLS).
  -verbose
    	Enable verbose (debug) logging.
```

### list-deliveries

List all the (ActivityPub activities) deliveries that have been recorded.

```
$> ./bin/list-deliveries -h
List all the (ActivityPub activities) deliveries that have been recorded.
Usage:
	 ./bin/list-deliveries [options]
Valid options are:
  -deliveries-database-uri string
    	A known sfomuseum/go-activitypub/DeliveriesDatabase URI.
  -verbose
    	Enable verbose (debug) logging.
```	

### retrieve-actor

Retrieve an ActivityPub actor by its @user@host address and emit it as a JSON-encoded string..

```
$> ./bin/retrieve-actor -h
Retrieve an ActivityPub actor by its @user@host address and emit it as a JSON-encoded string..
Usage:
	 ./bin/retrieve-actor [options]
Valid options are:
  -address string
    	The @user@host address of the actor to retrieve.
  -insecure
    	A boolean flag indicating whether the host that the -address flag resolves to is running without TLS enabled.
  -verbose
    	Enable verbose (debug) logging.
```

### retrieve-delivery

Retrieve and display a specific (ActivityPub activity) delivery.

```
$> ./bin/retrieve-delivery -h
Retrieve and display a specific (ActivityPub activity) delivery.
Usage:
	 ./bin/retrieve-delivery [options]
Valid options are:
  -deliveries-database-uri string
    	A registered sfomuseum/go-activitypub/DeliveriesDatabase URI.
  -delivery-id int
    	The unique ID of the delivery to retrieve.
  -verbose
    	Enable verbose (debug) logging.
```

### retrieve-note

Retrieve a message (note) by its unique go-activitypub ID.

```
$> ./bin/retrieve-note -h
Retrieve a message (note) by its unique go-activitypub ID.
Usage:
	 ./bin/retrieve-note [options]
Valid options are:
  -body
    	Display the (ActivityPub) body of the note.
  -note-id int
    	The unique 64-bit note ID to retrieve.
  -notes-database-uri string
    	A registered sfomuseum/go-activitypub/database.NotesDatabase URI.
  -verbose
    	Enable verbose (debug) logging.
```
 
### server

Start a HTTP (web) server to handle ActivityPub-related requests.

```
$> ./bin/server -h
Start a HTTP (web) server to handle ActivityPub-related requests.
Usage:
	 ./bin/server [options]
Valid options are:
  -accounts-database-uri string
    	A registered sfomuseum/go-activitypub/database.DeliveriesDatabase URI.
  -aliases-database-uri string
    	A registered sfomuseum/go-activitypub/database.AliasesDatabase URI.
  -allow-boosts
    	Enable support for ActivityPub "Announce" (boost) activities. (default true)
  -allow-create
    	Enable support for ActivityPub "Create" activities.
  -allow-follow
    	Enable support for ActivityPub "Follow" activities. (default true)
  -allow-likes
    	Enable support for ActivityPub "Like" activities. (default true)
  -allow-mentions
    	If enabled allows posts ("Create" activities) to accounts not followed by author but where account is mentioned in post. (default true)
  -allow-remote-icon-uri
    	Allow account icons hosted on a remote host.
  -blocks-database-uri string
    	A registered sfomuseum/go-activitypub/database.BlocksDatabase URI.
  -boosts-database-uri string
    	A registered sfomuseum/go-activitypub/database.BoostsDatabase URI.
  -disabled
    	Return a 503 Service unavailable response for all requests.
  -followers-database-uri string
    	A registered sfomuseum/go-activitypub/database.FollowersDatabase URI.
  -following-database-uri string
    	A registered sfomuseum/go-activitypub/database.FollowingDatabase URI.
  -hostname string
    	The hostname (domain) of the ActivityPub server delivering activities.
  -insecure
    	A boolean flag indicating the ActivityPub server delivering activities is insecure (not using TLS).
  -likes-database-uri string
    	A registered sfomuseum/go-activitypub/database.LikesDatabase URI.
  -messages-database-uri string
    	A registered sfomuseum/go-activitypub/database.MessagesDatabase URI.
  -notes-database-uri string
    	A registered sfomuseum/go-activitypub/database.NotesDatabase URI.
  -post-tags-database-uri string
    	A registered sfomuseum/go-activitypub/database.PostTagsDatabase URI.
  -posts-database-uri string
    	A registered sfomuseum/go-activitypub/database.PostsDatabase URI.
  -process-message-queue-uri string
    	A registered go-activitypub/queue.ProcessMessageQueue URI. (default "null://")
  -properties-database-uri string
    	A registered sfomuseum/go-activitypub/database.PropertiesDatabase URI.
  -server-uri string
    	A registered aaronland/go-http-server/server.Server URI. (default "http://localhost:8080")
  -verbose
    	Enable verbose (debug) logging.
```	