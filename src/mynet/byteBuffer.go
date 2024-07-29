package mynet

const (
	SIZE = 10
)

type ByteBuffer struct {
	_buffer     []byte
	_readIndex  int
	_writeIndex int
}

func NewByteBuffer() *ByteBuffer {
	return &ByteBuffer{
		_buffer:     make([]byte, SIZE),
		_readIndex:  0,
		_writeIndex: 0,
	}
}

func (this *ByteBuffer) RdBuf() []byte {
	return this._buffer[this._readIndex:this._writeIndex]
}

func (this *ByteBuffer) RdSize() int {
	return this._writeIndex - this._readIndex
}

func (this *ByteBuffer) RdFlip(step int) {
	if step < this.RdSize() {
		this._readIndex += step
	} else {
		this.ReSet()
	}
}

func (this *ByteBuffer) ReSet() {
	this._writeIndex = 0
	this._readIndex = 0
}

func (this *ByteBuffer) WrSize() int {
	return len(this._buffer) - this._writeIndex
}

func (this *ByteBuffer) WrBuf() []byte {
	return this._buffer[this._writeIndex:]
}

func (this *ByteBuffer) WrFlip(step int) {
	this._writeIndex += step
}

func (this *ByteBuffer) WrInc(size int) {
	if size > this.WrSize() {
		this.ReBuff(size)
	}
}

func (this *ByteBuffer) ReBuff(size int) {
	if this.WrSize()+this.RdSize() > size {
		tmpBuf := make([]byte, this._writeIndex+size)
		copy(tmpBuf, this._buffer)
		this._buffer = tmpBuf
	} else {
		step := this.RdSize()
		copy(this._buffer, this._buffer[this._readIndex:this._writeIndex])
		this._readIndex = 0
		this._writeIndex = this._readIndex + step
	}
}

func (this *ByteBuffer) Append(buf ...byte) {
	size := len(buf)
	if size == 0 {
		return
	}
	this.WrInc(size)
	copy(this.WrBuf(), buf)
	this.WrFlip(size)
}
