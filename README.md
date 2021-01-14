# Halia

Halia是一个基于组件化的网络应用框架，用于快速开发可维护的高性能协议服务器和客户端。

## Links
+ [文档地址](https://halia-group.github.io/halia/)

## Features

+ 组件化
+ 可扩展
+ 高性能

## 数据流

+ 入站数据流: HeadInboundHandler -> inboundHandler1 -> inboundHandlerN
+ 出站数据流: inboundHandlerN.Write -> HeadOutboundHandler -> outboundHandler1 -> outboundHandlerN -> TailOutboundHandler

## 内置解码器

+ DebugEncoder/DebugDecoder: 将流经数据打印在标准错误输出，方便调试
+ FixedLengthFrameDecoder: 定长报文解码器
+ LengthFieldBasedFrameDecoder: 基于长度字段的变长报文解码器
+ LineBasedFrameDecoder: 基于换行符的报文解码器

## 示例

+ [Echo](examples/echo)，回显服务器，基于LineBasedFrameDecoder实现
+ [Time](examples/time)，时间服务器，基于FixedLengthFrameDecoder实现，报文固定为8字节长度时间戳
+ [Chat](https://github.com/halia-group/halia-chat)，聊天服务器，基于LengthFieldBasedFrameDecoder实现，低耦合/高扩展性，扩展数据包只需要注册到`PacketFactory`和`ProcessorFactory`即可。
