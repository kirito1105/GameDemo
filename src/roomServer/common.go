package roomServer

func AddHeader(bytes []byte) []byte {
	size := len(bytes)
	buf := make([]byte, 0)
	buf = append(buf, byte(size), byte(size>>8), byte(size>>16), byte(size>>24))
	buf = append(buf, bytes...)
	return buf
}
