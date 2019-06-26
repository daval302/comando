package main

import (
	"log"
	"strconv"
	"testing"
)

func TestEncodeUUID(t *testing.T) {

	uuid := EncodeUUID("0xAA")

	// Get last two bytes
	uuid &= 0x000000FF

	if uuid != 0xAA {
		t.Errorf("Expected 0xAA, got %x\n", uuid)
	}
}

func TestDecodeUUID(t *testing.T) {

	id := DecodeUUID(0xabcde608) // last 2 byte : 08

	val, err := strconv.ParseInt(id, 16, 64)
	if err != nil {
		log.Fatal("Bad conversion: ", err)
	}

	// zeros 6 bytes left
	val &= 0x000000FF

	if val != 0x08 {
		t.Errorf("Expected 0x08, got %x\n", id)
	}
}
