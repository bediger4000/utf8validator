package main

import (
	"fmt"
	"os"
)

/*
Bytes   |           Byte format
-----------------------------------------------
   1     | 0xxxxxxx
   2     | 110xxxxx 10xxxxxx
   3     | 1110xxxx 10xxxxxx 10xxxxxx
   4     | 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
*/

func main() {
	var buf [1]byte

	byteCount := 0

	for {
		n, err := os.Stdin.Read(buf[:])
		if n == 0 {
			// end-of-file
			break
		}
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}
		byteCount++
	}
	fmt.Printf("Read %d bytes\n", byteCount)
}
