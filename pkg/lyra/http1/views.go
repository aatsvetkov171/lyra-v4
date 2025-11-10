package http1

import (
	"strings"
)

func NotFound(request *Request) *Response {
	response := NewResponse(404)
	response.AddString("<h1>Not Found</h1>")
	return response
}

func MethodNotAllowed(request *Request) *Response {
	response := NewResponse(405)
	response.AddString("<h1>Method Not Allowed</h1>")
	return response
}

func SendStaticFile(request *Request) *Response {
	response := NewResponse(200)
	fileEnd := strings.Split(request.GetPath(), ".")
	if len(fileEnd) == 2 {
		switch fileEnd[1] {
		case "css":
			response.GetHeaders()["Content-Type"] = "text/css; charset=UTF-8"
		case "js":
			response.GetHeaders()["Content-Type"] = "text/javascript; charset=UTF-8"
		}
	}
	response.AddFile(request.GetPath()[8:])
	return response
}
