package main

import (
	"github.com/aatsvetkov171/lyra-v4/pkg/lyra/http1"
	"github.com/aatsvetkov171/lyra-v4/pkg/lyra/server"
)

func main() {

	config := server.NewConfig("localhost:8000")
	router := http1.NewRouter()
	server := server.NewServer(config, router)
	server.ListenAdnServ()

}
