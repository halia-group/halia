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

package http

import (
	"bufio"
	"github.com/halia-group/halia/channel"
	"github.com/halia-group/halia/handler/codec"
	"io"
	"strconv"
	"strings"
)

type RequestDecoder struct{ codec.Decoder }

func (rd RequestDecoder) ChannelRead(c channel.HandlerContext, msg interface{}) {
	var (
		br = bufio.NewReader(c.Channel())
	)
	for {
		request := Request{
			Headers: make(map[string][]string),
			Body:    nil,
		}
		if err := rd.parseRequestLine(br, &request); err != nil {
			c.FireOnError(err)
			return
		}
		if err := rd.parseRequestHeaders(br, &request); err != nil {
			c.FireOnError(err)
			return
		}
		if err := rd.parseRequestBody(br, &request); err != nil {
			c.FireOnError(err)
			return
		}
		// 交付业务层
		c.FireChannelRead(&request)
	}
}

// 解析请求行
func (rd RequestDecoder) parseRequestLine(br *bufio.Reader, r *Request) error {
	line, _, err := br.ReadLine()
	if err != nil {
		return err
	}
	lines := strings.Split(string(line), " ")
	r.Method, r.Path, r.Version = lines[0], lines[1], lines[2]
	return nil
}

// 解析请求体，需要从header读取请求体长度
func (rd RequestDecoder) parseRequestBody(br *bufio.Reader, r *Request) error {
	contentLength, err := rd.getRequestContentLength(r)
	if err != nil {
		return err
	}
	if contentLength == 0 {
		return nil
	}
	r.Body = make([]byte, contentLength)
	if _, err := io.ReadFull(br, r.Body); err != nil {
		return err
	}
	return nil
}

// 获取请求体长度
func (rd RequestDecoder) getRequestContentLength(r *Request) (int, error) {
	// 没有请求体
	key := "content-length"
	if _, exists := r.Headers[key]; !exists {
		return 0, nil
	}
	return strconv.Atoi(r.Headers[key][0])
}

// 解析请求头
func (rd RequestDecoder) parseRequestHeaders(br *bufio.Reader, r *Request) error {
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			return err
		}
		// 请求头读取完毕
		if len(line) == 0 {
			return nil
		}
		pair := strings.Split(string(line), ":")
		name, value := strings.ToLower(strings.TrimSpace(pair[0])), strings.ToLower(strings.TrimSpace(pair[1]))
		if _, exists := r.Headers[name]; !exists {
			r.Headers[name] = make([]string, 0)
		}
		r.Headers[name] = append(r.Headers[name], value)
	}
}
