package go_framer

import (
	"encoding/binary"
	"io"
	"sync"
)

type Writer struct {
	mu sync.Mutex

	w         io.Writer
	frameSize int
	lenBuf    []byte
}

func NewWriter(w io.Writer, frameSize int, lenBuff []byte) (*Writer, error) {
	if w == nil {
		return nil, emptyWriterError
	}

	if frameSize < 1 || frameSize > MaxFrameSize {
		return nil, invalidFrameSize
	}

	if lenBuff == nil {
		lenBuff = make([]byte, MaxFrameSizeBytes)
	} else if len(lenBuff) < MaxFrameSizeBytes {
		return nil, lenBuffTooSmallError
	}

	return &Writer{
		w:         w,
		frameSize: frameSize,
		lenBuf:    lenBuff,
	}, nil
}

// TODO add docs
func (w *Writer) Write(b []byte) (int, error) {
	bLen := len(b)

	if bLen == 0 {
		return 0, nil
	}

	var wl, wc = 0, 0
	var err error = nil

	w.mu.Lock()
	defer w.mu.Unlock()

	for {
		wl = min(bLen-wc, w.frameSize)

		binary.BigEndian.PutUint32(w.lenBuf, uint32(wl))

		wn, we := w.w.Write(w.lenBuf)
		err = we

		if err != nil {
			// TODO not covered by unit tests
			break
		}

		if wn != MaxFrameSizeBytes {
			// TODO not covered by unit tests
			err = invalidWriteResultError

			break
		}

		wn, we = w.w.Write(b[wc : wc+wl])
		wc += wn
		err = we

		if wc == bLen || err != nil {
			break
		}

		if wc > bLen {
			// TODO not covered by unit tests
			err = invalidWriteResultError

			break
		}
	}

	return wc, err
}
