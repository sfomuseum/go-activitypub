# Databases

Databases in `go-activitypub` are really more akin to traditional RDBMS "tables" in that all of the database (or tables) listed below have one or more associations with each other. 

Each of these databases is defined as a Go-language interface with one or more default implementations provided by this package. For example:

```
type GetFollowersCallbackFunc func(context.Context, string) error

type FollowersDatabase interface {
	GetFollowersForAccount(context.Context, int64, GetFollowersCallbackFunc) error
	GetFollower(context.Context, int64, string) (*Follower, error)
	AddFollower(context.Context, *Follower) error
	RemoveFollower(context.Context, *Follower) error
	Close(context.Context) error
}
```

The idea here is that the various tools for performing actions (posting, serving ActivityPub requests, etc.) don't know anything about the underlying database implementation. Maybe you want to run things locally using a SQLite database, or you want to run it in "production" using a MySQL database or in a "serverless" environment using something like DynamoDB. The answer is: Yes. So long as your database of choice implements the different database (or table) interfaces then it is supported.

## Implementations

Like everything else in the `go-activitypub` package these databases implement a subset of all the functionality defined in the ActivityPub and ActivityStream standards and reflect the ongoing investigation in translating those standards in to working code.

As of this writing two "classes" of databases are supported: 

### database/sql

Anything that implements the built-in Go [database/sql](https://pkg.go.dev/database/sql) `DB` interface. As of this writing only SQLite databases have been tested using the [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) package. There are two things to note about the SQLite implementation:

* The use of the `mattn/go-sqlite3` package means you'll need to have a C complier to build the code.
* The SQLite databases themselves not infrequenely get themselves in to a state where they are locked preventing other operations from completing. This is a SQLite thing so you probably don't want to deploy it to "production" but it is generally good enough for testing things. It's also possible that I am simply misconfiguring SQLite and if I am I would appreciate any pointers on how to fix these mistakes.

To add a different database "driver", for example MySQL, you will need to clone the respective tools and add the relevant import statement. For example to update the [cmd/server](cmd/server/main.go) tool to use MySQL you would replace the `_ "github.com/mattn/go-sqlite3"` import statement with `_ "github.com/go-sql-driver/mysql"` like this:

```
package main

import (
	"context"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sfomuseum/go-activitypub/app/server"
	"github.com/sfomuseum/go-activitypub/slog"
)

func main() {
	ctx := context.Background()
	logger := slog.Default()
	server.Run(ctx, logger)
}
```

#### SQL schemas

* [SQLite](schema/sqlite)
* [MySQL](schema/mysql) – _Note: These have not been tested yet._

### gocloud.dev/docstore

Anything that implements the [gocloud.dev/docstore](https://pkg.go.dev/gocloud.dev/docstore) `Docstore` interface. As of this writing only DynamoDB document stores have been tested using the [awsdynamodb](https://gocloud.dev/howto/docstore/#dynamodb) and the [aaronland/gocloud-docstore](https://github.com/aaronland/gocloud-docstore/) packages. A few things to note:

* One side effect of using the `aaronland/gocloud-docstore` package is that the `gocloud.dev/docstore/awsdynamodb` "driver" is always imported and available regardless of whether you include in an `import` statement.
* The "global secondary indices" for the [schema/dynamodb](schema/dynamodb) definitions are inefficient and could stand to be optimized at some point in the future.

To add a different docstore "driver", for example MongoDB, you will need to clone the respective tools and add the relevant import statement. For example to update the [cmd/server](cmd/server/main.go) tool to use MongoDB you would replace the `_ "github.com/mattn/go-sqlite3"` import statement with `_ "gocloud.dev/docstore/mongodocstore"` like this:

```
package main

import (
	"context"
	"os"
	
	_ "gocloud.dev/docstore/mongodocstore"
	"github.com/sfomuseum/go-activitypub/app/server"
	"github.com/sfomuseum/go-activitypub/slog"
)

func main() {
	ctx := context.Background()
	logger := slog.Default()
	server.Run(ctx, logger)
}
```

#### Document Store table definitions

* [DynamoDB](schema/dynamodb)

## Interfaces

### AccountsDatabase

This is where accounts specific to an atomic instance of the `go-activitypub` package, for example `example1.social` versus `example2.social`, are stored.

### ActivitiesDatabase

This is where outbound (ActivityPub) actvities related to accounts are stored. This database stores both the raw (JSON-encoded) ActvityPub activity record as well as other properties (account id, activity type, etc.) specific to the `go-activitypub` package.

### AliasesDatabase

This is where aliases (alternate names) for accounts are stored.

### BlocksDatabase

This is where records describing external actors or hosts that are blocked by individual accounts are stored.

### DeliveriesDatabase

This is where a log of outbound deliveries of (ActivityPub) actvities for individual accounts are stored.

### FollowersDatabase

This is where records describing the external actors following (internal) accounts are stored.

### FollowingDatabase

This is where records describing the external actors that (internal) accounts are following are stored.

### LikesDatabase

This is where records describing boost activities by external actors (of activities by internal accounts) are stored.

_There are currently no database tables for storing boost events by internal accounts._

### MessagesDatabase

This is where the pointer to a "Note" activity (stored in an implementation of the `NotesDatabase`) delivered to a specific account is stored.

### NotesDatabase

This is where the details of a "Note" activity from an external actor is stored.

### PostTagsDatabase

This is where the details of the "Tags" associated with a post (or "Note") activitiy are stored.

_Currently, this interface is tightly-coupled with "Notes" (posts) which may be revisited in the future._

### PostsDatabase

This is where the details of a "Note" activity (or "post") by an account are stored. These details do not include ActivityStream "Tags" which are stored separately in an implementation of the `PostTagsDatabase`.

### PropertiesDatabase

This is where arbitrary key-value property records for individual accounts are stored.