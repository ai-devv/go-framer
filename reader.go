package go_framer

import (
	"encoding/binary"
	"io"
)

type Reader struct {
	r   io.Reader
	buf []byte
	pos int
	rc  int
}

func NewReader(r io.Reader, buf []byte) *Reader {
	return &Reader{
		r:   r,
		buf: buf,
		pos: 0,
		rc:  0,
	}
}

// TODO add docs
func (r *Reader) Read(b []byte) (int, error) {
	if len(r.buf) < 2 {
		return 0, bufferTooSmallError
	}

	var err error = nil
	var wc, mLen = 0, 0

	for {
		if r.pos == 0 {
			r.rc, err = r.r.Read(r.buf)

			if err != nil {
				break
			}
		}

		if mLen == 0 {
			if r.rc < 2 {
				err = srcTooSmallError

				break
			}

			mLen = int(binary.BigEndian.Uint16(r.buf[r.pos : r.pos+2]))

			if mLen == 0 {
				r.pos = 0

				continue
			}

			r.pos += 2

			if len(b) < mLen {
				err = dstTooSmallError

				break
			}
		}

		copyLen := min(r.rc-r.pos, mLen-wc)

		if copyLen == 0 {
			r.pos = 0

			continue
		}

		copy(b[wc:], r.buf[r.pos:r.pos+copyLen])

		r.pos += copyLen
		wc += copyLen

		if r.rc == r.pos {
			r.pos = 0
		}

		if wc == mLen {
			break
		}

		if wc > mLen {
			// TODO not covered by unit tests
			err = invalidReadResultError

			break
		}
	}

	return wc, err
}
