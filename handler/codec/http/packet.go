package http

import "fmt"

type Request struct {
	Version string
	Method  string
	Path    string
	Headers map[string][]string
	Body    []byte
}

func (r Request) String() string {
	return fmt.Sprintf("Request{Version=%s, Method=%s, Path=%s}", r.Version, r.Method, r.Path)
}

type Response struct {
	Version string
	Status  int
	Reason  string
	Headers map[string][]string
	Body    []byte
}
