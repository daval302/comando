package main

import (
	"encoding/binary"
	"testing"
)

func TestCombineIDWithMessage(t *testing.T) {

	// use time.Now().Unix() to generate the id
	var id int64 = 0x01
	message := "Hello Mr Robertson"
	buffer := combineIDWithMessage(id, []byte(message))

	// combinedId should be equals to id
	// combined message should be equal to message
	combinedID, _ := binary.Varint(buffer[:8])
	combinedMessage := string(buffer[8:])

	if combinedID != id {
		t.Errorf("Expected %d, got %d\n", id, combinedID)
	}

	if combinedMessage != message {
		t.Errorf("Expected %s, got %s\n", message, combinedMessage)
	}

}
