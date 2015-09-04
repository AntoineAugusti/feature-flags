package helpers

import (
	"strconv"
)

func Uint32ToBytes(u uint32) []byte {
	return []byte(strconv.FormatUint(uint64(u), 10))
}

func IntInSlice(a uint32, list []uint32) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
