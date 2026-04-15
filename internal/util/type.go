package util

import "io"

type BytesReader []byte

func (b BytesReader) Read(p []byte) (int, error) {
	if len(b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, b)
	b = b[n:]
	return n, nil
}
