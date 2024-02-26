# go-activitypub

An opionated (and incomplete) ActivityPub service implementation in Go.

## Motivation

I find the documentation for ActivityPub very confusing. I don't think I have any problem(s) with the underlying specification but I have not found any implementation guides that haven't left me feeling more confused than when I started. This includes the actual ActivityPub specifications published by the W3C which are no doubt thorough but, as someone with a finite of amount of competing time to devote to reading those specs, often feel counter-productive. Likewise, working implementations of the ActivityPub standards are often a confusing maze of abstractions that become necessary to do everything defined in the specs. There are some third-party guides, listed below, which are better than others but so far each one has felt incomplete in one way or another.

Importantly, the point is not that any of these things are "bad". They clearly aren't as evidenced by the many different working implementations of the ActivityPub standards in use today. The point is that the documentation, as it exists, hasn't been great for _me_. This repository is an attempt to understand all the moving pieces and their relationship to one another by working through the implementation of a simple ActivityPub service. It is incomplete by design and, if you are reading this, it's entirely possible that parts of it remain incorrect.

The goal is implement a basic web service and a set of command line tools which allow:

* Individual accounts to be created
* The ability for one account to follow, or unfollow, one another
* The ability for one account to block, or unblock, another account
* The ability for one account to post a message and to have that message relayed to one or more other accounts
* The ability for one account to see all the messages that have been delivered to them by other accounts

That's it, at least for now. It does have support for ActivityPub account migration, editing posts or notifications of changes to posts.

Importantly not all of those features have been implemented in both the web service and command line tools. This code is not something you can, or should, deploy as a hosted service for "strangers on the Internet". I have some fairly specific use-cases in mind for this code but the priority right now is just to understand the ActivityPub specification and the actual "brass tacks" of running a service that implements the specification.

The mechanics of the code are discussed later in this document.

## How does ActivityPub work?

Let's say there are two people, Bob and Alice, who want to exchange messages. A "message" might be text or images or video of some combination of all three. An "exchange" is the act of sending those messages from one person to another using an email-like addressing scheme but instead of using an email-specific protocol messages are sent over HTTP(S).

Both Bob and Alice have their own respective public-private key pairs. When Bob sends a message it is signed using Bob's _private key_. When Alice receives a message from Bob the authenticity of that message (trust that it was sent by Bob) is verified by Alice using Bob's _public_ key.

What needs to happen for this exchange of messages possible?

1. There needs to be one or more web servers (services) to broker the exchange of messages between Bob and Alice.
2. Those web services need to have the concept of "member" accounts, in this case Bob or Alice.
3. Each web service needs to implement an endpoint for looking up other ActivityPub-specific endpoints for each member account, namely there ActivityPub "inbox" and "outbox". The detail of the inbox and outbox are discussed below.
4. Some kind of persistent database for the web service to store information about member accounts, relationships between individual members and the people they want to send and receive messages from, the messages that have been sent and the messages that have been received.
5. Though not required an additional database to track accounts that an individual member does not want to interact with, referred to here as "blocking" is generally considered to be an unfortunate necessity.
6. A delivery mechanism to send messages published by Alice to all the people who have "followed" them (in this case Bob). The act of delivering a message consists of Alice sending that message to their "outbox" with a list of recipients. The "outbox" is resposible for coordinating the process of relaying that message to each recipient's ActivityPub (web service) "inbox".
7. In practice you will also need somewhere to store and serve account icon images from. This might be a filesystem, a remote hosting storage system (like AWS S3) or even by storing the images as base64-encoded blobs in one or your databases. The point is that there is a requirement for this whole other layer of generating, storing, tracking and serveing account icon images. _Note: The code included in this package has support for generating generic coloured-background-and-capital-letter style icons on demand but there are plenty of scenarios where those icons might be considered insufficient._

To recap, we've got:

1. A web server with a minimum of four endpoints: webfinger, actor, inbox and outbox
2. A database with the following tables: accounts, followers, following, posts, messages, blocks
3. Two member accounts: Bob and Alice
4. A delivery mechanism for sending messages; this might be an in-process loop or an asynchronous message queue but the point is that it is a sufficiently unique part of the process that it deserves to be thought of as distinct from the web server or the database.
5. A web server, or equivalent platform, for storing and serving account icon images.

