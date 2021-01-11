package codec

import (
	"bufio"
	"halia/channel"
	"io"
)

type FixedLengthFrameDecoder struct {
	Decoder
	length int
}

func (d FixedLengthFrameDecoder) ChannelRead(c channel.HandlerContext, _ interface{}) {
	var (
		br  = bufio.NewReader(c.Channel())
		buf = make([]byte, d.length)
	)
	for {
		_, err := io.ReadFull(br, buf)
		if err == io.EOF {
			return
		}
		if err != nil {
			d.OnError(c, err)
			return
		}
		c.FireChannelRead(buf)
	}
}

func NewFixedLengthFrameDecoder(length int) *FixedLengthFrameDecoder {
	return &FixedLengthFrameDecoder{length: length}
}
