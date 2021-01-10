package main

import (
	"encoding/binary"
	"halia/channel"
	"halia/handler/codec"
	"time"
)

type timestampDecoder struct {
	codec.FrameDecoder
}

func (d timestampDecoder) ChannelRead(c channel.HandlerContext, msg interface{}) {
	data := msg.([]byte)
	timestamp := binary.BigEndian.Uint64(data)
	c.FireChannelRead(time.Unix(int64(timestamp), 0))
}

func newTimestampDecoder() *timestampDecoder {
	return &timestampDecoder{}
}
