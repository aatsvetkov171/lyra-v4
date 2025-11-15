package server

import (
	"fmt"
	"net"
	"time"

	"github.com/aatsvetkov171/lyra-v4/pkg/lyra/http1"
)

type Paths struct {
	TemplateDir string
	StaticDir   string
	MediaDir    string
}

type Config struct {
	Addr               string
	Network            string
	KeepAlive          bool
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	ConnTimeout        time.Duration
	MaxConnMesgCount   int
	ReqContentLenLimit [2]int
	Path               *Paths
	DEBUG              bool
	BuferSizeFile      int
	ConnWorkerCount    int
	ConnTasksBuf       int
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
		Path: &Paths{
			TemplateDir: "templates",
			StaticDir:   "static",
			MediaDir:    "media",
		},
		DEBUG:           true,
		BuferSizeFile:   4096,
		ConnWorkerCount: 16,
		ConnTasksBuf:    500,
	}
	return &newConfig
}

type lyra struct {
	Name        string
	config      Config
	router      *http1.Router
	logger      *LyraLog
	connTasksCh chan net.Conn
}

func NewServer(conf *Config, router *http1.Router) *lyra {
	newLyra := lyra{
		Name:        "Lyra-v4",
		config:      *conf,
		router:      router,
		logger:      NewLyraLog(),
		connTasksCh: make(chan net.Conn, conf.ConnTasksBuf),
	}
	return &newLyra
}

func (l *lyra) SetLogDebug(flag bool) {
	l.logger.debug = flag
}

func (l *lyra) worker(id int, connTasks <-chan net.Conn, router *http1.Router) {
	for conn := range connTasks {
		l.connHandle(conn, router)
	}
}

func (l *lyra) ListenAdnServ() {
	listener, err := net.Listen(l.config.Network, l.config.Addr)
	if err != nil {
		fmt.Println("create listener error:", err.Error())
	}
	defer listener.Close()

	for i := 0; i < l.config.ConnWorkerCount; i++ {
		go l.worker(i, l.connTasksCh, l.router)
	}
	l.logger.Info("listening on %s", l.config.Addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err.Error())
		}
		select {
		case l.connTasksCh <- conn:
			// ok
		default:
			conn.Close()
		}
	}
}
