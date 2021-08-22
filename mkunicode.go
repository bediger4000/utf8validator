package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type OutString struct {
	str   string
	valid bool
}

// Some of these cases use Go's "\xWXYZ" literals,
// which end up using the compiler to get valid UTF-8 encoding.
// The invalid cases are "\xXY" bytes to get invalid byte sequencese.
var variants = []OutString{
	OutString{str: "\u20ac", valid: true},                                     // euro
	OutString{str: "\u0400\u0401\u0402\u0408\u0416\u0429", valid: true},       // cyrillic
	OutString{str: "\u0041\u006d\u0065\u0301\u006c\u0069\u0065", valid: true}, // ASCII+Accent
	OutString{str: "\u0041\u006d\u00e9\u006c\u0069\u0065", valid: true},       // ASCII+Accent
	OutString{str: "\xff\x41", valid: false},                                  // invalid
	OutString{str: "\u20ac\xff\x41\u20ac", valid: false},                      // code point, invalid, code point
	OutString{str: "\u20ac\xAA\xAA\u20ac", valid: false},                      // euro, 2 invalid bytes, euro
	OutString{str: "\u2655\u265f", valid: true},                               // chess pieces
	OutString{str: "\xf0\x92\x80\x80\xf0\x93\x88\xa8", valid: true},           // Valid 4-byte encodings, Sumerian and Hieroglypics
	OutString{str: "\xf0\x42\x80\x80\xf0\x93\x88\xa8", valid: false},          // Invalid 4-byte encodings
	OutString{str: "\xf0\x92\x41\x80\xf0\x93\x88\xa8", valid: false},          // Valid 4-byte encodings, Sumerian and Hieroglypics
	OutString{str: "\xf0\x92\x80\x41\xf0\x93\x88\xa8", valid: false},          // Valid 4-byte encodings, Sumerian and Hieroglypics
	OutString{str: "aaa\xf0\x92\x80\x70\xf0\x93\x88\xa8bbb", valid: false},    // Invalid 4-byte encodings
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s [0-%d]\n", os.Args[0], len(variants)-1)
		fmt.Printf("Print one of %d variant stream of bytes, most valid UTF-8, to stdout\n", len(variants))
		return
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	validity := "valid"
	if !variants[n].valid {
		validity = "invalid"
	}
	fmt.Fprintf(os.Stderr, "%s\n", validity)
	fmt.Printf("%s", variants[n].str)
}
