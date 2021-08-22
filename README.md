# Daily Coding Problem: Problem #965 [Easy] 

This problem was asked by Google.

UTF-8 is a character encoding that maps each symbol to one, two, three, or four bytes.

For example, the Euro sign, €,
corresponds to the three bytes 11100010 10000010 10101100.
The rules for mapping characters are as follows:

* For a single-byte character, the first bit must be zero.
* For an n-byte character, the first byte starts with n ones and a zero.

The other n - 1 bytes all start with 10.

Visually, this can be represented as follows.

```
 Bytes   |           Byte format
-----------------------------------------------
   1     | 0xxxxxxx
   2     | 110xxxxx 10xxxxxx
   3     | 1110xxxx 10xxxxxx 10xxxxxx
   4     | 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
```

Write a program that takes in an array of integers representing byte values,
and returns whether it is a valid UTF-8 encoding.

## Analysis

I chose to use a state machine:

![UTF-8 validator state machine](states.png)

## Interview Analysis

The problem statement doesn't cover the whole of UTF-8:
a program can re-synchronize with a UTF-8 bytestream after encountering
a goofed up code point.
Some 0xxxxxxx bytes are also invalid: not all of ASCII is UTF-8.
