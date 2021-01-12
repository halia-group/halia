package codec

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/halia-group/halia/channel"
)

var (
	ErrIllegalLengthField = errors.New("illegal length field")
)

type LengthFieldBasedFrameDecoder struct {
	Decoder
	lengthFieldLength int              // 包长度字段长度
	lengthFieldOffset int              // 包长度偏移
	byteOrder         binary.ByteOrder // 字节序

}

func NewLengthFieldBasedFrameDecoder(lengthFieldLength, lengthFieldOffset int, byteOrder binary.ByteOrder) *LengthFieldBasedFrameDecoder {
	return &LengthFieldBasedFrameDecoder{lengthFieldLength: lengthFieldLength, lengthFieldOffset: lengthFieldOffset, byteOrder: byteOrder}
}

func (d LengthFieldBasedFrameDecoder) ChannelRead(c channel.HandlerContext, _ interface{}) {
	var (
		scanner = bufio.NewScanner(c.Channel())
	)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if !atEOF && len(data) >= (d.lengthFieldLength+d.lengthFieldOffset) {
			var (
				lengthFieldBuf    = make([]byte, d.lengthFieldLength)
				lengthFieldReader = bytes.NewReader(data[d.lengthFieldOffset : d.lengthFieldOffset+d.lengthFieldLength])
				bodyLength        = 0
			)

			if err = binary.Read(lengthFieldReader, d.byteOrder, &lengthFieldBuf); err != nil {
				return
			}
			if bodyLength = d.adjustFrameLength(lengthFieldBuf); bodyLength == -1 {
				err = ErrIllegalLengthField
				return
			}
			fullPacketLength := bodyLength + d.lengthFieldOffset + d.lengthFieldLength
			if len(data) >= fullPacketLength {
				return fullPacketLength, data[:fullPacketLength], nil
			}
		}
		return
	})
	for scanner.Scan() {
		c.FireChannelRead(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		c.FireOnError(err)
	}
}

func (d LengthFieldBasedFrameDecoder) adjustFrameLength(buf []byte) int {
	switch len(buf) {
	case 1:
		return int(buf[0])
	case 2:
		return int(d.byteOrder.Uint16(buf))
	case 4:
		return int(d.byteOrder.Uint32(buf))
	case 8:
		return int(d.byteOrder.Uint64(buf))
	default:
		return -1
	}
}
