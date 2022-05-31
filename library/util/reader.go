package util

import (
	"io"
)

func ReadLimitBytes(r io.Reader, maxReadNum int) []byte {
	var yetReadNum  = 0
	b := make([]byte, 8) // 8 这里控制每次读取的字节数
	total := make([]byte, 0)
	for {
		n, err := r.Read(b)
		total = append(total, b[:n]...)
		if err == io.EOF {
			break
		}
		if maxReadNum > 0 {
			yetReadNum += n
			if yetReadNum >= maxReadNum {
				break
			}
		}
	}

	return total
}
