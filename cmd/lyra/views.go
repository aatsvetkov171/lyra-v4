package main

import (
	"fmt"

	"github.com/aatsvetkov171/lyra-v4/pkg/lyra/http1"
)

func IndexPage(request *http1.Request) *http1.Response {
	response := http1.NewResponse(200)
	response.AddFile("index.html")
	body := fmt.Sprintf("Метод: %s<br>Путь: %s<br> Заголовки:<br>",
		request.GetMethod(), request.GetPath())
	for key, val := range request.GetHeaders() {
		body += key + " : " + val + "<br>"
	}
	response.SetParams(map[string]string{
		"res": body,
	})
	return response
}
