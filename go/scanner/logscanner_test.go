package scanner

import (
	"fmt"
	"strings"
	"testing"
)

func TestLogScanner(t *testing.T) {
	testInput := "123457890\n123457890ABCDEF\n123457890ABCDE\n123457890ABCDEFG\n123457890ABCDEF"
	bufferSize := 16

	fmt.Println("buffer size: ", bufferSize)
	fmt.Println("test input:")
	fmt.Println(testInput)

	logScanner := NewLogScanner(strings.NewReader(testInput), bufferSize)

	fmt.Println("test output:")
	for logScanner.Scan() {
		fmt.Println(string(logScanner.Bytes()))
	}

	if err := logScanner.Err(); err != nil {
		fmt.Println("scan err: ", err)
	} else {
		fmt.Println("scan err: none")
	}
}
