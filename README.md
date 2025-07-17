# Halia

Halia is a component-based network application framework for rapid development of maintainable high-performance servers
and clients.

Reference netty to implementation. [Netty](https://netty.io/)

[中文文档](README-CN.md)

## Links

+ [Document](https://halia-group.github.io/halia/)

## Features

+ Component-Based
+ Extensible
+ High-Performance

## Data Stream

```
-------------------------------------------
        -> handler1 -> ... handlerN -> 
head                                    tail    
        <- handler1 <- ... handlerN <-       
--------------------------------------------
```

## Built-in Decoder

+ DebugEncoder/DebugDecoder: Print the data to stderr for easy debugging.
+ FixedLengthFrameDecoder
+ LengthFieldBasedFrameDecoder
+ LineBasedFrameDecoder

## Built-in Protocol Decoder

+ HTTP Protocol

## Examples

+ [Echo Application](examples/echo)，Implementation based on `LineBasedFrameDecoder`
+ [Time Application](examples/time)，Implementation based on `FixedLengthFrameDecoder`, The packet is fixed to an 8-byte
  length timestamp
+ [Chat Application](https://github.com/halia-group/halia-chat)
  Implementation based on `LengthFieldBasedFrameDecoder`, low coupling/high scalability, To add custom package 
  only needs to be registered in PacketFactory and ProcessorFactory.