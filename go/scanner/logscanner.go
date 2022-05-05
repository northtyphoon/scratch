package scanner

import (
	"bufio"
	"io"
)

type LogScanner struct {
	reader *bufio.Reader
	bytes  []byte
	err    error
}

// NOTE: bufio.Reader supports minimal 16 bytes buffer
func NewLogScanner(r io.Reader, bufferSize int) *LogScanner {
	return &LogScanner{
		reader: bufio.NewReaderSize(r, bufferSize),
	}
}

func (s *LogScanner) Scan() bool {
	// initialize the skip current line flag
	tooLargeNeedToSkipCurrentLine := false

	for {
		var prefix bool
		s.bytes, prefix, s.err = s.reader.ReadLine()

		// if encouter any error, exit
		if s.err != nil {
			return false
		}

		// if unable to read into the buffer, set skip current line flag and continue to read
		if prefix {
			tooLargeNeedToSkipCurrentLine = true
			continue
		}

		// if read to the line end, check if need to skip
		if tooLargeNeedToSkipCurrentLine {
			// reset the flag and continue to read next line
			tooLargeNeedToSkipCurrentLine = false
			continue
		}

		return true
	}
}

func (s *LogScanner) Bytes() []byte {
	return s.bytes
}

func (s *LogScanner) Err() error {
	if s.err != nil && s.err != io.EOF {
		return s.err
	} else {
		return nil
	}
}
