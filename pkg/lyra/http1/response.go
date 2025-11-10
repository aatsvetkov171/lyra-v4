package http1

import (
	"fmt"
	"strconv"
	"strings"
)

var statusStrings = map[int]string{
	200: "OK",
	404: "Not Found",
	405: "Method Not Allowed",
}

type Response struct {
	statusCode int
	proto      string
	headers    map[string]string
	body       []byte
	filename   string
}

func NewResponse(statusCode int) *Response {
	newResponse := Response{
		statusCode: statusCode,
		proto:      "HTTP/1.1",
		headers:    make(map[string]string),
		filename:   "nofile",
	}
	newResponse.headers["Content-Type"] = "text/html; charset=UTF-8"
	newResponse.headers["Server"] = "Lyra/0.1"
	newResponse.headers["Connection"] = "keep-alive"
	return &newResponse
}

func (response *Response) AddHeader(key string, val string) {
	response.headers[key] = val
}

func (response *Response) GetHeaders() map[string]string {
	return response.headers
}

func (response *Response) GetHeadersBytes() []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s %d %s\r\n", response.proto, response.statusCode, statusStrings[response.statusCode])
	for k, v := range response.headers {
		fmt.Fprintf(&sb, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(&sb, "\r\n")
	result := []byte(sb.String())
	return result
}

func (response *Response) GetBody() []byte {
	return response.body
}

func (response *Response) GetFileName() string {
	return response.filename
}

func (response *Response) AddString(str string) {
	response.body = []byte(str)
	size := strconv.Itoa(len(response.body))
	response.AddHeader("Content-Length", size)
}

func (response *Response) AddFile(filename string) {
	response.filename = filename
}
