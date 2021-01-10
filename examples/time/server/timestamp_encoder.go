package main

import (
	"encoding/binary"
	"halia/channel"
	"time"
)

type timestampEncoder struct{}

func (e timestampEncoder) OnError(c channel.HandlerContext, err error) {
	c.FireOnError(err)
}

func (e timestampEncoder) Flush(c channel.HandlerContext) error {
	return c.Channel().Flush()
}

func (e timestampEncoder) Write(c channel.HandlerContext, msg interface{}) error {
	t := msg.(time.Time)
	value := t.Unix()
	return binary.Write(c.Channel(), binary.BigEndian, &value)
}
