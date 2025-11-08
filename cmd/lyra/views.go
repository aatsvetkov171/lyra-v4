package main

import "github.com/aatsvetkov171/lyra-v4/pkg/lyra/http1"

func IndexPage(request *http1.Request) *http1.Response {
	response := http1.NewResponse(200)
	response.AddFile("index.html")
	return response
}
