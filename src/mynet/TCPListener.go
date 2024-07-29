package mynet

import (
	"net"
)

type TCPListener struct {
	listener *net.TCPListener
}

func (this *TCPListener) Addr() net.Addr {
	return this.listener.Addr()
}

func (this *TCPListener) ListenTCP() error {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	var err error
	this.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	return nil
}

func (this *TCPListener) AcceptTCP() (*net.TCPConn, error) {
	conn, err := this.listener.AcceptTCP()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
