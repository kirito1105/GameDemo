package roomServer

import (
	"net"
)

type message struct{}

type roomRunable interface {
}

type Communication struct {
	udpCli *net.UDPConn
	tcpCli *net.TCPListener
}

func NewCommunication() *Communication {
	return &Communication{}
}

func (c *Communication) Listen() string {

	for i := 0; i < 10; i++ {
		addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
		var err1 error
		var err2 error
		c.tcpCli, err2 = net.ListenTCP("tcp", addr)
		udpaddr, _ := net.ResolveUDPAddr("udp", c.tcpCli.Addr().String())
		c.udpCli, err1 = net.ListenUDP("udp", udpaddr)
		if err1 == nil && err2 == nil {
			return c.tcpCli.Addr().String()
		}
		if err1 == nil {
			c.udpCli.Close()
		}
		if err2 == nil {
			c.tcpCli.Close()
		}
	}
	return ""
}

// 处理每个玩家的连接
// 阻塞
func (c *Communication) process(conn net.Conn) {
	defer conn.Close()

}

func (c *Communication) Serve() { //监听端口
	//TODO 目前只实现了TCP通信
	for {
		conn, err := c.tcpCli.Accept() // 阻塞建立连接
		if err != nil {
			return
		}
		go c.process(conn)
	}
}
