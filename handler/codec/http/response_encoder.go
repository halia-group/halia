package http

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/halia-group/halia/channel"
	"strconv"
)

type ResponseEncoder struct{}

func (r ResponseEncoder) OnError(c channel.HandlerContext, err error) {
	c.FireOnError(err)
}

func (r ResponseEncoder) Write(c channel.HandlerContext, msg interface{}) error {
	resp, ok := msg.(*Response)
	if !ok {
		return errors.New("invalid msg type")
	}
	buf := bytes.Buffer{}
	// 编码响应行
	if err := r.writeResponseLine(&buf, resp); err != nil {
		return err
	}
	// 编码响应头
	if err := r.writeResponseHeaders(&buf, resp); err != nil {
		return err
	}
	// 编码响应体
	if err := r.writeResponseBody(&buf, resp); err != nil {
		return err
	}
	return c.Write(buf.Bytes())
}

func (r ResponseEncoder) Flush(c channel.HandlerContext) error {
	return c.Flush()
}

// 编码响应行
func (r ResponseEncoder) writeResponseLine(w *bytes.Buffer, resp *Response) error {
	if _, err := w.WriteString(fmt.Sprintf("%s %d %s\r\n", resp.Version, resp.Status, resp.Reason)); err != nil {
		return err
	}
	return nil
}

// 写入响应头
func (r ResponseEncoder) writeResponseHeaders(w *bytes.Buffer, resp *Response) error {
	if resp.Headers == nil {
		resp.Headers = make(map[string][]string)
	}
	key := "Content-Length"
	resp.Headers[key] = append(resp.Headers[key], strconv.Itoa(len(resp.Body)))
	for name, vals := range resp.Headers {
		for _, val := range vals {
			if _, err := w.WriteString(fmt.Sprintf("%s: %s\r\n", name, val)); err != nil {
				return err
			}
		}
	}
	if _, err := w.WriteString("\r\n"); err != nil {
		return err
	}
	return nil
}

// 编码响应体
func (r ResponseEncoder) writeResponseBody(w *bytes.Buffer, resp *Response) error {
	if _, err := w.Write(resp.Body); err != nil {
		return err
	}
	return nil
}
