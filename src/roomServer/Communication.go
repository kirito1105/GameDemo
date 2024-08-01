package roomServer

//
//import (
//	"bufio"
//	"fmt"
//	"net"
//)
//
//type message struct{}
//
//type header struct {
//	len int
//}
//
//type roomRunable interface {
//}
//
//type Communication struct {
//	room   roomRunable
//	udpCli *net.UDPConn
//	tcpCli *net.TCPListener
//}
//
//func NewCommunication(room roomRunable) *Communication {
//	return &Communication{
//		room: room,
//	}
//}
//
//func (c *Communication) Listen() string {
//
//	for i := 0; i < 10; i++ {
//		addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
//		var err1 error
//		var err2 error
//		c.tcpCli, err2 = net.ListenTCP("tcp", addr)
//		udpaddr, _ := net.ResolveUDPAddr("udp", c.tcpCli.Addr().String())
//		c.udpCli, err1 = net.ListenUDP("udp", udpaddr)
//		if err1 == nil && err2 == nil {
//			return c.tcpCli.Addr().String()
//		}
//		if err1 == nil {
//			c.udpCli.Close()
//		}
//		if err2 == nil {
//			c.tcpCli.Close()
//		}
//	}
//	return ""
//}
//
//// 处理每个玩家的连接
//// 阻塞
//func (c *Communication) process(conn net.Conn) {
//	defer conn.Close()
//	for {
//		reader := bufio.NewReader(conn)
//		var buf [128]byte
//		n, err := reader.Read(buf[:]) // 读取数据
//		if err != nil {
//			fmt.Println("read from testClient failed, err:", err)
//			break
//		}
//		recvStr := string(buf[:n])
//		fmt.Println("收到client端发来的数据：", recvStr)
//		conn.Write([]byte(recvStr)) // 发送数据
//	}
//}
//
//func (c *Communication) Serve() { //监听端口
//	//TODO 目前只实现了TCP通信
//	for {
//		conn, err := c.tcpCli.Accept() // 阻塞建立连接
//		if err != nil {
//			return
//		}
//		go c.process(conn)
//	}
//}