For the purposes of these examples and for testing the assumption is that Bob and Alice have member accounts on the same server.

Importantly, please note that there is no mention of how Bob or Alice are authenticated or authorized on the web server itself. The public-private key pairs, mentioned above, that are assigned to each member are soley for the purposes of signing and verifiying messages send bewteen one or more ActivityPub endpoints.

_As a practical matter what that means is: For the purposes of running a web service that implements ActivityPub-based message exchange you will need to implement some sort of credentialing system to distinguish Bob from Alice and to prevent Alice from sending messages on Bob's behalf._

### Accounts

Accounts are the local representation of an individual or what ActivityPub refers to as "actors". Accounts are distinguished from one another by the use of a unique name, for example `bob` or `alice.

Actors are distinguised from one another by the use of a unique "address" which consists of a name (`bob` or `alice`) and a hostname (`bob.com` or `alice.com`). For example `alice@alice.com` and `alice@bob.com` are two distinct "actors". In this example there are web services implementing the ActivityPub protocal available at both `bob.com` and `alice.com`.

Each actor (or account) has a pair of public-private encryption keys. As the name suggests the public key is available for anyone to view. Bob is authorized to see Alice's public key and vice versa. The private key however is considered sensitive and should only be visible to Alice or a service acting on Alice's behalf.

_The details of how any given private key is kept secure are not part of the ActivityPub specification and are left as implementation details to someone building a ActivityPub-based webs service._

### Exchanging messages

#### Identifiers

_TBW_

#### Signatures

_TBW_

#### Call and response

_TBW_

### Looking up and following accounts

So let's say that Doug is on a Mastodon instance called `mastodon.server` and wants to follow `bob@bob.com`. To do this Doug would start by searching for the address `@bob@bob.com`.

_Note: I am just using `bob.com` and `mastodon.server` as examples. They are not an actual ActivityPub or Mastodon endpoints._

The code that runs Mastodon will then derive the hostname (`bob.com`) from the address and construct a URL in the form of:

```
https://bob.com/.well-known/webfinger?resource=acct:bob@bob.com
```

Making a `GET` request to that URL is expected to return a [Webfinger](#) document which will look like this:

```
$> curl -s 'https://bob.com/.well-known/webfinger?resource=acct:bob@bob.com' | jq
{
  "subject": "acct:bob@bob.com",
  "links": [
    {
      "href": "https://bob.com/ap/bob",
      "type": "text/html",
      "rel": "http://webfinger.net/rel/profile-page"
    },
    {
      "href": "https://bob.com/ap/bob",
      "type": "application/activity+json",
      "rel": "self"
    }
  ]
}
```

The code will then iterate through the `links` element of the response searching for `rel=self` and `type=application/activity+json`. It will take the value of the corresponding `href` attribute and issue a second `GET` request assigning the HTTP `Accept` header to be `application/ld+json; profile="https://www.w3.org/ns/activitystreams"`.

(There's a lot of "content negotiation" going on in ActivityPub and is often the source of confusion and mistakes.)

This `GET` request is expected to return a "person" or "actor" resource in the form of:

```
$> curl -s -H 'Accept: application/ld+json; profile="https://www.w3.org/ns/activitystreams"' https://bob.com/ap/bob | jq
{
  "@context": [
    "https://www.w3.org/ns/activitystreams",
    "https://w3id.org/security/v1"
  ],
  "id": "https://bob.com/ap/bob",
  "type": "Person",
  "preferredUsername": "bob",
  "inbox": "https://bob.com/ap/bob/inbox",
  "outbox": "https://bob.com/ap/bob/outbox",  
  "publicKey": {
    "id": "https://bob.com/ap/bob#main-key",
    "owner": "https://bob.com/ap/bob",
    "publicKeyPem": "-----BEGIN RSA PUBLIC KEY-----\nMIICCgKCAgEAvzo9pTyEGXl9jbJT6zv1p+cEfDP2vVN8bbgBYsltYw5A8LutZD7A\nspATOPJ3i9w43dZCORjmyuAX/0qyljbLfwzx1IEBmeg/3EAs0ON8A8tIbfcmI9JE\nn47UVR+Vn1h6o1dsRFx7X+fGefRIm005f7H/GLbJYTAvTgW3HJcakQI9rbFhaqnT\nmq6E+eEVhFqORVRrBjFMmAMNv6kJHSDtJie2YW76Nd9lqgR1FKV5B2M3a6gtIWv4\nNLOnwHxc266kqllmVUW79LB/2yI9KogMXjbp+MB7NhbtndJTpn1vAMYvUYSwxPhW\nJbWTqq7yhQi7zNaEDmzgOUhDiehHmm2XAqyIhlFEVvdKdOXUpJuIzEyHyxfCTA8Q\nNB9kncrS+L8TNDwdraNBQzgL68sKGp9eE3Rv/H4oNsqDD0/N8FyYwIOy+1BDGa9E\nPlsd/8vDi/3Mf3OBjfj64QwQj3V689jq2S+M1JCX/3EC77p2thT61GZUIFy/VfFZ\nuHUpiPvaxMo9KehsjCNTeRyGwRDBnLv/MWgRwFNGrT2w/m+cafiYoALOI4YB2RF0\ntWS8wK+559zfkV8T+UuQNzZbGAa0q+IpuBMlQhhfiwhEb3Olw7SvTXQUnwPBwmQb\nbbg3Lffg2N2Qz7QN9G99MjFDHIXXSyKyO+/kLsM28pLbitAHmP2KeuUCAwEAAQ==\n-----END RSA PUBLIC KEY-----\n"
  },
  "following": "https://bob.com/ap/bob/following",
  "followers": "https://bob.com/ap/bob/followers",
  "discoverable": true,
  "published": "2024-02-20T15:55:17-08:00",
  "icon": {
    "type": "Image",
    "mediaType": "image/png",
    "url": "https://bob.com/ap/bob/icon.png"
  }
}
```

At this point Doug's Mastodon server (`mastodon.server`) will issue a `POST` request to `https://bob.com/ap/bob/inbox` (or whatever the value is of the `inbox` property in the document that is returned). The body of that request will be a "Follow" sctivity that looks like this:

