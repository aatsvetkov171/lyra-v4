package http1

import (
	"bytes"
	"sync"
)

var ResponseBuf = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

func GetResponseBuf() *bytes.Buffer {
	buf := ResponseBuf.Get().(*bytes.Buffer)
	return buf
}

func PutResponseBuf(buf *bytes.Buffer) {
	buf.Reset()
	ResponseBuf.Put(buf)
}
