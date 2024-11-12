# Example

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
		-dynamodb-client-uri 'awsdynamodb://?region=localhost&credentials=anon:&local=true'
```

Note that we passing a `TABLE_PREFIX` argument. This is to demonstrate how you can assign custom prefixes to the tables created in DynamoDB. You might want to do that because there are already one or more tables with the same names used by this package or because you want to run multiple, but distinct, ActivityPub services in the same DynamoDB environment.

Start the ActivityPub server:

```
$> make server TABLE_PREFIX=custom_
go run cmd/server/main.go \
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-followers-database-uri 'awsdynamodb://custom_followers?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-following-database-uri 'awsdynamodb://custom_following?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-notes-database-uri 'awsdynamodb://custom_notes?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-messages-database-uri 'awsdynamodb://custom_messages?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-blocks-database-uri 'awsdynamodb://custom_blocks?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
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
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-account-name bob \
		-account-icon-uri fixtures/icons/bob.jpg
go run cmd/add-account/main.go \
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-account-name alice \
		-allow-remote-icon-uri \
		-account-icon-uri https://static.sfomuseum.org/media/172/956/659/5/1729566595_kjcAQKRw176gxIieIWZySjhlNzgKNxoA_s.jpg
```

Next `bob` follows `alice`:

```
$> make follow TABLE_PREFIX=custom_
go run cmd/follow/main.go \
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-following-database-uri 'awsdynamodb://custom_following?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-messages-database-uri 'awsdynamodb://custom_messages?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
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
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-following-database-uri 'awsdynamodb://custom_following?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-messages-database-uri 'awsdynamodb://custom_messages?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
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
		-accounts-database-uri 'awsdynamodb://accounts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-followers-database-uri 'awsdynamodb://followers?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-posts-database-uri 'awsdynamodb://posts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-deliveries-database-uri 'awsdynamodb://deliveries?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
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
		-deliveries-database-uri 'awsdynamodb://deliveries?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-delivery-id 1762205712303263744 \
		-verbose
{"time":"2024-02-26T12:00:47.086377-08:00","level":"DEBUG","msg":"Verbose logging enabled"}
{"id":1762205712303263744,"activity_id":"http://localhost:8080/ap#as-f555a16d-9b09-452a-886f-0aae2cd52506","post_id":1762205712257126400,"account_id":1762201778486513664,"recipient":"bob@localhost:8080","inbox":"http://localhost:8080/ap/bob/inbox","created":1708977556,"completed":1708977556,"success":true}
```

Checking `bob`'s inbox we see the message from Alice:

```
$> make inbox TABLE_PREFIX=custom_ ACCOUNT=bob
go run cmd/inbox/main.go \
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-messages-database-uri 'awsdynamodb://custom_messages?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-notes-database-uri 'awsdynamodb://custom_notes?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-account-name bob
{"time":"2024-02-20T10:33:31.588716-08:00","level":"INFO","msg":"Get Note","message":1760009464636772352,"id":1760009464624189440}
{"time":"2024-02-20T10:33:31.592025-08:00","level":"INFO","msg":"NOTE","body":"{\"attributedTo\":\"fix me\",\"content\":\"This post left intentionally blank\",\"id\":\"be26c823-b75e-461c-adea-7260299a7434\",\"published\":\"2024-02-20T10:32:10-08:00\",\"to\":\"https://www.w3.org/ns/activitystreams#Public\",\"type\":\"Note\",\"url\":\"x-urn:fix-me#1760009464481583104\"}"}
```

`bob` unfollows `alice` and then removes all of Alice's posts from (Bob's) inbox:

```
$> make unfollow TABLE_PREFIX=custom_
go run cmd/follow/main.go \
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-following-database-uri 'awsdynamodb://custom_following?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-messages-database-uri 'awsdynamodb://custom_messages?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
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
		-accounts-database-uri 'awsdynamodb://custom_accounts?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-messages-database-uri 'awsdynamodb://custom_messages?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-notes-database-uri 'awsdynamodb://custom_notes?partition_key=Id&allow_scans=true&local=true&region=localhost&credentials=anon:' \
		-account-name bob
```
