package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func readFirsLine(reader *bufio.Reader) ([]byte, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return line, err
	}
	return line, nil
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

func readReqBody(reader *bufio.Reader, c int) ([]byte, error) {
	if c == 0 {
		return []byte(""), nil
	}
	buffer := make([]byte, c)
	_, err := reader.Read(buffer)
	return buffer, err
}
