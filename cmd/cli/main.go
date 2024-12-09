package main

import (
	"fmt"
	"log"

	"github.com/vitaminmoo/poe-get-version/internal/version"
)

func main() {
	poe2, err := version.Poe2()
	if err != nil {
		log.Printf("getting poe2 version: %v", err)
		poe2 = "error"
	}
	fmt.Printf("poe2: %s\n", poe2)
}

func hexdump(data []byte) {
	for i := 0; i < len(data); i += 16 {
		fmt.Printf("%08x  ", i)

		// Hex dump
		for j := 0; j < 16; j++ {
			if i+j < len(data) {
				fmt.Printf("%02x ", data[i+j])
			} else {
				fmt.Print("   ")
			}
			if j == 7 {
				fmt.Print(" ")
			}
		}

		// ASCII dump
		fmt.Print(" |")
		for j := 0; j < 16; j++ {
			if i+j < len(data) {
				if data[i+j] >= 32 && data[i+j] <= 126 {
					fmt.Printf("%c", data[i+j])
				} else {
					fmt.Print(".")
				}
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("|")
	}
}
