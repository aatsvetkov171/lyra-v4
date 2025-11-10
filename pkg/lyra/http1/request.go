package http1

import (
	"net/url"
	"strings"
)

type Request struct {
	method   string
	path     string
	proto    string
	headers  map[string]string
	body     []byte
	DataPOST map[string]string
}

func (r *Request) parseHeaders(headers []byte) {
	headersStrings := strings.Split(string(headers), "\r\n")
	for i := 0; i < len(headersStrings); i++ {
		header := strings.SplitN(headersStrings[i], ":", 2)
		if len(header) < 2 {
			continue
		}

		r.headers[strings.TrimSpace(strings.ToLower(header[0]))] = strings.TrimSpace(header[1])
		//fmt.Println(r.headers[strings.ToLower(header[0])])
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
		req.method = "UNKNOWN"
		req.path = "/"
		req.proto = "HTTP/1.1"
	}

	req.parseHeaders(headers)

	if req.method == "POST" {
		req.createDataPOST()
	}

	return &req
}

// name=alex&fam=true
func parseQuery(body []byte) map[string]string {
	dict := make(map[string]string)
	if len(body) == 0 {
		return dict
	}
	slice1 := strings.Split(string(body), "&")
	for i := 0; i < len(slice1); i++ {
		parts := strings.SplitN(slice1[i], "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, err1 := url.QueryUnescape(parts[0])
		val, err2 := url.QueryUnescape(parts[1])
		if err1 == nil && err2 == nil {
			dict[key] = val
		}

	}
	return dict
}
func (r *Request) createDataPOST() {
	if r.method == "POST" {
		if r.headers["content-type"] == "application/x-www-form-urlencoded" {
			r.DataPOST = parseQuery(r.body)
		}
	}
}

func (r *Request) GetHeaders() map[string]string {
	return r.headers
}

func (r *Request) GetMethod() string {
	return r.method
}

func (r *Request) GetPath() string {
	return r.path
}

func (r *Request) GetBody() []byte {
	return r.body
}
