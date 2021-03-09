# Halia

Halia是一个基于组件化的网络应用框架，用于快速开发可维护的高性能协议服务器和客户端。

## 基础概念

### Handler

数据包处理器，分为`InboundHandler`和`OutboundHandler`。`InboundHandler`负责处理入站数据，`OutboundHandler`负责处理出站数据。

### Decoder

解码器，`InboundHandler`的一个实现，负责从数据流中读取`[]byte`将其解析为特定格式的消息。

### Encoder

编码器，`OutboundHandler`的一个实现，负责将特定格式的消息编码为`[]byte`输出到下层协议。

### Channel

连接包装器，提供属性设置/获取方法。比如在用户登录成功后需要设置当前用户的连接为登录状态，则可以调用`Channel`的`SetAttribute`和`GetAttribute`方法。

### HandlerContext

Handler与Channel关联上下文对象，负责Handler的链式调用以及入站和出站方向转换。

### Pipeline

持有入站和出站Handler链，负责数据的链式调用。