package go_framer

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestReader_Read(t *testing.T) {
	cases := []struct {
		name    string
		bufSize int
		m       []byte
		en      int
		ee      error
		em      []byte
	}{
		{
			name:    "very small",
			bufSize: 10,
			m:       []byte{0, 1, 1},
			en:      1,
			ee:      io.EOF,
			em:      []byte{1},
		},
		{
			name:    "small",
			bufSize: 10,
			m:       []byte{0, 5, 1, 2, 3, 4, 5},
			en:      5,
			ee:      io.EOF,
			em:      []byte{1, 2, 3, 4, 5},
		},
		{
			name:    "exactly one",
			bufSize: 7,
			m:       []byte{0, 5, 1, 2, 3, 4, 5},
			en:      5,
			ee:      io.EOF,
			em:      []byte{1, 2, 3, 4, 5},
		},
		{
			name:    "big",
			bufSize: 7,
			m: []byte{
				0, 5, 1, 2, 3, 4, 5,
				0, 3, 6, 7, 8,
			},
			en: 8,
			ee: io.EOF,
			em: []byte{
				1, 2, 3, 4, 5,
				6, 7, 8,
			},
		},
		{
			name:    "very big",
			bufSize: 7,
			m: []byte{
				0, 5, 1, 2, 3, 4, 5,
				0, 5, 6, 7, 8, 9, 0,
				0, 5, 1, 2, 3, 4, 5,
				0, 3, 6, 7, 8,
			},
			en: 18,
			ee: io.EOF,
			em: []byte{
				1, 2, 3, 4, 5,
				6, 7, 8, 9, 0,
				1, 2, 3, 4, 5,
				6, 7, 8,
			},
		},
		{
			name:    "empty message",
			bufSize: 2,
			m:       nil,
			en:      0,
			ee:      io.EOF,
			em:      nil,
		},
		{
			name:    "zero message",
			bufSize: 2,
			m:       []byte{0, 0},
			en:      0,
			ee:      io.EOF,
			em:      nil,
		},
		{
			name:    "small buffer",
			bufSize: 0,
			m:       nil,
			en:      0,
			ee:      bufferTooSmallError,
			em:      nil,
		},
		{
			name:    "too small source",
			bufSize: 2,
			m:       []byte{0},
			en:      0,
			ee:      srcTooSmallError,
			em:      nil,
		},
		{
			name:    "too small destination",
			bufSize: 103,
			m:       append([]byte{0, 101}, make([]byte, 101)...),
			en:      0,
			ee:      dstTooSmallError,
			em:      nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var err error

			buf := make([]byte, 100)
			rc := 0
			mr := bytes.NewReader(c.m)
			r := NewReader(mr, make([]byte, c.bufSize))

			for {
				rn, re := r.Read(buf[rc:])
				rc += rn
				err = re

				if err != nil {
					break
				}
			}

			if c.en != rc {
				t.Errorf("expected readed: %d, got: %d\n", c.en, rc)
			}

			if !errors.Is(c.ee, err) {
				t.Errorf("expected err: %s, got: %s\n", c.ee, err)
			}

			if !bytes.Equal(c.em, buf[:rc]) {
				t.Error("expected bytes not equal to actual bytes")
			}
		})
	}
}
