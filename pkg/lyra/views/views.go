package views

import "github.com/aatsvetkov171/lyra-v4/pkg/lyra/http1"

func NotFound(request *http1.Request) *http1.Response

func MethodNotAllowed(request *http1.Request) *http1.Response

func Hello(request *http1.Request) *http1.Response {
	response := http1.NewResponse()
	response.AddBody([]byte("<p>Lyyyyraaa</p>"))
	return &response
}
