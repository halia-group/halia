package channel

type HandlerContext interface {
	Channel() Channel
	Pipeline() Pipeline
	Write(msg interface{}) error
	WriteAndFlush(msg interface{}) error
	Flush() error
	Close() error
}
