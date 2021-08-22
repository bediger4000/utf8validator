package main

/* UTF-8 validator.
 * Doesn't re-synchronize on encoding errors,
 * doesn't calculate code points.
 */

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

// Byte prefixes
const (
	Ascii = iota
	FirstOf2
	FirstOf3
	FirstOf4
	Subsequent
	Other
)

func readByte() int {
	var b [1]byte
	n, err := os.Stdin.Read(b[:])
	if n == 0 {
		// end-of-file
		return 0
	}
	if err != nil {
		// mechanical error, not erroneous UTF-8 encoding
		fmt.Printf("Error: %v\n", err)
		return -1
	}
	return int(b[0])
}

func main() {
	byteCount := 0
	validUTF8 := true

	for i := readByte(); validUTF8 && i > 0; i = readByte() {
		byteCount++

		t := byteType(i)

		if t == Other || t == Subsequent {
			validUTF8 = false
			break
		}

		var readN int
		switch t {
		case Ascii:
			readN = 0
		case FirstOf2:
			readN = 1
		case FirstOf3:
			readN = 2
		case FirstOf4:
			readN = 3
		}

		for j := 0; j < readN; j++ {
			k := readByte()
			if k <= 0 {
				break
			}

			t := byteType(k)

			if t == Other || t != Subsequent {
				validUTF8 = false
				break
			}
		}
	}

	phrase := "Valid"
	if !validUTF8 {
		phrase = "Invalid"
	}
	fmt.Printf("%s UTF-8, Read %d bytes\n", phrase, byteCount)
}

func byteType(b int) int {
	if (b>>7)&0b1 == 0 {
		return Ascii
	}
	if (b>>5)&0b111 == 0b110 {
		return FirstOf2
	}
	if (b>>4)&0b1111 == 0b1110 {
		return FirstOf3
	}
	if (b>>4)&0b1111 == 0b1111 {
		return FirstOf4
	}
	if (b>>6)&0b11 == 0b10 {
		return Subsequent
	}

	return Other
}
