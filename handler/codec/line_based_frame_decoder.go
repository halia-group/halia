package codec

import (
	"bufio"
	"halia/channel"
	"io"
)

type LineBasedFrameDecoder struct {
	FrameDecoder
}

func NewLineBasedFrameDecoder() *LineBasedFrameDecoder {
	return &LineBasedFrameDecoder{}
}

func (d *LineBasedFrameDecoder) ChannelRead(c channel.HandlerContext, _ interface{}) {
	var br = bufio.NewReader(c.Channel())
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			return
		}
		if err != nil {
			d.OnError(c, err)
			return
		}
		c.FireChannelRead(line)
	}
}
