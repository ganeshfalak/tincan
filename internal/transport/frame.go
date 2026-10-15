package transport

import (
	"encoding/binary"
	"fmt"
	"io"
)

const MaxFrame = 1 << 20 // 1 MiB

func WriteFrame(w io.Writer, payload []byte) error {
	n := len(payload)
	if n == 0 || n > MaxFrame {
		return fmt.Errorf("frame length %d", n)
	}

	var hdr [4]byte
	binary.BigEndian.PutUint32(hdr[:], uint32(n))
	if err := writeFull(w, hdr[:]); err != nil {
		return err
	}

	return writeFull(w, payload)
}

func ReadFrame(r io.Reader) ([]byte, error) {
	var hdr [4]byte
	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint32(hdr[:])
	if n == 0 || n > MaxFrame {
		return nil, fmt.Errorf("frame length %d", n)
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func writeFull(w io.Writer, p []byte) error {
	for len(p) > 0 {
		n, err := w.Write(p)
		p = p[n:]
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
	}
	return nil
}
