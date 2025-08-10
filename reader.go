package go_framer

import (
	"encoding/binary"
	"io"
)

type Reader struct {
	r       io.Reader
	lenBuff []byte
	pos     int
	end     int
}

func NewReader(r io.Reader, lenBuff []byte) (*Reader, error) {
	if r == nil {
		return nil, emptyReaderError
	}

	if lenBuff == nil {
		lenBuff = make([]byte, MaxFrameSizeBytes)
	} else if len(lenBuff) < MaxFrameSizeBytes {
		return nil, lenBuffTooSmallError
	}

	return &Reader{
		r:       r,
		lenBuff: lenBuff,
		pos:     0,
		end:     0,
	}, nil
}

// TODO add docs
func (r *Reader) Read(b []byte) (int, error) {
	var err error = nil

	mLen, mLenPos, wc := 0, 0, 0

	for {
		if r.pos == r.end {
			r.pos = wc
		}

		if r.pos == wc {
			rn, re := r.r.Read(b[wc:])
			r.end = wc + rn
			err = re

			if err != nil {
				break
			}

			if rn == 0 {
				return 0, dstTooSmallError
			}
		}

		if mLen == 0 {
			copyLen := min(r.end-r.pos, MaxFrameSizeBytes-mLenPos)

			copy(r.lenBuff[mLenPos:], b[r.pos:r.pos+copyLen])

			r.pos += copyLen
			mLenPos += copyLen

			if mLenPos != MaxFrameSizeBytes {
				continue
			}

			mLenPos = 0
			mLen = int(binary.BigEndian.Uint32(r.lenBuff))

			if mLen == 0 {
				continue
			}

			if len(b) < mLen {
				err = dstTooSmallError

				break
			}

			if r.pos == r.end {
				continue
			}
		}

		copyLen := min(r.end-r.pos, mLen-wc)

		copy(b[wc:], b[r.pos:r.pos+copyLen])

		r.pos += copyLen
		wc += copyLen

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