{
   "@context" : "https://www.w3.org/ns/activitystreams",
   "actor" : "https://mastodon.server/users/doug",
   "id" : "https://mastodon.server/52c7a999-a6bb-4ce5-82ca-5f21aec51811",
   "object" : "https://bob.bom/ap/bob",
   "type" : "Follow"
}


Bob's server `bob.com` will then verify the request from Doug to follow Bob is valid by... _TBW_.

Bob's server will then create a local entry indiciating that Doug is following Bob and then post (as in HTTP `POST` method) an "Accept" message to Doug's inbox:

```
POST /users/doug/inbox HTTP/1.1
Host: mastodon.server
Content-Type: application/ld+json; profile="https://www.w3.org/ns/activitystreams"
Date: 2024-02-24T02:28:21Z
Digest: SHA-256=DrqW7OcDFoVsm/1G9mRx5576MkWm5rK5BwI0NglugJo=
Signature: keyId="https://bob.com/ap/bob",algorithm="hs2019",headers="(request-target) host date",signature="..."

{
  "@context": "https://www.w3.org/ns/activitystreams",
  "id": "0b8f64a3-2ab1-46c8-9f2c-4230a9f62689",
  "type": "Accept",
  "actor": "https://bob.com/ap/bob",
  "object": {
    "id" : "https://mastodon.server/52c7a999-a6bb-4ce5-82ca-5f21aec51811",  
    "type": "Follow",
    "actor": "https://mastodon.server/users/doug",
    "object": "https://bob.com/ap/bob"
  }
}
```

There are a fews things to note:

1. It appears that ActivityPub services sending messages to an inbox don't care about, and don't evaluate, responses that those inboxes return. Basically inboxes return a 2XX HTTP status code if everything went okay and everyone waits for the next message to arrive in an inbox before deciding what to do next. I am unclear if this is really true or not.
2. There is no requirement to send the `POST` right away. In fact many services don't because they want to allow people to manually approve followers and so final "Accept" messages are often sent "out-of-band".

For the purposes of this example the code is sending the "Accept" message immediately after the `HTTP 202 Accepted` response is sent in a Go language deferred (`defer`) function. As mentioned, it is unclear whether it is really necessary to send the "Accept" message in a deferred function (or whether it can be sent inline before the HTTP 202 response is sent). On the other there are accept activities which are specifically meant to happen "out-of-band", like follower requests that are manually approved, so the easiest way to think about things is that they will (maybe?) get moved in to its own delivery queue (distinct from posts) to happen after the inbox handler has completed.

