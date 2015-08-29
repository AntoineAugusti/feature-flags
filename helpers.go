package main

import (
	"strconv"
)

func uint32ToBytes(u uint32) []byte {
	return []byte(strconv.FormatUint(uint64(u), 10))
}

func intInSlice(a uint32, list []uint32) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func mapToSlice(m map[string]interface{}) []string {
	slice := make([]string, len(m))
	i := 0
	for _, value := range m {
		str, _ := value.(string)
		slice[i] = str
		i++
	}
	return slice
}
