package go_framer

import (
	"encoding/binary"
	"io"
	"math"
)

type Writer struct {
	w   io.Writer
	buf []byte
}

func NewWriter(w io.Writer, buf []byte) *Writer {
	return &Writer{
		w:   w,
		buf: buf,
	}
}

// TODO add docs
func (w *Writer) Write(b []byte) (int, error) {
	bufCap := len(w.buf) - 2

	if bufCap < 1 {
		return 0, bufferTooSmallError
	}

	bLen := len(b)

	if bLen == 0 {
		return 0, nil
	}

	var wc = 0
	var err error = nil

	for {
		writeLen := min(bufCap, bLen-wc)

		if writeLen > math.MaxUint16 {
			err = srcTooLargeError

			break
		}

		binary.BigEndian.PutUint16(w.buf, uint16(writeLen))
		copy(w.buf[2:], b[wc:wc+writeLen])

		wn, we := w.w.Write(w.buf[:writeLen+2])
		wc += max(0, wn-2)
		err = we

		if wc == bLen || err != nil {
			break
		}

		if wn < writeLen+2 || wc > bLen {
			// TODO not covered by unit tests
			err = invalidWriteResultError

			break
		}
	}

	return wc, err
}
