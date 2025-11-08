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
	ConnTimeout        time.Duration
	MaxConnMesgCount   int
	ReqContentLenLimit [2]int
	templateDir        string
	DEBUG              bool
	BuferSizeFile      int
}

func NewConfig(addr string) *Config {
	newConfig := Config{
		Addr:               addr,
		Network:            "tcp",
		KeepAlive:          true,
		ReadTimeout:        10 * time.Second,
		WriteTimeout:       10 * time.Second,
		ConnTimeout:        10 * time.Second,
		MaxConnMesgCount:   100,
		ReqContentLenLimit: [2]int{0, 0},
		templateDir:        "templates\\",
		DEBUG:              true,
		BuferSizeFile:      4096,
	}
	return &newConfig
}

type lyra struct {
	Name   string
	config Config
	router *http1.Router
}

func NewServer(conf *Config, router *http1.Router) *lyra {
	newLyra := lyra{
		Name:   "Lyra-v4",
		config: *conf,
		router: router,
	}
	return &newLyra
}

func (l *lyra) ListenAdnServ() {
	listener, err := net.Listen(l.config.Network, l.config.Addr)
	if err != nil {
		fmt.Println("create listener error:", err.Error())
	}
	defer listener.Close()

	fmt.Println("Lyra listening on", l.config.Addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err.Error())
		}
		go l.connHandle(conn, l.router)
	}
}
