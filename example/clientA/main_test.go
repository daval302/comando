package main

import (
	"testing"
)

func TestEncodeUUID(t *testing.T) {
	uuid := EncodeUUID("0xAA")

	// Get last two bytes
	uuid &= 0x000000FF

	if uuid != 0xAA {
		t.Errorf("Expected 0xAA, got %x", uuid)
	}
}
