package channel

type OutboundInvoker interface {
	Write(msg interface{}) error
	Flush() error
	WriteAndFlush(msg interface{}) error
}

type InboundInvoker interface {
	FireChannelActive()
	FireChannelInActive()
	FireChannelRead(msg interface{})
	FireErrorCaught(err error)
}
