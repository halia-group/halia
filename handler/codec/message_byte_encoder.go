package codec

import (
	"halia/channel"
	"io"
)

type MessageToByteEncoder struct{ channel.OutboundHandlerAdapter }

func (m MessageToByteEncoder) Write(c channel.HandlerContext, msg interface{}) error {
	return m.Encode(c, msg, c.Channel())
}

func (MessageToByteEncoder) Encode(c channel.HandlerContext, msg interface{}, w io.Writer) error {
	panic("implement me")
}
