# CHANGELOG

## v0.2.0

1. 添加HTTP协议编解码器

## v0.2.1

1. `DefaultChannel`添加`buffer`处理，本次更新为兼容性更新

## v0.3.0

> 本次更新对业务层影响不大，注册handler时使用`AddLast`和`AddFirst`方法即可。

1. `Pipeline`: 
    1. `ChannelHandlerContext`使用双向链表代替两个单链表
    2. 使用`AddLast`和`AddFirst`添加处理器
2. `AttributeMap`: 添加`GetIntAttribute`方法