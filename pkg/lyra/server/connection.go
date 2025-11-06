package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/aatsvetkov171/lyra-v4/pkg/lyra/http1"
)

func newReader(conn net.Conn) *bufio.Reader {
	reader := bufio.NewReader(conn)
	return reader
}

func newWriter(conn net.Conn) *bufio.Writer {
	writer := bufio.NewWriter(conn)
	return writer
}

func isBlank(fline []byte) bool {
	return len(fline) == 0
}

func keepAliveTimer(conn net.Conn, timeout time.Duration, activeCh, doneCh chan struct{}) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case <-activeCh:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(timeout)
		case <-timer.C:
			fmt.Println("keep-alive timeout, conn closed")
			conn.Close()
			close(doneCh)
			return
		case <-doneCh:
			return
		}
	}
}

func (l *lyra) connHandle(conn net.Conn, router *http1.Router) {
	defer func() {
		conn.Close()
	}()
	reader := newReader(conn)
	messageCount := 0
	doneCh := make(chan struct{})
	activeCh := make(chan struct{})

	connDeadTimeout := l.config.ConnTimeout + (5 * time.Second)

	go keepAliveTimer(conn, l.config.ConnTimeout, activeCh, doneCh)

	for {
		conn.SetReadDeadline(time.Now().Add(connDeadTimeout))
		//--------------------------------------------
		fLine, err := readFirsLine(reader)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("SetReadDeadline timeout")
				return
			}
			if err == io.EOF {
				return
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			fmt.Println("some reading error", err.Error())
			return
		}
		if isBlank(fLine) {
			return
		}
		//--------------------------
		headersLines, err := readHeadersLines(reader)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("неизвестная ошибка чтения строк:", err.Error())
		}

		contentLen := FindReqContentLength(headersLines)

		if l.config.ReqContentLenLimit[0] != 0 {
			if contentLen > l.config.ReqContentLenLimit[1] {
				fmt.Println("conten len > max")
				return
			}
		}
		body, err := readReqBody(reader, contentLen)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("vontent length error:::::", err.Error())
			return
		}

		request := http1.NewRequest(fLine, headersLines, body)
		if connVal, ok := request.GetHeaders()["connection"]; ok && connVal == "close" {
			l.config.KeepAlive = false
		}
		for k, v := range request.GetHeaders() {
			fmt.Println(k, "---", v)
		}

		messageCount += 1

		writer := newWriter(conn)

		_, responseFunc := router.GetResponse(request)

		response := responseFunc(request)

		writer.Write(response.Build())

		writer.Flush()
		select {
		case activeCh <- struct{}{}:
		default:
		}
		if !l.config.KeepAlive || messageCount >= l.config.MaxConnMesgCount {
			break
		}
	}
}
