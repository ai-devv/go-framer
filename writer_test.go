package go_framer

import (
	"bytes"
	"errors"
	"math"
	"testing"
)

func TestWriter_Write(t *testing.T) {
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
			m:       []byte{1},
			en:      1,
			ee:      nil,
			em:      []byte{0, 1, 1},
		},
		{
			name:    "small",
			bufSize: 10,
			m:       []byte{1, 2, 3, 4, 5},
			en:      5,
			ee:      nil,
			em:      []byte{0, 5, 1, 2, 3, 4, 5},
		},
		{
			name:    "exactly one",
			bufSize: 7,
			m:       []byte{1, 2, 3, 4, 5},
			en:      5,
			ee:      nil,
			em:      []byte{0, 5, 1, 2, 3, 4, 5},
		},
		{
			name:    "big",
			bufSize: 7,
			m: []byte{
				1, 2, 3, 4, 5,
				6, 7, 8,
			},
			en: 8,
			ee: nil,
			em: []byte{
				0, 5, 1, 2, 3, 4, 5,
				0, 3, 6, 7, 8,
			},
		},
		{
			name:    "very big",
			bufSize: 7,
			m: []byte{
				1, 2, 3, 4, 5,
				6, 7, 8, 9, 0,
				1, 2, 3, 4, 5,
				6, 7, 8,
			},
			en: 18,
			ee: nil,
			em: []byte{
				0, 5, 1, 2, 3, 4, 5,
				0, 5, 6, 7, 8, 9, 0,
				0, 5, 1, 2, 3, 4, 5,
				0, 3, 6, 7, 8,
			},
		},
		{
			name:    "empty message",
			bufSize: 3,
			m:       nil,
			en:      0,
			ee:      nil,
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
			name:    "too large source",
			bufSize: math.MaxUint16 + 3,
			m:       make([]byte, math.MaxUint16+1),
			en:      0,
			ee:      srcTooLargeError,
			em:      nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer

			w := NewWriter(&buf, make([]byte, c.bufSize))
			n, err := w.Write(c.m)

			if c.en != n {
				t.Errorf("expected written: %d, got: %d\n", c.en, n)
			}

			if !errors.Is(c.ee, err) {
				t.Errorf("expected err: %s, got: %s\n", c.ee, err)
			}

			if !bytes.Equal(c.em, buf.Bytes()) {
				t.Error("expected bytes not equal to actual bytes")
			}
		})
	}
}
