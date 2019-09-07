package main

import (
	"encoding/binary"
)

func combineIDWithMessage(id int64, message []byte) []byte {

	// make a buffer of 64 bytes
	buffer := make([]byte, 64)

	// add an int64 id into a []byte buffer
	binary.PutVarint(buffer, id)

	// concatenate message to buffer
	n := copy(buffer[8:], message)

	// return the combined message and truncate
	return buffer[:n+8]

}
