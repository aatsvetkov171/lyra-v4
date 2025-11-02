package http1

func NotFound(request *Request) *Response

func MethodNotAllowed(request *Request) *Response

func Hello(request *Request) *Response {
	response := NewResponse()
	response.AddBody([]byte("<p>Lyyyyraaa</p>"))
	return response
}
