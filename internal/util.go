package internal

import "bytes"

func ReadString(reader *bytes.Reader) string {
	var buffer []byte
	for {
		buf := make([]byte, 1)
		reader.Read(buf)
		if buf[0] == 0x0 {
			break
		}
		buffer = append(buffer, buf[0])
	}
	return string(buffer)
}
