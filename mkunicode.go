package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	variants := []string{
		"\u20ac",                               // euro
		"\u0400\u0401\u0402\u0408\u0416\u0429", // cyrillic
		"\u0041\u006d\u0065\u0301\u006c\u0069\u0065", // ASCII+Accent
		"\u0041\u006d\u00e9\u006c\u0069\u0065",       // ASCII+Accent
		"\xff\x41",                                   // invalid
		"\u20ac\xff\x41\u20ac",                       // code point, invalid, code point
		"\u20ac\xAA\xAA\u20ac",
	}
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s [0-%d]\n", os.Args[0], len(variants)-1)
		fmt.Printf("Print one of %d variant stream of bytes, most valid UTF-8, to stdout\n", len(variants)-1)
		return
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", variants[n])
}
