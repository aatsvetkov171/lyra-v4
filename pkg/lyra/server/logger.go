package server

import (
	"fmt"
	"log"
	"os"
)

type LyraLog struct {
	log   *log.Logger
	debug bool
}

func NewLyraLog() *LyraLog {
	newLyraLog := LyraLog{
		log:   log.New(os.Stdout, "[Lyra] ", log.Ldate|log.Ltime),
		debug: true,
	}
	return &newLyraLog
}

func (l *LyraLog) Info(msg string, args ...any) {
	l.log.Println("INFO:", fmt.Sprintf(msg, args...))
}

func (l *LyraLog) Debug(msg string, args ...any) {
	if l.debug {
		l.log.Println("DEBUG:", fmt.Sprintf(msg, args...))
	}
}
