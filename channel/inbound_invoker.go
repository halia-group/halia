package channel

type InboundInvoker interface {
	FireChannelActive()
	FireChannelInActive()
	FireChannelRead(msg interface{})
}
