package channel

type OutboundInvoker interface {
	Write(msg interface{}) error
	Flush() error
	WriteAndFlush(msg interface{}) error
	Close() error
}
