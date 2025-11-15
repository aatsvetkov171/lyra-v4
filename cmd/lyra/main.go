package main

import (
	"github.com/aatsvetkov171/lyra-v4/pkg/lyra/http1"
	"github.com/aatsvetkov171/lyra-v4/pkg/lyra/server"
)

func main() {

	config := server.NewConfig("localhost:8000")
	config.BuferSizeFile = 1024
	config.MaxConnMesgCount = 50
	router := http1.NewRouter(config.Path.StaticDir, config.Path.MediaDir)

	router.GET("/", IndexPage)
	router.GET("/ab", AboutPage)

	server := server.NewServer(config, router)
	server.ListenAdnServ()

}
