# How does ActivityPub work?

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

## Accounts

Accounts are the local representation of an individual or what ActivityPub refers to as "actors". Accounts are distinguished from one another by the use of a unique name, for example `bob` or `alice.

Actors are distinguised from one another by the use of a unique "address" which consists of a name (`bob` or `alice`) and a hostname (`bob.com` or `alice.com`). For example `alice@alice.com` and `alice@bob.com` are two distinct "actors". In this example there are web services implementing the ActivityPub protocal available at both `bob.com` and `alice.com`.

Each actor (or account) has a pair of public-private encryption keys. As the name suggests the public key is available for anyone to view. Bob is authorized to see Alice's public key and vice versa. The private key however is considered sensitive and should only be visible to Alice or a service acting on Alice's behalf.

_The details of how any given private key is kept secure are not part of the ActivityPub specification and are left as implementation details to someone building a ActivityPub-based webs service._

## Exchanging messages

### Identifiers

_TBW_

### Signatures

_TBW_

### Call and response

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
 
## Posting messages (to followers)

This works (see the [#example](example section) below). I am still trying to work out the details.

## Endpoints

_To be written._

## Signing and verifying messages

_To be written. In the meantime consult [inbox.go](inbox.go), [actor.go](actor.go) and [www/inbox_post.go](www/inbox_post.go)._

## See also

* https://github.com/w3c/activitypub/blob/gh-pages/activitypub-tutorial.txt
* https://shkspr.mobi/blog/2024/02/activitypub-server-in-a-single-file/
* https://blog.joinmastodon.org/2018/07/how-to-make-friends-and-verify-requests/
* https://seb.jambor.dev/posts/understanding-activitypub/
* https://justingarrison.com/blog/2022-12-06-mastodon-files-instance/
* https://paul.kinlan.me/adding-activity-pub-to-your-static-site/