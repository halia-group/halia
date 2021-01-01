package channel

// 每一个handler都有一个context
type HandlerContext interface {
	OutboundInvoker

	Channel() Channel
	Name() string
	Handler() Handler
	Pipeline() Pipeline
	FireChannelRead(msg interface{})
	FireChannelActive()
	FireChannelInActive()
}
