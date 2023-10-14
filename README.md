# udp_pubsub

A publish/subscribe message bus using UDP broadcast packets.

<img width="831" alt="image" src="https://github.com/nicwaller/udp_pubsub/assets/2850248/2ae818ae-cb68-40b1-b2ec-d2854d3e7bb2">

## Why

The [publish-subscribe channel](https://www.enterpriseintegrationpatterns.com/patterns/messaging/PublishSubscribeChannel.html) is a useful design pattern, and there are great tools like [Redis Pub/Sub](https://redis.io/docs/interact/pubsub/) to support it, but sometimes I just want something **very** lightweight.

The usual ways of designing inter-process communication (IPC) don't allow the fan-out afforded by publish-subscribe channels. Unix domain sockets and unix pipes are both designed to work with a single listener.

## How

UDP broadcast is actually perfect for this. The publisher sends broadcast packets on a known port, and any process on any host of the local network broadcast segment can become a subscriber just by listening for packets on that same port.

The thing that makes this work with multiple processes on the same host is use of [SO_REUSEPORT](https://lwn.net/Articles/542629/).

## Usage

```shell
(pushd publisher && go build && ./publisher) &
(pushd subscriber && go build && ./subscriber)
```

## Limitations

Broadcasting to `255.255.255.255` _really does_ send to every machine on the local broadcast segment of the network. This behaviour may be desirable (or not) depending on your needs and your network architecture.
