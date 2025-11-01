package http1

import (
	"strings"
)

type Request struct {
	method  string
	path    string
	proto   string
	headers map[string]string
	body    []byte
}

func (r *Request) parseHeaders(headers []byte) {
	headersStrings := strings.Split(string(headers), "\r\n")
	for i := 0; i < len(headersStrings); i++ {
		header := strings.SplitN(strings.TrimSpace(headersStrings[i]), ":", 2)
		if len(header) < 2 {
			continue
		}
		r.headers[strings.ToLower(header[0])] = header[1]
	}
}

func NewRequest(Fline []byte, headers []byte, body []byte) *Request {
	req := Request{
		headers: make(map[string]string),
		body:    body,
	}
	firstline := strings.Fields(strings.TrimSpace(string(Fline)))
	if len(firstline) >= 3 {
		req.method = strings.ToUpper(firstline[0])
		req.path = firstline[1]
		req.proto = strings.ToUpper(firstline[2])
	} else {
		req.method = "UNKHOW"
		req.path = "/"
		req.proto = "HTTP/1.1"
	}

	req.parseHeaders(headers)

	return &req
}
