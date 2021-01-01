package codec

import (
	"bytes"
	"halia/channel"
	"io"
)

type ByteToMessageDecoder struct {
	channel.InboundHandlerAdapter
}

func (m *ByteToMessageDecoder) ChannelRead(c channel.HandlerContext, msg interface{}) (out interface{}, err error) {
	buf, ok := msg.([]byte)
	if !ok {
		return
	}
	out, err = m.Decode(c, bytes.NewReader(buf))
	return
}

func (ByteToMessageDecoder) Decode(c channel.HandlerContext, r io.Reader) (out interface{}, err error) {
	panic("implement me")
}
