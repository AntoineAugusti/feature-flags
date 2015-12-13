package helpers

import (
	"strconv"
)

// Transform a uint32 to a byte slice
func Uint32ToBytes(u uint32) []byte {
	return []byte(strconv.FormatUint(uint64(u), 10))
}

// Check if an int is in a slice
func IntInSlice(a uint32, list []uint32) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Check if a string is in a slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
