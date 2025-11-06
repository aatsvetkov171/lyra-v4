package main

import (
	"time"

	"github.com/aatsvetkov171/lyra-v4/pkg/lyra/server"
)

func main() {

	var config server.Config = server.Config{
		Addr:               "localhost:8004",
		Network:            "tcp",
		ReadTimeout:        time.Second,
		WriteTimeout:       time.Second,
		MaxConnMesgCount:   60,
		MaxConnTime:        10,
		KeepAlive:          true,
		ReqContentLenLimit: [2]int{0, 0},
	}

	server := server.NewServer(&config)
	server.ListenAdnServ()
}
