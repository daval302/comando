package main

import (
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

	id := DecodeUUID(0xabcde608)

	if id != "8" {
		t.Errorf("Expected 0x08, got %x\n", id)
	}
}
