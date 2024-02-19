# go-activitypub

An opionated (and incomplete) ActivityPub service implementation in Go.

## Motivation

I find the documentation for ActivityPub very confusing. I don't think I have any problem(s) with the underlying specification but I have not found any implementation guides that haven't left me feeling more confused that when I started. This includes the actual ActivityPub specifications published by the W3C which are no doubt thorough but, as someone with a finite of amount of competing time to devote to reading those specs, often feel counter-productive. There are some third-party guides, listed below, which are better than others but so far each one has felt incomplete in one way or another.

This repository is an attempt to working through the implementation of a simple ActivityPub service. It is incomplete by design and, if you are reading this, it's entirely possible that parts of it remain incorrect. The goal is implement a basic web service and a set of command line tools which allow:

* Individual accounts to be created
* The ability for one account to follow, or unfollow, one another
* The ability for one account to block, or unblock, another account
* The ability for one account to post a message and to have that message relayed to one or more other accounts
* The ability for one account to see all the messages that have been delivered to them by other accounts

That's it, at least for now. Importantly not all of those features have been implemented in both the web service and command line tools. This code is not something you can, or should, deploy as a hosted service for "strangers on the Internet". I have some fairly specific use-cases in mind for this code but the priority right now is just to understand the ActivityPub specification and the actual "brass tacks" of running a service that implements the specification.

The mechanics of the code are discussed later in this document.

## How does ActivityPub work?

Let's say there are two people, Bob and Alice, who want to exchange messages. A "message" might be text or images or video of some combination of all three. An "exchange" is the act of sending those messages from one person to another an email-like addressing scheme but instead of using an email-specific protocol messages are sent over HTTP(S).

Both Bob and Alice have their own respective public-private key pairs. When Bob sends a message it is signed using Bob's _private key_. When Alice receives a message from Bob the authenticity of that message (trust that it was sent by Bob) is verified by Alice using Bob's _public_ key.

What needs to happen for this exchange of messages possible?

1. There needs to be one or more web servers (services) to broker the exchange of messages between Bob and Alice.
2. Those web services need to have the concept of "member" accounts, in this case Bob or Alice.
3. Each web service needs to implement an endpoint for looking up other ActivityPub-specific endpoints for each member account, namely there ActivityPub "inbox" and "outbox". The detail of the inbox and outbox are discussed below.
4. Some kind of persistent database for the web service to store information about member accounts, relationships between individual members and the people they want to send and receive messages from, the messages that have been sent and the messages that have been received.
4a. Though not required an additional database to track accounts that an individual member does not want to interact with, referred to here as "blocking" is generally considered to be an unfortunate necessity.
5. A delivery mechanism to send messages published by Alice to all the people who have "followed" them (in this case Bob). The act of delivering a message consists of Alice sending that message to their "outbox" with a list of recipients. The "outbox" is resposible for coordinating the process of relaying that message to each recipient's ActivityPub (web service) "inbox".

To recap, we've got:

1. A web server with a minimum of three endpoints: webfinger, inbox, outbox
2. A database with the following tables: accounts, followers, following, posts, messages, blocks
3. Two member accounts: Bob and Alice
4. A delivery mechanism for sending messages; this might be an in-process loop or an asynchronous message queue but the point is that it is a sufficiently unique part of the process that it deserves to be thought of as distinct from the web server or the database.

For the purposes of these examples and for testing the assumption is that Bob and Alice have member accounts on the same server.

Importantly, please note that there is no mention of how Bob or Alice are authenticated or authorized on the web server itself. The public-private key pairs, mentioned above, that are assigned to each member are soley for the purposes of signing and verifiying messages send bewteen one or more ActivityPub endpoints.

_As a practical matter what that means is: For the purposes of running a web service that implements ActivityPub-based message exchange you will need to implement some sort of credentialing system to distinguish Bob from Alice and to prevent Alice from sending messages on Bob's behalf._

### Accounts

Accounts are the local representation of an individual or what ActivityPub refers to as "actors". Accounts are distinguished from one another by the use of a unique name, for example `bob` or `alice.

Actors are distinguised from one another by the use of a unique "address" which consists of a name (`bob` or `alice`) and a hostname (`bob.com` or `alice.com`). For example `alice@alice.com` and `alice@bob.com` are two distinct "actors". In this example there are web services implementing the ActivityPub protocal available at both `bob.com` and `alice.com`.

Each actor (or account) has a pair of public-private encryption keys. As the name suggests the public key is available for anyone to view. Bob is authorized to see Alice's public key and vice versa. The private key however is considered sensitive and should only be visible to Alice or a service acting on Alice's behalf.

_The details of how any given private key is kept secure are not part of the ActivityPub specification and are left as implementation details to someone building a ActivityPub-based webs service._

### Endpoints

#### Webfinger

#### Inbox

#### Outbox

## The Code

### Databases

### Accounts

### Follow

### Unfollow

### Block

### Unblock

### Posting (and delivering) a message

### Retrieving messages (received)

## See also

* https://github.com/w3c/activitypub/blob/gh-pages/activitypub-tutorial.txt
* https://shkspr.mobi/blog/2024/02/activitypub-server-in-a-single-file/
* https://blog.joinmastodon.org/2018/07/how-to-make-friends-and-verify-requests/
* https://seb.jambor.dev/posts/understanding-activitypub/