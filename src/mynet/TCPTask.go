package mynet

import (
	"fmt"
	"io"
	"net"
	"sync"
	"unsafe"
)

type TCPTaskInter interface {
	ParseMsg(data []byte) bool
}

type TCPTask struct {
	Conn     net.Conn
	recvBuf  *ByteBuffer
	sendSign chan struct{}
	sendMux  sync.Mutex
	sendBuf  *ByteBuffer
	Task     TCPTaskInter
}

func NewTCPTask(conn net.Conn) *TCPTask {
	return &TCPTask{
		Conn:     conn,
		recvBuf:  NewByteBuffer(),
		sendBuf:  NewByteBuffer(),
		sendSign: make(chan struct{}),
	}
}

func (this *TCPTask) Start() {
	go this.sendloop()
	go this.recvloop()
}

func (this *TCPTask) Close() {
	this.Conn.Close()
	this.recvBuf.ReSet()
	this.sendBuf.ReSet()
}

func (this *TCPTask) recvloop() {
	defer this.Close()

	var (
		tolalSize int
		msgBuf    []byte
	)

	for {
		tolalSize = this.recvBuf.RdSize()
		if tolalSize < int(unsafe.Sizeof(new(TCPHeader))) {

			neednum := int(unsafe.Sizeof(new(TCPHeader))) - tolalSize

			err := this.readAtLeast(this.recvBuf, neednum)
			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
					return
				}
				continue
			}

			tolalSize = this.recvBuf.RdSize()
		}

		msgBuf = this.recvBuf.RdBuf()

		HeadBuf := msgBuf[0:int(unsafe.Sizeof(new(TCPHeader)))]
		Head := *(**TCPHeader)(unsafe.Pointer(&HeadBuf))
		if tolalSize < Head.Size+int(unsafe.Sizeof(new(TCPHeader))) {
			neednum := Head.Size + int(unsafe.Sizeof(new(TCPHeader))) - tolalSize
			err := this.readAtLeast(this.recvBuf, neednum)
			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
					return
				}
				continue
			}
			msgBuf = this.recvBuf.RdBuf()
		}

		this.Task.ParseMsg(msgBuf[int(unsafe.Sizeof(new(TCPHeader))) : int(unsafe.Sizeof(new(TCPHeader)))+Head.Size])

		this.recvBuf.RdFlip(Head.Size + int(unsafe.Sizeof(new(TCPHeader))))

	}
}

func (this *TCPTask) sendloop() {
	defer this.Close()

	var tmpBuf = NewByteBuffer()

	for {
		select {
		case <-this.sendSign:
			for {
				this.sendMux.Lock()
				tmpBuf.Append(this.sendBuf.RdBuf()...)
				this.sendBuf.ReSet()
				this.sendMux.Unlock()

				num, err := this.Conn.Write(tmpBuf.RdBuf())
				if err != nil {
					return
				}
				tmpBuf.RdFlip(num)
			}

		}
	}
}

func (this *TCPTask) readAtLeast(buf *ByteBuffer, neednum int) error {
	buf.WrInc(neednum)
	n, err := io.ReadAtLeast(this.Conn, buf.WrBuf(), neednum)
	buf.WrFlip(n)
	return err
}

func (this *TCPTask) SendMsg(msg []byte) {
	_, err := this.Conn.Write(msg)
	if err != nil {
		return
	}
}
