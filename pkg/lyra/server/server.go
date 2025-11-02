package server

import (
	"fmt"
	"net"
	"time"

	"github.com/aatsvetkov171/lyra-v4/pkg/lyra/http1"
)

type Config struct {
	Addr               string
	Network            string
	KeepAlive          bool
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	MaxConnTime        int
	MaxConnMesgCount   int
	ReqContentLenLimit [2]int
}

type lyra struct {
	Name   string
	config Config
	//router
}

func NewServer(conf *Config) *lyra {
	newLyra := lyra{
		Name:   "Lyra-v4",
		config: *conf,
	}
	return &newLyra
}

func (l *lyra) ListenAdnServ() {
	listener, err := net.Listen(l.config.Network, l.config.Addr)
	if err != nil {
		fmt.Println("create listener error:", err.Error())
	}
	defer listener.Close()

	router := http1.NewRouter()
	router.Handle("GET", "/", http1.Hello)
	fmt.Println("Lyra listening on", l.config.Addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err.Error())
		}
		go l.connHandle(conn, router)
	}
}
