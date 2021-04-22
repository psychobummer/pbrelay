## PBRelay

PsychoBummer's generic gRPC message relay for streaming ed25519-signed `[]byte` payloads from n-producers to n-subscribers.

## Installation 

`go get github.com/psychobummer/pbrelay`

## Usage

```
$ go build 
$ ./pbrelay server --help

Usage:
  pbrelay server [flags]

Flags:
  -h, --help            help for server
  -l, --listen string   ip:port to server should bind to (e.g: 0.0.0.0:9999) (required)
  -s, --validate        validate signature of producer messages (default: true) (default true)
```

## Overview (long version)

`pbrelay` consists of two gRPC services: a `SignedRelay` ([gRPC Service](https://github.com/psychobummer/pbrelay/blob/master/rpc/relay/relay.proto), [Implementation](https://github.com/psychobummer/pbrelay/blob/master/signed_relay/server.go)), and a `Keystore` ([gRPC Service](https://github.com/psychobummer/pbrelay/blob/master/rpc/keystore/keystore.proto), [Implementation](https://github.com/psychobummer/pbrelay/blob/master/keystore/server.go)).

To use the relay, a producer will first send an `ed25519` public key to the `Keystore`. The `Keystore` will generate a random ID for that session, store the producer's public key under that ID in a backing storage engine (e.g: [MemoryStorage](https://github.com/psychobummer/pbrelay/blob/master/keystore/memory_storage.go)) and return the ID to the producer; this ID acts as the "stream id," and can be shared with anyone out of band to start receiving data.

The producer can then start streaming `<id, data, signature(data)>` `relay.StreamMessage` protobuffs to the relay. The relay will validate the `signature` of the first `relay.StreamMessage` received against the key stored under `id`, and if the signature looks good, the relay will forward `data` to any subscribers for the stream `id`. A full producer client example can be found [here](https://github.com/psychobummer/pbrelay-producer/blob/master/midiproducer/midiproducer.go)

To use the relay as a `subscriber`, things are a bit easier. The subscriber just needs to call the `GetStream(context.Context, pb.GetStreamRequest{})` method, and if they're so inclined, (optionally) validate the signature of any data it receives.

The distribution of `relay.StreamMessage` protobuffs from producers to subscribers is mediated through [the broker](https://github.com/psychobummer/pbrelay/blob/master/signed_relay/broker.go)

A full subscriber client example can be found [here](https://github.com/psychobummer/pbrelay-subscriber/blob/master/cmd/midibtle.go)

## Overview (short version)

* pbrelay producer sends public a key to the keystore, gets back an `id`
* pbrelay producer streams `relay.StreamMessage` protobuffs to relay server
* pbrelay subscriber can subscribe to `id` and receive `relay.StreamMessage` protobuffs.
* [example producer](https://github.com/psychobummer/pbrelay-producer/blob/master/midiproducer/midiproducer.go)
* [example subscriber](https://github.com/psychobummer/pbrelay-subscriber/blob/master/cmd/midibtle.go)

## I also a need a "back of the napkin whatever it's good enough why isnt there a good 1:n grpc relay package already?" solution for my thing

There's nothing magic about this; you can just attach the exposed services to your own gRPC server, with the caveat that if you want the relay to valiate the signature of messages from a producer, you'll need to provide it a functioning keystore client. e.g:

```golang

import (
    ...
    ...
    pbk "github.com/psychobummer/pbrelay/rpc/keystore"
    pbr "github.com/psychobummer/pbrelay/rpc/relay"
)

...
...

keystoreServer := keystore.NewServer()
relayServer := NewServer(Config{
	KeystoreClient:            pbk.NewKeystoreServiceClient(conn),
	ValidateProducerSignature: true,
})

grpcServer := grpc.NewServer()
pbr.RegisterRelayServiceServer(grpcServer, relayServer)
pbk.RegisterKeystoreServiceServer(grpcServer, keystoreServer)
if err := grpcServer.Serve(listener); err != nil {
	panic(err)
}
```

## Caveats

* We don't make any assumptions about the form of `StreamMessage.data` payloads
* If you want to run multiple relay servers, you'll need a load balancing layer in front of them which hashes subscribers to the apporpriate relays. Our scale doesn't require this yet, so we don't support mesh'd `relay<->relay` topologies, where you can connect to one relay, and it'll go figure out which relay actually has that stream. Honestly, if you need that, you most likely want to use NATs, Kafka, zeroMQ, rabbit, redis or some other pubsub-capable system and simply build a facade atop that. We have limited need and didn't want additional infra/compute overhead.
