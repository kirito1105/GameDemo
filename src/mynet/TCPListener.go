package mynet

import (
	"net"
)

const URL = "192.168.96.182:0"

type TCPListener struct {
	listener *net.TCPListener
}

func (this *TCPListener) Addr() net.Addr {
	return this.listener.Addr()
}

func (this *TCPListener) ListenTCP() error {
	addr, _ := net.ResolveTCPAddr("tcp", URL)
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
	err = conn.SetNoDelay(true)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