Basically: Treat every message sent to the ActivityPub inbox as an offline task. I am still trying to determine if that's an accurate assumption but what that suggests is, especially for languages that don't have deferred functions (for example PHP), the minimal viable ActivityPub service needs an additional database and delivery queue for these kinds of activities.
 
### Posting messages (to followers)

This works (see the [#example](example section) below). I am still trying to work out the details.

### Endpoints

_To be written._

### Signing and verifying messages

_To be written. In the meantime consult [inbox.go](inbox.go), [actor.go](actor.go) and [www/inbox_post.go](www/inbox_post.go)._

## The Code

### Databases

The package liberally mixes up the terms "database" and "table". Generally each aspect of the ActivityPub service has been separated in to distinct "tables" each with its own Go language interface. For example the interface for adding and removing followers looks like this:

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

As of this writing two "classes" of databases are supported: 

#### database/sql

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

##### SQL schemas

* [SQLite](schema/sqlite)
* [MySQL](schema/mysql) – _Note: These have not been tested yet._

#### gocloud.dev/docstore

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

##### Document Store table definitions

* [DynamoDB](schema/dynamodb)

### Delivery Queues

_TBW_

The default delivery queue is [SynchronousDeliveryQueue](delivery_queue_synchronous.go) which delivers posts to each follower in the order they are received.

### Example

What follows are the output of the different "targets" in the [Makefile](Makefile) that is included with this package. These targets are designed to make it easier to test common scenarios and to provide a reference of how things need to be configured.

If you want to follow along and run these examples your self you will need the following:

* The [Go programming language](https://go.dev/dl)
* The [Docker Desktop](https://www.docker.com/products/docker-desktop/) runtime environment to launch a local instance of DynamoDB.
* (3) different terminal windows (or "consoles")

In console (1) start a local instance of DynamoDB. The easiest way to do this is using the Dockerfile that AWS provides:

```
$> docker run --rm -it -p 8000:8000 amazon/dynamodb-local
Initializing DynamoDB Local with the following configuration:
Port:	8000
InMemory:	true
Version:	2.2.1
DbPath:	null
SharedDb:	false
shouldDelayTransientStatuses:	false
CorsParams:	null
```

_Note: This is an ephemeral Docker container so when you shut it down all the data that has been saved (like accounts and ActivityPub activity below) will be deleted._

In console (2) create the necessary tables for the ActivityPub service in DynamoDB.

```
$> make dynamo-tables-local TABLE_PREFIX=custom_
go run -mod vendor cmd/create-dynamodb-tables/main.go \
		-refresh \
		-table-prefix custom_ \
		-dynamodb-client-uri 'awsdynamodb://?local=true'
```

Note that we passing a `TABLE_PREFIX` argument. This is to demonstrate how you can assign custom prefixes to the tables created in DynamoDB. You might want to do that because there are already one or more tables with the same names used by this package or because you want to run multiple, but distinct, ActivityPub services in the same DynamoDB environment.

Start the ActivityPub server:

```
$> make server TABLE_PREFIX=custom_
go run cmd/server/main.go \
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true' \
		-followers-database-uri 'awsdynamodb://custom_followers?partition_key=Id&allow_scans=true&local=true' \
		-following-database-uri 'awsdynamodb://custom_following?partition_key=Id&allow_scans=true&local=true' \
		-notes-database-uri 'awsdynamodb://custom_notes?partition_key=Id&allow_scans=true&local=true' \
		-messages-database-uri 'awsdynamodb://custom_messages?partition_key=Id&allow_scans=true&local=true' \
		-blocks-database-uri 'awsdynamodb://custom_blocks?partition_key=Id&allow_scans=true&local=true' \
		-allow-create \
		-verbose \
		-allow-remote-icon-uri \
		-hostname localhost:8080 \
		-insecure
{"time":"2024-02-20T10:29:49.505754-08:00","level":"DEBUG","msg":"Verbose logging enabled"}
{"time":"2024-02-20T10:29:49.506312-08:00","level":"INFO","msg":"Listening for requests","address":"http://localhost:8080"}
```

Note the `-insecure` flag. Normally it is expected that all ActivityPub communications will happen over an encrypted (HTTPS) connection but since we are testing things locally and all of our accounts (Bob and Alice) will reside on the same server, and because setting up self-signed TLS certificates locally is a chore, we're going to exchange messages over an insecure connection.

Also note the `-hostname` flag. This is when you are running the `server` tool in a "serverless" environment which likely has a different domain name than the one associated with public-facing ActivityPub server. If the `-hostname` flag is left empty then its value is derived from the `-server-uri` flag which defaults to "http://localhost:8080".

Switch to console (3) and create account records for `bob` and `alice`:

```
$> make accounts TABLE_PREFIX=custom_

> make accounts
go run cmd/add-account/main.go \
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true' \
		-account-name bob \
		-account-icon-uri fixtures/icons/bob.jpg
go run cmd/add-account/main.go \
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true' \
		-account-name alice \
		-allow-remote-icon-uri \
		-account-icon-uri https://static.sfomuseum.org/media/172/956/659/5/1729566595_kjcAQKRw176gxIieIWZySjhlNzgKNxoA_s.jpg
```

Next `bob` follows `alice`:

```
$> make follow TABLE_PREFIX=custom_
go run cmd/follow/main.go \
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true' \
		-following-database-uri 'awsdynamodb://custom_following?partition_key=Id&allow_scans=true&local=true' \
		-messages-database-uri 'awsdynamodb://custom_messages?partition_key=Id&allow_scans=true&local=true' \
		-account-name bob \
		-follow alice@localhost:8080 \
		-hostname localhost:8080 \
		-verbose \
		-insecure
{"time":"2024-02-20T10:31:12.002118-08:00","level":"DEBUG","msg":"Verbose logging enabled"}
{"time":"2024-02-20T10:31:12.034109-08:00","level":"DEBUG","msg":"Webfinger URL for resource","resource":"alice","url":"http://localhost:8080/well-known/.webfinger?resource=alice"}
{"time":"2024-02-20T10:31:12.044169-08:00","level":"DEBUG","msg":"Profile page for actor","actor":"alice","url":"http://localhost:8080/ap/alice"}
{"time":"2024-02-20T10:31:12.04633-08:00","level":"DEBUG","msg":"Post to inbox","inbox":"http://localhost:8080/ap/alice/inbox","key_id":"http://localhost:8080/ap/bob"}
{"time":"2024-02-20T10:31:12.080122-08:00","level":"INFO","msg":"Following successful"}
```

Then `bob` unfollows `alice`:

```
$> make unfollow TABLE_PREFIX=custom_
go run cmd/follow/main.go \
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true' \
		-following-database-uri 'awsdynamodb://custom_following?partition_key=Id&allow_scans=true&local=true' \
		-messages-database-uri 'awsdynamodb://custom_messages?partition_key=Id&allow_scans=true&local=true' \
		-account-name bob \
		-follow alice@localhost:8080 \
		-hostname localhost:8080 \
		-insecure \
		-verbose \
		-undo
{"time":"2024-02-20T10:31:26.454195-08:00","level":"DEBUG","msg":"Verbose logging enabled"}
{"time":"2024-02-20T10:31:26.474536-08:00","level":"DEBUG","msg":"Webfinger URL for resource","resource":"alice","url":"http://localhost:8080/well-known/.webfinger?resource=alice"}
{"time":"2024-02-20T10:31:26.479316-08:00","level":"DEBUG","msg":"Profile page for actor","actor":"alice","url":"http://localhost:8080/ap/alice"}
{"time":"2024-02-20T10:31:26.482626-08:00","level":"DEBUG","msg":"Post to inbox","inbox":"http://localhost:8080/ap/alice/inbox","key_id":"http://localhost:8080/ap/bob"}
{"time":"2024-02-20T10:31:26.521846-08:00","level":"INFO","msg":"Unfollowing successful"}
```

At some point `alice` posts a message and then delivers it to all of their followers (including `bob` who has followed `alice` again):

```
$> make post MESSAGE='This message left intentionally blank'
go run cmd/post/main.go \
		-accounts-database-uri 'awsdynamodb://accounts?partition_key=Id&allow_scans=true&local=true' \
		-followers-database-uri 'awsdynamodb://followers?partition_key=Id&allow_scans=true&local=true' \
		-posts-database-uri 'awsdynamodb://posts?partition_key=Id&allow_scans=true&local=true' \
		-deliveries-database-uri 'awsdynamodb://deliveries?partition_key=Id&allow_scans=true&local=true' \
		-account-name alice \
		-message "This message left intentionally blank" \
		-hostname localhost:8080 \
		-insecure \
		-verbose
{"time":"2024-02-26T11:59:16.636181-08:00","level":"DEBUG","msg":"Verbose logging enabled"}
{"time":"2024-02-26T11:59:16.66917-08:00","level":"DEBUG","msg":"Deliver post","post":1762205712257126400,"from":1762201778486513664,"to":"bob@localhost:8080"}
{"time":"2024-02-26T11:59:16.669221-08:00","level":"DEBUG","msg":"Webfinger URL for resource","resource":"bob","url":"http://localhost:8080/.well-known/webfinger?resource=acct%3Abob%40localhost%3A8080"}
{"time":"2024-02-26T11:59:16.673629-08:00","level":"DEBUG","msg":"Profile page for actor","actor":"bob","url":"http://localhost:8080/ap/bob"}
{"time":"2024-02-26T11:59:16.676805-08:00","level":"DEBUG","msg":"Post to inbox","inbox":"http://localhost:8080/ap/bob/inbox"}
{"time":"2024-02-26T11:59:16.676888-08:00","level":"DEBUG","msg":"Post to inbox","inbox":"http://localhost:8080/ap/bob/inbox","key_id":"http://localhost:8080/ap/alice"}
{"time":"2024-02-26T11:59:16.706987-08:00","level":"DEBUG","msg":"Response","inbox":"http://localhost:8080/ap/bob/inbox","code":202,"content-type":""}
{"time":"2024-02-26T11:59:16.707027-08:00","level":"DEBUG","msg":"Add delivery for post","delivery id":1762205712303263744,"post id":1762205712257126400,"recipient":"bob@localhost:8080","success":true}
```

Switching back to the console running the `server` tool you should see something like this:

```
{"time":"2024-02-20T10:32:10.414033-08:00","level":"INFO","msg":"Fetch key for sender","method":"POST","accept":"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"","path":"/ap/bob/inbox","remote_addr":"127.0.0.1:55271","account":"bob","account id":1760009124885565440,"sender_address":"alice@localhost:8080","activity-type":"Create","key id":"http://localhost:8080/ap/alice","key_id":"http://localhost:8080/ap/alice"}
{"time":"2024-02-20T10:32:10.416587-08:00","level":"DEBUG","msg":"Get following","account":1760009124885565440,"following":"alice@localhost:8080"}
{"time":"2024-02-20T10:32:10.419441-08:00","level":"DEBUG","msg":"Add note","uuid":"be26c823-b75e-461c-adea-7260299a7434","author":"alice@localhost:8080"}
{"time":"2024-02-20T10:32:10.419451-08:00","level":"DEBUG","msg":"Create new note","uuid":"be26c823-b75e-461c-adea-7260299a7434","author":"alice@localhost:8080"}
{"time":"2024-02-20T10:32:10.421078-08:00","level":"DEBUG","msg":"Return new note","id":1760009464624189440}
{"time":"2024-02-20T10:32:10.42109-08:00","level":"DEBUG","msg":"Get message","account":1760009124885565440,"note":1760009464624189440}
{"time":"2024-02-20T10:32:10.422455-08:00","level":"DEBUG","msg":"Add message","account":1760009124885565440,"note":1760009464624189440,"author":"alice@localhost:8080"}
{"time":"2024-02-20T10:32:10.422462-08:00","level":"DEBUG","msg":"Create new message","account":1760009124885565440,"note":1760009464624189440,"author":"alice@localhost:8080"}
{"time":"2024-02-20T10:32:10.424299-08:00","level":"INFO","msg":"Note has been added to messages","method":"POST","accept":"application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"","path":"/ap/bob/inbox","remote_addr":"127.0.0.1:55271","account":"bob","account id":1760009124885565440,"sender_address":"alice@localhost:8080","activity-type":"Create","key id":"http://localhost:8080/ap/alice","note uuid":"be26c823-b75e-461c-adea-7260299a7434","note id":1760009464624189440,"message id":1760009464636772352}
```

Did you notice the "Add delivery for post" debug message when posting the message? Deliveries for post messages are logged in a "deliveries database". These logs record where and when posts were sent and whether the delivery was successful. For example:

```
$> make delivery ID=1762205712303263744 
go run cmd/retrieve-delivery/main.go \
		-deliveries-database-uri 'awsdynamodb://deliveries?partition_key=Id&allow_scans=true&local=true' \
		-delivery-id 1762205712303263744 \
		-verbose
{"time":"2024-02-26T12:00:47.086377-08:00","level":"DEBUG","msg":"Verbose logging enabled"}
{"id":1762205712303263744,"activity_id":"http://localhost:8080/ap#as-f555a16d-9b09-452a-886f-0aae2cd52506","post_id":1762205712257126400,"account_id":1762201778486513664,"recipient":"bob@localhost:8080","inbox":"http://localhost:8080/ap/bob/inbox","created":1708977556,"completed":1708977556,"success":true}
```

Checking `bob`'s inbox we see the message from Alice:

```
$> make inbox TABLE_PREFIX=custom_ ACCOUNT=bob
go run cmd/inbox/main.go \
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true' \
		-messages-database-uri 'awsdynamodb://custom_messages?partition_key=Id&allow_scans=true&local=true' \
		-notes-database-uri 'awsdynamodb://custom_notes?partition_key=Id&allow_scans=true&local=true' \
		-account-name bob
{"time":"2024-02-20T10:33:31.588716-08:00","level":"INFO","msg":"Get Note","message":1760009464636772352,"id":1760009464624189440}
{"time":"2024-02-20T10:33:31.592025-08:00","level":"INFO","msg":"NOTE","body":"{\"attributedTo\":\"fix me\",\"content\":\"This post left intentionally blank\",\"id\":\"be26c823-b75e-461c-adea-7260299a7434\",\"published\":\"2024-02-20T10:32:10-08:00\",\"to\":\"https://www.w3.org/ns/activitystreams#Public\",\"type\":\"Note\",\"url\":\"x-urn:fix-me#1760009464481583104\"}"}
```

`bob` unfollows `alice` and then removes all of Alice's posts from (Bob's) inbox:

```
$> make unfollow TABLE_PREFIX=custom_
go run cmd/follow/main.go \
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true' \
		-following-database-uri 'awsdynamodb://custom_following?partition_key=Id&allow_scans=true&local=true' \
		-messages-database-uri 'awsdynamodb://custom_messages?partition_key=Id&allow_scans=true&local=true' \
		-account-name bob \
		-follow alice@localhost:8080 \
		-hostname localhost:8080 \
		-insecure \
		-verbose \
		-undo
{"time":"2024-02-20T10:33:57.141869-08:00","level":"DEBUG","msg":"Verbose logging enabled"}
{"time":"2024-02-20T10:33:57.161661-08:00","level":"DEBUG","msg":"Webfinger URL for resource","resource":"alice","url":"http://localhost:8080/well-known/.webfinger?resource=alice"}
{"time":"2024-02-20T10:33:57.168161-08:00","level":"DEBUG","msg":"Profile page for actor","actor":"alice","url":"http://localhost:8080/ap/alice"}
{"time":"2024-02-20T10:33:57.171233-08:00","level":"DEBUG","msg":"Post to inbox","inbox":"http://localhost:8080/ap/alice/inbox","key_id":"http://localhost:8080/ap/bob"}
{"time":"2024-02-20T10:33:57.197086-08:00","level":"INFO","msg":"Remove message","id":1760009464636772352}
{"time":"2024-02-20T10:33:57.19868-08:00","level":"INFO","msg":"Unfollowing successful"}
```

Checking `bob`'s inbox again yields no posts:

```
$> make inbox TABLE_PREFIX=custom_ ACCOUNT=bob
go run cmd/inbox/main.go \
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true' \
		-messages-database-uri 'awsdynamodb://custom_messages?partition_key=Id&allow_scans=true&local=true' \
		-notes-database-uri 'awsdynamodb://custom_notes?partition_key=Id&allow_scans=true&local=true' \
		-account-name bob
```

## See also

* https://github.com/w3c/activitypub/blob/gh-pages/activitypub-tutorial.txt
* https://shkspr.mobi/blog/2024/02/activitypub-server-in-a-single-file/
* https://blog.joinmastodon.org/2018/07/how-to-make-friends-and-verify-requests/
* https://seb.jambor.dev/posts/understanding-activitypub/
* https://justingarrison.com/blog/2022-12-06-mastodon-files-instance/
