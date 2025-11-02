package http1

func NotFound(request *Request) *Response {
	response := NewResponse()
	response.AddBody([]byte("<p>Lyyyyraaa</p>"))
	return response //потом исправлю на правильное поведение
}

func MethodNotAllowed(request *Request) *Response {
	response := NewResponse()
	response.AddBody([]byte("<p>Lyyyyraaa</p>"))
	return response //потом исправлю на правильное поведение
}

func Hello(request *Request) *Response {
	response := NewResponse()
	response.AddBody([]byte("<p>Lyyyyraaa Helloooo</p>"))
	return response
}
