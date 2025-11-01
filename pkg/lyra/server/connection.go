package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/aatsvetkov171/lyra-v4/pkg/lyra/http1"
)

func newReader(conn net.Conn) *bufio.Reader {
	reader := bufio.NewReader(conn)
	return reader
}

func FindReqContentLength(headers []byte) int {

	headers = bytes.ToLower(headers)
	index1 := bytes.Index(headers, []byte("content-length"))
	if index1 == -1 {
		return 0
	}
	headersIndex := headers[index1:]
	index2 := bytes.Index(headersIndex, []byte("\r\n")) + index1
	if index2 == -1 {
		return 0
	}
	res := headers[index1:index2]
	resBytes := bytes.TrimSpace(bytes.Split(res, []byte(":"))[1])
	resInt, err := strconv.Atoi(string(resBytes))
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	return resInt
}

func readFirsLine(reader *bufio.Reader) ([]byte, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return line, err
	}
	return line, nil
}

func isBlank(fline []byte) bool {
	return len(fline) == 0
}

func readHeadersLines(reader *bufio.Reader) ([]byte, error) {
	headers := make([]byte, 0)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("соединение закрыто")
				return headers, err
			}
			return headers, err
		}

		if len(line) == 2 && line[0] == '\r' && line[1] == '\n' {
			break
		}
		headers = append(headers, line...)
	}
	return headers, nil
}

func readReqBody(reader *bufio.Reader, c int) ([]byte, error) {
	if c == 0 {
		return []byte(""), nil
	}
	buffer := make([]byte, c)
	_, err := reader.Read(buffer)
	return buffer, err
}

func (l *lyra) connHandle(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	keepAlive := l.config.KeepAlive
	reader := newReader(conn)
	messageCount := 0
	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(l.config.MaxConnTime) * time.Second))
		//--------------------------------------------
		fLine, err := readFirsLine(reader)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Time out")
				return
			}
			if err == io.EOF {
				return
			}
			fmt.Println(err.Error())
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
			fmt.Println(err.Error())
			return
		}

		request := http1.NewRequest(fLine, headersLines, body)
		if connVal, ok := request.GetHeaders()["connection"]; ok && connVal == "close" {
			keepAlive = false
		}
		fmt.Println(request)

		messageCount += 1
		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length:5\r\nConnection:keep-alive\r\n\r\nhello"))
		if !keepAlive || messageCount >= l.config.MaxConnMesgCount {
			break
		}
	}
}
