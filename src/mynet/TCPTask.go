package mynet

import (
	"fmt"
	"io"
	"net"
	"sync"
)

const HEADER_SIZE = 4

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
		if tolalSize < HEADER_SIZE {

			neednum := HEADER_SIZE - tolalSize

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

		HeadBuf := msgBuf[0:HEADER_SIZE]
		var Head = int(HeadBuf[0]) + int(HeadBuf[1])<<8 + int(HeadBuf[2])<<16 + int(HeadBuf[3])<<24
		fmt.Println(Head)
		if tolalSize < Head+HEADER_SIZE {
			neednum := Head + HEADER_SIZE - tolalSize
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

		this.Task.ParseMsg(msgBuf[HEADER_SIZE : HEADER_SIZE+Head])

		this.recvBuf.RdFlip(Head + HEADER_SIZE)

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
