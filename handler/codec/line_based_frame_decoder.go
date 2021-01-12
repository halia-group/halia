package codec

import (
	"bufio"
	"github.com/halia-group/halia/channel"
	"io"
)

type LineBasedFrameDecoder struct {
	Decoder
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
