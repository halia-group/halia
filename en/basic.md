# Halia

Halia is a component-based network application framework for rapid development of maintainable high-performance servers
and clients.

Reference netty to implementation. [Netty](https://netty.io/)

## Basic Concept

### Handler

Packet handler, There are two handlers: `InboundHandler` and `OutboundHandler`. `InboundHandler` handles inbound data
and `OutboundHandler` handles outbound data.

### Decoder

Packet Decoder，An implementation of `InboundHandler`, reading `[]byte` from connection and parsing it into a message of
specific format.

### Encoder

Packet Encoder，An implementation of `OutboundHandler`，encoding message in a specific format as `[]byte` then sending to
client

### Channel

A connection wrapper，Provides property setter/getter methods. For example, after the user logs in successfully, you need
to set the current user's connection to the `logged-in` state, you can call the `SetAttribute` and
`GetAttribute` methods of `Channel`.

### HandlerContext

The Handler is associated with the Channel context object. which is responsible for the chain invocation of the Handler
and the conversion of inbound and outbound directions.

### Pipeline

Holds inbound and outbound Handler chains, responsible for chained calls of data.