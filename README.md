# PubSub Examples

This repositories illustrates a few things about PubSub systems, how
they work, how they differ, and how to implement them in Go.


## Prerequisites

Although also mentioned in the README files of the subfolders, you
should probably just start open at least 4 console windows and set
these environment variables for the rest of the examples.

### Google PubSub 
These make [Google PubSub emulator](https://cloud.google.com/pubsub/docs/emulator)
work locally.

```
$ export PUBSUB_EMULATOR_HOST=localhost:8086
$ export PUBSUB_PROJECT_ID=demo
```

### NATS PubSub 
This one is for [NATS Streaming Server](https://github.com/nats-io/nats-streaming-server).

```
$ export NATS_SERVER_URL=localhost:4222
```

Also, make sure you have a recent version of Docker and Docker Compose
installed. We start the infrastructure through Docker Compose, so you
don't have to install anything besides Go on your machine.



## Differences in PubSub systems

Some PubSub systems create topics on the fly, some systems require the caller
to create the topic before sending (or trying to receive) the first message.

On the subscriber side, some systems use a pull model (the consumer asks the
PubSub system) or a push model (the PubSub system forwards new messages), and
some offer both. On top of that, both models often have a synchronous or an
asynchronous implementation. Google PubSub, for example, has all of this: Push,
pull, sync/async.

What's also interesting is the guarantees that PubSub systems give you in terms
of ordering, delivery, and duplicates. In general, it is very hard to get
"exactly-once" delivery, and I think only Kafka can give this guarantee at the
time of writing this. So your code must be prepared for the case that you get
a message twice, or a message delivered later to come first to a subscriber.

Notice that some PubSub systems remove messages when they are consumed. Kafka,
on the other hand, doesn't do that. Instead, Kafka has a TTL for each message
and will remove them later. So setting the TTL to "never" will basically store
all messages forever, allowing you to iterate over historic messages years later.


## Examples

Now, this repository comes with 3 examples, each on in a subdirectory.

The [`gcp`](https://github.com/olivere/pubsub-example/tree/master/gcp)
and [`nats`](https://github.com/olivere/pubsub-example/tree/master/nats)
subdirectories illustrate how to implement a rather simply publisher and
subscriber for Google PubSub (via the Emulator) and NATS Streaming Server.

In [`gocloud`](https://github.com/olivere/pubsub-example/tree/master/gocloud)
you can see a single implementation that uses the wonderful
[`gocloud.dev` library by Google](https://github.com/google/go-cloud) to
do PubSub with a single implementation. You can use the code with any of the
supported PubSub systems: AWS SNS/SQS, Azure SB, Google PubSub, Kafka,
RabbitMQ, NATS Streaming Server, and an in-memory implementation which is
perfect for testing. You control the system to use by passing a URL.

Now, open 4 terminals and start the infrastructure in the 1st one:

### Start All Examples -- Run the GCP one locally
```
$ docker-compose up
```

Head to the 2nd terminal and run:

```
$ cd gcp
$ make
$ ./pub -h
Usage of ./pub:
  -p string
    	Project ID
  -t string
    	Topic name
$ ./pub -t messages
2019/07/29 17:51:12 pub.go:66: Send: 00000001 2019-07-29T17:51:12+02:00
2019/07/29 17:51:12 pub.go:66: Send: 00000002 2019-07-29T17:51:12+02:00
2019/07/29 17:51:13 pub.go:66: Send: 00000003 2019-07-29T17:51:13+02:00
2019/07/29 17:51:14 pub.go:66: Send: 00000004 2019-07-29T17:51:14+02:00
2019/07/29 17:51:15 pub.go:66: Send: 00000005 2019-07-29T17:51:15+02:00
2019/07/29 17:51:15 pub.go:66: Send: 00000006 2019-07-29T17:51:15+02:00
...
```

Go to the 3rd terminal and run:

```
$ cd gcp
$ ./sub -h
Usage of ./sub:
  -p string
    	Project ID
  -s string
    	Subscription name
  -t string
    	Topic name
$ ./sub -t messages -s subscriber-1
2019/07/29 17:51:12 sub.go:87: Recv: 00000001 2019-07-29T17:51:12+02:00
2019/07/29 17:51:12 sub.go:87: Recv: 00000002 2019-07-29T17:51:12+02:00
2019/07/29 17:51:13 sub.go:87: Recv: 00000003 2019-07-29T17:51:13+02:00
2019/07/29 17:51:14 sub.go:87: Recv: 00000004 2019-07-29T17:51:14+02:00
2019/07/29 17:51:15 sub.go:87: Recv: 00000005 2019-07-29T17:51:15+02:00
2019/07/29 17:51:16 sub.go:87: Recv: 00000006 2019-07-29T17:51:15+02:00
```

Up to the 4th terminal and run:

```
$ cd gcp
$ ./sub -t messages -s subscriber-2
2019/07/29 17:51:12 sub.go:87: Recv: 00000001 2019-07-29T17:51:12+02:00
2019/07/29 17:51:12 sub.go:87: Recv: 00000002 2019-07-29T17:51:12+02:00
2019/07/29 17:51:13 sub.go:87: Recv: 00000003 2019-07-29T17:51:13+02:00
2019/07/29 17:51:14 sub.go:87: Recv: 00000004 2019-07-29T17:51:14+02:00
2019/07/29 17:51:15 sub.go:87: Recv: 00000005 2019-07-29T17:51:15+02:00
2019/07/29 17:51:16 sub.go:87: Recv: 00000006 2019-07-29T17:51:15+02:00
...
```

Notice how both `subscriber-1` and `subscriber-2` get a copy of each message
being sent by the publisher on topic `messages`.

Now, stop the `subscriber-2` in terminal 4, and run it under the same
name as `subscriber-1` instead:

```
$ ./sub -t messages -s subscriber-1
2019/07/29 17:52:05 sub.go:87: Recv: 00000007 2019-07-29T17:52:05+02:00
2019/07/29 17:52:07 sub.go:87: Recv: 00000010 2019-07-29T17:52:07+02:00
...
```

Notice how messages will be split between terminal 3 and 4, load-balanced if you will.

# License

MIT.
