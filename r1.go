package main

import (
	"fmt"
	"os"
	"unicode"
)

/*
Bytes   |           Byte format
-----------------------------------------------
   1     | 0xxxxxxx
   2     | 110xxxxx 10xxxxxx
   3     | 1110xxxx 10xxxxxx 10xxxxxx
   4     | 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
*/

// States
type State int

const (
	Fail = iota
	Begin
	Read1of2
	Read1of3
	Read2of3
	Read1of4
	Read2of4
	Read3of4
)

// Byte prefixes
type BytePrefix int

const (
	Ascii = iota
	FirstOf2
	FirstOf3
	FirstOf4
	Subsequent
	Other
)

// nextState = transitionTable[State][BytePrefix]
var transitionTable = [][]State{
	//                    Ascii FirstOf2 FirstOf3 FirstOf4 Subsequent Other
	/* Fail     */ []State{Begin, Read1of2, Read1of3, Read1of4, Fail, Fail},
	/* Begin    */ []State{Begin, Read1of2, Read1of3, Read1of4, Fail, Fail},
	/* Read1of2 */ []State{Fail, Fail, Fail, Fail, Begin, Fail},
	/* Read1of3 */ []State{Fail, Fail, Fail, Fail, Read2of3, Fail},
	/* Read2of3 */ []State{Fail, Fail, Fail, Fail, Begin, Fail},
	/* Read1of4 */ []State{Fail, Fail, Fail, Fail, Read2of4, Fail},
	/* Read2of4 */ []State{Fail, Fail, Fail, Fail, Read3of4, Fail},
	/* Read3of4 */ []State{Fail, Fail, Fail, Fail, Begin, Fail},
}

func main() {
	var buf [1]byte

	byteCount := 0

	var state State = Begin
	var codePoint uint
	validUTF8 := true // is stream of bytes valid unicode or not?

	// This loop conflates reading stdin a byte at a time,
	// and working the state machine.
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

		b := buf[0]

		t, v := byteType(b)

		if !(state == Begin || state == Fail) {
			codePoint <<= 6
		}
		codePoint |= v

		state = transitionTable[state][t]

		switch state {
		case Begin:
			if unicode.IsPrint(rune(codePoint)) {
				fmt.Printf("Codepoint %6x %6d %c\n", codePoint, codePoint, codePoint)
			} else {
				fmt.Printf("Codepoint %6x %6d\n", codePoint, codePoint)
			}
			codePoint = 0
		case Fail:
			// traditionally code point 0xfffd
			fmt.Printf("Invalid UTF-8 byte value %x at byte number %d\n", b, byteCount)
			codePoint = 0
			validUTF8 = false
		}

	}

	phrase := "Valid"
	if !validUTF8 {
		phrase = "Invalid"
	}
	fmt.Printf("%s UTF-8, Read %d bytes\n", phrase, byteCount)
}

// byteType returns one of the bit-prefix types,
// and the value that's the rest of the byte.
func byteType(b byte) (BytePrefix, uint) {
	if (b>>7)&0b1 == 0 {
		return Ascii, uint(b)
	}
	if (b>>6)&0b11 == 0b10 {
		return Subsequent, uint(b & 0b00111111)
	}
	if (b>>5)&0b111 == 0b110 {
		return FirstOf2, uint(b & 0b00011111)
	}
	if (b>>4)&0b1111 == 0b1110 {
		return FirstOf3, uint(b & 0b00001111)
	}
	if (b>>3)&0b11111 == 0b11110 {
		return FirstOf4, uint(b & 0b00000111)
	}

	return Other, 0
}
