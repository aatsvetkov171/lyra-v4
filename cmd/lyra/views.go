package main

import (
	"fmt"

	"github.com/aatsvetkov171/lyra-v4/pkg/lyra/http1"
)

func IndexPage(request *http1.Request) *http1.Response {
	response := http1.NewResponse(200)
	response.AddFile("index.html")
	return response
}

func SubmitForm(request *http1.Request) *http1.Response {
	fmt.Println(request.DataPOST)
	for k, v := range request.DataPOST {
		fmt.Println(k, "()", v)
	}
	res := http1.NewResponse(200)
	res.AddString("<h1>ДАнные получены</h1>")
	return res
}
