package http1

import (
	"strings"
)

type Request struct {
	method  string
	path    string
	proto   string
	headers map[string]string
	body    string
}

func NewRequest(Fline []byte, lines []byte, body []byte) *Request {
	firstLineString := strings.TrimSpace(string(Fline))
	firstLineStringSlice := strings.Fields(firstLineString)
	headers := make(map[string]string)
	linesString := strings.TrimSpace(string(lines))
	linesStringSlice := strings.Split(linesString, "\r\n")

	for i := 0; i < len(linesStringSlice); i++ {
		header := strings.SplitN(linesStringSlice[i], ":", 2)
		headers[strings.ToLower(header[0])] = header[1]

	}

	req := Request{
		method:  firstLineStringSlice[0],
		path:    firstLineStringSlice[1],
		proto:   firstLineStringSlice[2],
		headers: headers,
		body:    string(body),
	}
	return &req
}
