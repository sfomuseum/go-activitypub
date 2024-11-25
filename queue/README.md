# Queues

There two type of queues in the `go-activitypub` package. "Delivery" queues handle the details of delivery ActivityPub activities to one or more recipients (inboxes). "Processing" queues handle additional, or custom, processing of events related to ActivityPub messages received by the `server` application.

There are currently two "processing" queues:

* A message processing queue with processes a message (which resolves to an ActivityPub "note") after its been received and recorded an account's inbox.
* A follower processing queue with processes a follow event (as in a remote actor following an account) after its been received.

It is an open question whether or not to support multiple processing queues which bring with it the hassle and complexity of managing an equal number of queue endpoints. On the other hand it's too soon to know what sort of information would need to be passed, and how, to a single user-defined processing endpoint. For now, the decision is to be explicit and configure each processing queue with its own dispatcher and receiver.

## Delivery queues

Delivery queues implement the `DeliveryQueue` interface:

```
type DeliveryQueue interface {
	DeliverActivity(context.Context, *deliver.DeliverActivityOptions) error
	Close(context.Context) error
}
```

### Implementations

The following implementations of the `DeliveryQueue` interface are available by default:

#### null://

This implementation will receive an activity but not do anything with it. It is akin to writing data to `/dev/null`.

#### pubsub://

This implementation will dispatch the activity unique `ActivityId` property to an underlying implementation of the `sfomuseum/go-pubsub/publisher.Publisher` interface. That ID is expected to have been recorded in the `ActivitiesDatabase` table and that it can be retrieved by whatever code receives the message.

See also:

* https://github.com/sfomuseum/go-pubsub

#### slog://

The implementation will log the activity using the default `log/slog` logger.

#### synchronous://

## Message processing queues

Message processing queues implement the `ProcessMessageQueue` interface:

```
type ProcessMessageQueue interface {
	ProcessMessage(context.Context, int64) error
	Close(context.Context) error
}
```

Currently, "messages" are considered to be ActivityPub "Create" activities with type "Note". Remember a "message" in the `go-activitypub` is a pointer to a note associated with a specific account. Messages are dispatched to a `ProcessMessageQueue` as a final step in the [www.InboxPostHandler](../www/inbox_post.go) in the [server](../app/server) application.

There is no default endpoint, or code, for receiving or processing those messages after they have been dispatched. That is left up to individual users to implement, out of bounds, as their needs suit them. There is an [example application for processing messages](../app/message/process/example) that you can use as "starter code" which can run from the command line or as a Lambda function. It does nothing more than validate the message, recipient account and associated note and logging those details.

### Implementations

#### null://

This implementation will receive a message (ID) but not do anything with it. It is akin to writing data to `/dev/null`.

#### pubsub://

This implementation will dispatch a message ID to an underlying implementation of the `sfomuseum/go-pubsub/publisher.Publisher` interface. That ID is expected to have been recorded in the `MessagesDatabase` table and that it can be retrieved by whatever code receives the message.

See also:

* https://github.com/sfomuseum/go-pubsub

#### slog://

The implementation will log the activity using the default `log/slog` logger.

## Follower processing queues

Follower processing queues implement the `ProcessFollowerQueue` interface:

```
type ProcessFollowerQueue interface {
	ProcessFollower(context.Context, int64) error
	Close(context.Context) error
}
```

This queue is dispatched to with the unique 64-bit ID of the [Follower](../follower.go) record created in the [FollowersDatabase](../database/followers_database.go) when a remote actor follows an account hosted by the `server` application. Messages are dispatched to a `ProcessFollowerQueue` as a final step processing "Follow" events in the [www.InboxPostHandler](../www/inbox_post.go) in the [server](../app/server) application.

There is no default endpoint, or code, for receiving or processing those messages after they have been dispatched. That is left up to individual users to implement, out of bounds, as their needs suit them. There is an [example application for processing messages](../app/follower/process/example) that you can use as "starter code" which can run from the command line or as a Lambda function. It does nothing more than validate the message, recipient account and associated note and logging those details.

### Implementations

#### null://

This implementation will receive a follower (ID) but not do anything with it. It is akin to writing data to `/dev/null`.

#### pubsub://

This implementation will dispatch a follower ID to an underlying implementation of the `sfomuseum/go-pubsub/publisher.Publisher` interface. That ID is expected to have been recorded in the `FollowersDatabase` table and that it can be retrieved by whatever code receives the event.

See also:

* https://github.com/sfomuseum/go-pubsub

#### slog://

The implementation will log the follow(er) activity using the default `log/slog` logger.