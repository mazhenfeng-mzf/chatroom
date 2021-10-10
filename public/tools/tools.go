package tools

import (
	"os"

	"github.com/lxn/walk"
)

func IsValueInSlice(testint int, testslice []int) bool {
	for _, v := range testslice {
		if v == testint {
			return true
		}
	}
	return false
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func SetBlue() walk.Color {
	return walk.RGB(0x00, 0x99, 0xff)
}

func SetGrey() walk.Color {
	return walk.RGB(0xcc, 0xcc, 0xcc)
}

func SetGreen() walk.Color {
	return walk.RGB(0x00, 0xcc, 0x66)
}
