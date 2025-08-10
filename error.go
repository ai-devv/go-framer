package framer

import "errors"

var emptyReaderError = errors.New("empty reader")

var lenBuffTooSmallError = errors.New("length buffer too small")

var dstTooSmallError = errors.New("destination too small")

var invalidReadResultError = errors.New("invalid read result")

var emptyWriterError = errors.New("empty writer")

var invalidFrameSize = errors.New("invalid frame size")

var invalidWriteResultError = errors.New("invalid write result")
