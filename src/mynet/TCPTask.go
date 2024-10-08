package mynet

import (
	"io"
	"net"
	"sync"
	"time"
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
				return
			}

			tolalSize = this.recvBuf.RdSize()
		}

		msgBuf = this.recvBuf.RdBuf()

		HeadBuf := msgBuf[0:HEADER_SIZE]
		var Head = int(HeadBuf[0]) + int(HeadBuf[1])<<8 + int(HeadBuf[2])<<16 + int(HeadBuf[3])<<24
		if tolalSize < Head+HEADER_SIZE {
			neednum := Head + HEADER_SIZE - tolalSize
			err := this.readAtLeast(this.recvBuf, neednum)
			if err != nil {
				return
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
				if this.sendBuf.RdReady() {
					tmpBuf.Append(this.sendBuf.RdBuf()...)
					this.sendBuf.ReSet()
				}
				this.sendMux.Unlock()
				if !tmpBuf.RdReady() {
					break
				}
				num, err := this.Conn.Write(tmpBuf.RdBuf()[:tmpBuf.RdSize()])
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
	time.Sleep(0)
	buf.WrFlip(n)
	return err
}

func (this *TCPTask) SendMsg(msg []byte) {
	this.sendMux.Lock()
	this.sendBuf.Append(msg...)
	this.Sighal()
	defer this.sendMux.Unlock()

}

func (this *TCPTask) Sighal() {
	select {
	case this.sendSign <- struct{}{}:
	default:
	}
}
