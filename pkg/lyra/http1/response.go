package http1

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// версия http статус код статус текст \r\n
// headers \r\n
// \r\n
// body

var statusMap = map[string]string{
	"200": "OK",
	"404": "Not Found",
}

type Response struct {
	proto        string
	statusCode   string
	statusString string
	headers      map[string]string
	body         []byte
}

func NewResponse() *Response {
	response := Response{
		proto:        "HTTP/1.1",
		statusCode:   "200",
		statusString: statusMap["200"], // if ... ; ok ?
		headers: map[string]string{
			"Content-Type": "text/html; charset=utf-8",
			"Server":       "Lyra/1.0",
			"Connection":   "keep-alive",
		},
	}
	return &response
}

func (r *Response) AddBody(body []byte) {
	r.body = body
}

func (r *Response) AddHeader(key string, val string) {
	r.headers[key] = val
}

func (r *Response) Build() []byte {
	if _, ok := r.headers["Date"]; !ok {
		r.headers["Date"] = time.Now().Format(time.RFC1123)
	}
	if _, ok := r.headers["Content-Length"]; !ok {
		r.headers["Content-Length"] = strconv.Itoa(len(r.body))
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s %s %s\r\n", r.proto, r.statusCode, r.statusString))
	for k, v := range r.headers {
		b.WriteString(fmt.Sprintf("%s:%s\r\n", k, v))
	}
	b.WriteString("\r\n")
	b.Write(r.body)
	return []byte(b.String())
}
