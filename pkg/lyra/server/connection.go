package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aatsvetkov171/lyra-v4/pkg/lyra/http1"
)

var readerPool = sync.Pool{
	New: func() any {
		return bufio.NewReader(nil)
	},
}

var writerPool = sync.Pool{
	New: func() any {
		return bufio.NewWriter(nil)
	},
}

func newReader(conn net.Conn) *bufio.Reader {
	r := readerPool.Get().(*bufio.Reader)
	r.Reset(conn)
	return r
}

func putReader(r *bufio.Reader) {
	r.Reset(nil)
	readerPool.Put(r)
}

func newWriter(conn net.Conn) *bufio.Writer {
	w := writerPool.Get().(*bufio.Writer)
	w.Reset(conn)
	return w
}

func putWriter(w *bufio.Writer) {
	w.Reset(nil)
	writerPool.Put(w)
}

func isBlank(fline []byte) bool {
	return len(fline) == 0
}

func getPathFile(filename string, templateDir string, debug bool) (string, error) {
	if debug {
		path, err := os.Getwd()

		return filepath.Join(path, templateDir, filename), err
	} else {
		path, err := os.Getwd()

		return filepath.Join(path, templateDir, filename), err
	}

	//strings.Count(exe, "/")

}

func sendFile(response *http1.Response, config *Config, writer *bufio.Writer) error {
	var path string
	var err error
	if response.GetMimeType() == "static" {
		path, err = getPathFile(response.GetFileName(), config.Path.StaticDir, config.DEBUG)
	} else if response.GetMimeType() == "media" {
		path, err = getPathFile(response.GetFileName(), config.Path.MediaDir, config.DEBUG)
	} else {
		path, err = getPathFile(response.GetFileName(), config.Path.TemplateDir, config.DEBUG)
	}
	if err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	size := int(info.Size())
	//fmt.Println(size)
	for k, v := range response.GetParams() {
		size = size - (4 + len(k)) + len(v)
	}
	//fmt.Println(size)

	response.AddHeader("Content-Length", strconv.Itoa(size))
	writer.Write(response.GetHeadersBytes())
	writer.Flush()
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	scanner := bufio.NewReader(file)

	buf := make([]byte, config.BuferSizeFile)
	for {
		n, err := scanner.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if n > 0 {
			chank := buf[:n]
			for key, val := range response.GetParams() {
				pattern := []byte("{{" + key + "}}")
				chank = bytes.ReplaceAll(chank, pattern, []byte(val))
			}
			writer.Write(chank)
		}
	}
	writer.Flush()
	file.Close()
	return nil
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
			conn.Close()
			close(doneCh)
			return
		case <-doneCh:
			return
		}
	}
}

func (l *lyra) connHandle(conn net.Conn, router *http1.Router) {

	l.logger.Info("new connection %s", conn.RemoteAddr().String())

	reader := newReader(conn)
	writer := newWriter(conn)

	defer func() {
		l.logger.Info("connection closed %s", conn.RemoteAddr().String())

		putReader(reader)
		putWriter(writer)
		conn.Close()

	}()

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
				l.logger.Debug("keep-alive timeout")
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

		l.logger.Debug("new request %s %s %s", request.GetMethod(), request.GetPath(), request.GetProto())
		if connVal, ok := request.GetHeaders()["connection"]; ok && connVal == "close" {
			l.config.KeepAlive = false
		}

		messageCount += 1

		response := router.GetResponse(request)
		l.logger.Debug("mime type %s", response.GetMimeType())

		if response.GetFileName() != "nofile" {
			err := sendFile(response, &l.config, writer)
			if err != nil {
				fmt.Println("file not found", err.Error())
				response = router.NotFound(request)
				writer.Write(response.GetHeadersBytes())
				writer.Write(response.GetBody())
				writer.Flush()
			}
		} else {
			writer.Write(response.GetHeadersBytes())
			writer.Write(response.GetBody())
			writer.Flush()
		}
		select {
		case activeCh <- struct{}{}:
		default:
		}
		if !l.config.KeepAlive || messageCount >= l.config.MaxConnMesgCount {
			break
		}
	}
}
