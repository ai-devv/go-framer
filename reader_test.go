package go_framer

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestNewReader(t *testing.T) {
	cases := []struct {
		name    string
		r       io.Reader
		lenBuff []byte
		ei      bool
		ee      error
	}{
		{
			name:    "full",
			r:       bytes.NewReader(nil),
			lenBuff: make([]byte, 2),
			ei:      true,
			ee:      nil,
		},
		{
			name:    "without len buff",
			r:       bytes.NewReader(nil),
			lenBuff: nil,
			ei:      true,
			ee:      nil,
		},
		{
			name:    "empty reader",
			r:       nil,
			lenBuff: nil,
			ei:      false,
			ee:      emptyReaderError,
		},
		{
			name:    "small len buff",
			r:       bytes.NewReader(nil),
			lenBuff: make([]byte, 0),
			ei:      false,
			ee:      lenBuffTooSmallError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r, err := NewReader(c.r, c.lenBuff)
			i := r != nil

			if c.ei != i {
				t.Errorf("expected instance: %v, got: %v\n", c.ei, i)
			}

			if !errors.Is(c.ee, err) {
				t.Errorf("expected err: %s, got: %s\n", c.ee, err)
			}
		})
	}
}

func TestReader_Read(t *testing.T) {
	cases := []struct {
		name    string
		m       []byte
		dstSize int
		en      int
		ee      error
		em      []byte
	}{
		// ======== GOOD / MESSAGES ========
		{
			name:    "very small one message",
			m:       []byte{0, 1, 1},
			dstSize: 10,
			en:      1,
			ee:      io.EOF,
			em:      []byte{1},
		},
		{
			name: "very small three messages",
			m: []byte{
				0, 1, 1,
				0, 1, 1,
				0, 1, 1,
			},
			dstSize: 10,
			en:      3,
			ee:      io.EOF,
			em:      []byte{1, 1, 1},
		},
		{
			name:    "small one message",
			m:       []byte{0, 3, 1, 2, 3},
			dstSize: 10,
			en:      3,
			ee:      io.EOF,
			em:      []byte{1, 2, 3},
		},
		{
			name:    "exactly one message",
			m:       []byte{0, 5, 1, 2, 3, 4, 5},
			dstSize: 5,
			en:      5,
			ee:      io.EOF,
			em:      []byte{1, 2, 3, 4, 5},
		},
		{
			name: "more than one message",
			m: []byte{
				0, 5, 1, 2, 3, 4, 5,
				0, 3, 6, 7, 8,
			},
			dstSize: 5,
			en:      8,
			ee:      io.EOF,
			em: []byte{
				1, 2, 3, 4, 5,
				6, 7, 8,
			},
		},
		{
			name: "big more than one message",
			m: []byte{
				0, 5, 1, 2, 3, 4, 5,
				0, 5, 6, 7, 8, 9, 0,
				0, 5, 1, 2, 3, 4, 5,
				0, 3, 6, 7, 8,
			},
			dstSize: 5,
			en:      18,
			ee:      io.EOF,
			em: []byte{
				1, 2, 3, 4, 5,
				6, 7, 8, 9, 0,
				1, 2, 3, 4, 5,
				6, 7, 8,
			},
		},
		{
			name:    "empty message",
			m:       nil,
			dstSize: 0,
			en:      0,
			ee:      io.EOF,
			em:      nil,
		},
		{
			name:    "zero message",
			m:       []byte{0, 0},
			dstSize: 1,
			en:      0,
			ee:      io.EOF,
			em:      nil,
		},
		{
			name: "zero message in stream",
			m: []byte{
				0, 1, 1,
				0, 0,
				0, 2, 3, 4,
			},
			dstSize: 2,
			en:      3,
			ee:      io.EOF,
			em: []byte{
				1,
				3, 4,
			},
		},
		// ======== GOOD / DESTINATION ========
		{
			name:    "small destination for one message",
			m:       []byte{0, 1, 1},
			dstSize: 1,
			en:      1,
			ee:      io.EOF,
			em:      []byte{1},
		},
		{
			name: "small destination for three messages",
			m: []byte{
				0, 1, 1,
				0, 1, 1,
				0, 1, 1,
			},
			dstSize: 1,
			en:      3,
			ee:      io.EOF,
			em: []byte{
				1,
				1,
				1,
			},
		},
		{
			name: "message size in end of destination",
			m: []byte{
				0, 5, 1, 2, 3, 4, 5,
				0, 3, 6, 7, 8,
			},
			dstSize: 9,
			en:      8,
			ee:      io.EOF,
			em: []byte{
				1, 2, 3, 4, 5,
				6, 7, 8,
			},
		},
		{
			name: "exploded message size",
			m: []byte{
				0, 5, 1, 2, 3, 4, 5,
				0, 3, 6, 7, 8,
			},
			dstSize: 8,
			en:      8,
			ee:      io.EOF,
			em: []byte{
				1, 2, 3, 4, 5,
				6, 7, 8,
			},
		},
		// ========== BAD ==========
		{
			name:    "too small destination for read",
			m:       []byte{0, 0},
			dstSize: 0,
			en:      0,
			ee:      dstTooSmallError,
			em:      nil,
		},
		{
			name:    "too small destination for handle message",
			m:       []byte{0, 5, 1, 2, 3, 4, 5},
			dstSize: 4,
			en:      0,
			ee:      dstTooSmallError,
			em:      nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var err error
			var am []byte

			b := make([]byte, c.dstSize)
			mr := bytes.NewReader(c.m)
			r, _ := NewReader(mr, nil)

			for {
				rn, re := r.Read(b)
				err = re
				am = append(am, b[:rn]...)

				if err != nil {
					break
				}
			}

			if c.en != len(am) {
				t.Errorf("expected readed: %d, got: %d\n", c.en, len(am))
			}

			if !errors.Is(c.ee, err) {
				t.Errorf("expected err: %s, got: %s\n", c.ee, err)
			}

			if !bytes.Equal(c.em, am) {
				t.Error("expected bytes not equal to actual bytes")
			}
		})
	}
}
