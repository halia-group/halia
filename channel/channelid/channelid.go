/*
 *
 *  * MIT License
 *  *
 *  * Copyright (c) [2021] [xialeistudio]
 *  *
 *  * Permission is hereby granted, free of charge, to any person obtaining a copy
 *  * of this software and associated documentation files (the "Software"), to deal
 *  * in the Software without restriction, including without limitation the rights
 *  * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  * copies of the Software, and to permit persons to whom the Software is
 *  * furnished to do so, subject to the following conditions:
 *  *
 *  * The above copyright notice and this permission notice shall be included in all
 *  * copies or substantial portions of the Software.
 *  *
 *  * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  * SOFTWARE.
 *
 */

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
