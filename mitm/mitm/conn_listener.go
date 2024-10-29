package mitm

import "net"

type ConnListener struct {
	connChan chan net.Conn
}

func NewConnListener() *ConnListener {
	return &ConnListener{
		connChan: make(chan net.Conn, 100),
	}
}

func (ln *ConnListener) Accept() (net.Conn, error) {
	conn := <-ln.connChan
	if conn == nil {
		return nil, net.ErrClosed
	}
	return conn, nil
}

func (ln *ConnListener) Close() error {
	ln.connChan <- nil
	return nil
}

func (ln *ConnListener) Addr() net.Addr {
	return &net.TCPAddr{}
}

func (ln *ConnListener) Add(conn net.Conn) {
	ln.connChan <- conn
}
