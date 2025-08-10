package go_framer

import (
	"bytes"
	"errors"
	"io"
	"math"
	"testing"
)

func TestNewWriter(t *testing.T) {
	cases := []struct {
		name      string
		w         io.Writer
		frameSize int
		lenBuff   []byte
		ee        error
	}{
		{
			name:      "full",
			w:         io.Discard,
			frameSize: 1,
			lenBuff:   make([]byte, 2),
			ee:        nil,
		},
		{
			name:      "too small frame size",
			w:         io.Discard,
			frameSize: 0,
			lenBuff:   make([]byte, 2),
			ee:        invalidFrameSize,
		},
		{
			name:      "too large frame size",
			w:         io.Discard,
			frameSize: math.MaxUint16 + 1,
			lenBuff:   make([]byte, 2),
			ee:        invalidFrameSize,
		},
		{
			name:      "too small len buff",
			w:         io.Discard,
			frameSize: 1,
			lenBuff:   make([]byte, 0),
			ee:        lenBuffTooSmallError,
		},
		{
			name:      "without len buff",
			w:         io.Discard,
			frameSize: 1,
			lenBuff:   nil,
			ee:        nil,
		},
		{
			name:      "empty writer",
			w:         nil,
			frameSize: 1,
			lenBuff:   make([]byte, 2),
			ee:        emptyWriterError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := NewWriter(c.w, c.frameSize, c.lenBuff)

			if !errors.Is(c.ee, err) {
				t.Errorf("expected err: %s, got: %s\n", c.ee, err)
			}
		})
	}
}

func TestWriter_Write(t *testing.T) {
	cases := []struct {
		name      string
		m         []byte
		frameSize int
		en        int
		ee        error
		em        []byte
	}{
		// ======== GOOD / MESSAGES ========
		{
			name:      "very small one message",
			m:         []byte{1},
			frameSize: 1,
			en:        1,
			ee:        nil,
			em:        []byte{0, 1, 1},
		},
		{
			name: "very small three messages",
			m: []byte{
				1,
				1,
				1,
			},
			frameSize: 1,
			en:        3,
			ee:        nil,
			em: []byte{
				0, 1, 1,
				0, 1, 1,
				0, 1, 1,
			},
		},
		{
			name:      "small one message",
			m:         []byte{1, 2, 3},
			frameSize: 10,
			en:        3,
			ee:        nil,
			em:        []byte{0, 3, 1, 2, 3},
		},
		{
			name:      "exactly one message",
			m:         []byte{1, 2, 3, 4, 5},
			frameSize: 5,
			en:        5,
			ee:        nil,
			em:        []byte{0, 5, 1, 2, 3, 4, 5},
		},
		{
			name: "more than one message",
			m: []byte{
				1, 2, 3, 4, 5,
				6, 7, 8,
			},
			frameSize: 5,
			en:        8,
			ee:        nil,
			em: []byte{
				0, 5, 1, 2, 3, 4, 5,
				0, 3, 6, 7, 8,
			},
		},
		{
			name: "big more than one message",
			m: []byte{
				1, 2, 3, 4, 5,
				6, 7, 8, 9, 0,
				1, 2, 3, 4, 5,
				6, 7, 8,
			},
			frameSize: 5,
			en:        18,
			ee:        nil,
			em: []byte{
				0, 5, 1, 2, 3, 4, 5,
				0, 5, 6, 7, 8, 9, 0,
				0, 5, 1, 2, 3, 4, 5,
				0, 3, 6, 7, 8,
			},
		},
		{
			name:      "empty message",
			m:         nil,
			frameSize: 1,
			en:        0,
			ee:        nil,
			em:        nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer

			w, _ := NewWriter(&buf, c.frameSize, nil)
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
