package go_framer

import "errors"

var bufferTooSmallError = errors.New("buffer too small")

var srcTooSmallError = errors.New("source too small")

var dstTooSmallError = errors.New("destination too small")

var invalidReadResultError = errors.New("invalid read result")

var srcTooLargeError = errors.New("source too large")

var invalidWriteResultError = errors.New("invalid write result")
