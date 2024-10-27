# Databases

Databases in `go-activitypub` are really more akin to traditional RDBMS "tables" in that all of the database (or tables) listed below have one or more associations with each other. 

Each of these databases is defined as a Go-language interface with one or more default implementations provided by this package: A `database/sql` implementation that has been tested with SQLite databases; A `gocloud.dev/docstore` implementation that has been tested with DynamoDB; A "null" implementation which supports all the methods defined by its interface but which either returns no value or an error.

Like everything else in the `go-activitypub` package these databases implement a subset of all the functionality defined in the ActivityPub and ActivityStream standards and reflect the ongoing investigation in translating those standards in to working code.

## AccountsDatabase

This is where accounts specific to an atomic instance of the `go-activitypub` package, for example `example1.social` versus `example2.social`, are stored.

## ActivitiesDatabase

This is where outbound (ActivityPub) actvities related to accounts are stored. This database stores both the raw (JSON-encoded) ActvityPub activity record as well as other properties (account id, activity type, etc.) specific to the `go-activitypub` package.

## AliasesDatabase

This is where aliases (alternate names) for accounts are stored.

## BlocksDatabase

This is where records describing external actors or hosts that are blocked by individual accounts are stored.

## DeliveriesDatabase

This is where a log of outbound deliveries of (ActivityPub) actvities for individual accounts are stored.

## FollowersDatabase

This is where records describing the external actors following (internal) accounts are stored.

## FollowingDatabase

This is where records describing the external actors that (internal) accounts are following are stored.

## LikesDatabase

This is where records describing boost activities by external actors (of activities by internal accounts) are stored.

_There are currently no database tables for storing boost events by internal accounts._

## MessagesDatabase

This is where the pointer to a "Note" activity (stored in an implementation of the `NotesDatabase`) delivered to a specific account is stored.

## NotesDatabase

This is where the details of a "Note" activity from an external actor is stored.

## PostTagsDatabase

This is where the details of the "Tags" associated with a post (or "Note") activitiy are stored.

_Currently, this interface is tightly-coupled with "Notes" (posts) which may be revisited in the future._

## PostsDatabase

This is where the details of a "Note" activity (or "post") by an account are stored. These details do not include ActivityStream "Tags" which are stored separately in an implementation of the `PostTagsDatabase`.

## PropertiesDatabase

This is where arbitrary key-value property records for individual accounts are stored.