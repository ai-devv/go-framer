package go_framer

import "errors"

var emptyReaderError = errors.New("empty reader")

var lenBuffTooSmallError = errors.New("length buffer too small")

var bufferTooSmallError = errors.New("buffer too small")

var dstTooSmallError = errors.New("destination too small")

var invalidReadResultError = errors.New("invalid read result")

var invalidWriteResultError = errors.New("invalid write result")
