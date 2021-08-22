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
const (
	Ascii = iota
	FirstOf2
	FirstOf3
	FirstOf4
	Subsequent
	Other
)

type BytePrefix int

// nextState = transitionTable[State][BytePrefix]
var transitionTable = [][]State{
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
	var value uint

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

		switch state {
		case Begin:
			value = v
		case Fail:
			value = v
		case Read1of2, Read1of3, Read1of4, Read2of3, Read2of4, Read3of4:
			value = (value << 5) | value
		}

		state = transitionTable[state][t]

		switch state {
		case Begin:
			if unicode.IsPrint(rune(value)) {
				fmt.Printf("Codepoint %6x %6d %c\n", value, value, value)
			} else {
				fmt.Printf("Codepoint %6x %6d %c\n", value, value, value)
			}
			value = 0
		case Fail:
			fmt.Printf("Invalid UTF-8 byte value %x at byte number %d\n", b, byteCount)
			value = 0
		}

	}
	fmt.Printf("Read %d bytes\n", byteCount)
}

func byteType(b byte) (BytePrefix, uint) {
	if b>>7 == 0 {
		return Ascii, uint(b)
	}
	if b>>5 == 0b110 {
		return FirstOf2, uint(b & 0b00011111)
	}
	if b>>4 == 0b1110 {
		return FirstOf3, uint(b & 0b00001111)
	}
	if b>>4 == 0b1111 {
		return FirstOf4, uint(b & 0b00001111)
	}
	if b>>6 == 0b10 {
		return Subsequent, uint(b & 0b00111111)
	}

	return Other, 0
}
