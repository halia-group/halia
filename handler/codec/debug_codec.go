package codec

import (
	"encoding/hex"
	"fmt"
	"halia/channel"
	"os"
)

type DebugDecoder struct {
	Decoder
}

func NewDebugDecoder() *DebugDecoder {
	return &DebugDecoder{}
}

func (d DebugDecoder) ChannelRead(c channel.HandlerContext, msg interface{}) {
	buf := msg.([]byte)
	if _, err := hex.Dumper(os.Stderr).Write(buf); err != nil {
		d.OnError(c, err)
	}
	fmt.Fprintln(os.Stderr)
	c.FireChannelRead(msg)
}

type DebugEncoder struct{}

func NewDebugEncoder() *DebugEncoder {
	return &DebugEncoder{}
}

func (d DebugEncoder) OnError(c channel.HandlerContext, err error) {
	c.FireOnError(err)
}

func (d DebugEncoder) Write(c channel.HandlerContext, msg interface{}) error {
	buf := msg.([]byte)
	if _, err := hex.Dumper(os.Stderr).Write(buf); err != nil {
		d.OnError(c, err)
	}
	fmt.Fprintln(os.Stderr)
	return c.Write(msg)
}

func (d DebugEncoder) Flush(c channel.HandlerContext) error {
	return c.Flush()
}
