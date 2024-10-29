package mitm

import (
	"bufio"
	"net"
)

type BufferIoConn struct {
	net.Conn
	reader *bufio.Reader
}

func NewBufferIoConn(conn net.Conn, reader *bufio.Reader) *BufferIoConn {
	return &BufferIoConn{
		Conn:   conn,
		reader: reader,
	}
}

func (conn *BufferIoConn) Read(p []byte) (n int, err error) {
	return conn.reader.Read(p)
}
