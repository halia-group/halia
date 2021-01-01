package channelid

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"sync/atomic"
	"time"
)

var (
	machineId []byte
	processId = os.Getpid()
	sequence  uint32
	logger    = log.WithField("component", "channelId")
)

type ChannelId interface {
	AsShortText() string
	AsLongText() string
}

type DefaultId struct {
	data       []byte
	shortValue string
	longValue  string
}

func (p *DefaultId) AsShortText() string {
	if p.shortValue == "" {
		p.shortValue = hex.EncodeToString(p.data[len(p.data)-4:])
	}
	return p.shortValue
}

func (p *DefaultId) AsLongText() string {
	if p.longValue == "" {
		p.longValue = hex.EncodeToString(p.data)
	}
	return p.longValue
}

// 生成新的channelId
func New() ChannelId {
	var (
		buf          = bytes.Buffer{}
		nextSequence = atomic.AddUint32(&sequence, 1)
		timestamp    = time.Now().Unix()
		random       = rand.Int31()
	)

	binary.Write(&buf, binary.BigEndian, &machineId)
	binary.Write(&buf, binary.BigEndian, &processId)
	binary.Write(&buf, binary.BigEndian, &nextSequence)
	binary.Write(&buf, binary.BigEndian, &timestamp)
	binary.Write(&buf, binary.BigEndian, &random)

	return &DefaultId{data: buf.Bytes()}
}
