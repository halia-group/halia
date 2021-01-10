package main

import (
	"halia/channel"
)

type StringToByteEncoder struct{}

func (e *StringToByteEncoder) OnError(c channel.HandlerContext, err error) {
	c.FireOnError(err)
}

func (e *StringToByteEncoder) Write(c channel.HandlerContext, msg interface{}) error {
	if str, ok := msg.(string); ok { // string才转换
		return c.Write([]byte(str))
	}
	return c.Write(msg)
}

func (e *StringToByteEncoder) Flush(c channel.HandlerContext) error {
	return c.Flush()
}
