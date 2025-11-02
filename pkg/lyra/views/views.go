package views

import "github.com/aatsvetkov171/lyra-v4/pkg/lyra/http1"

func NotFound(request *http1.Request) *http1.Response

func MethodNotAllowed(request *http1.Request) *http1.Response
