package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: romdump <rom_file>")
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ROM size: %d bytes\n", len(data))
	fmt.Printf("First 32 bytes:\n")
	for i := 0; i < 32 && i < len(data); i += 16 {
		end := i + 16
		if end > len(data) {
			end = len(data)
		}
		
		fmt.Printf("%04X: ", i)
		for j := i; j < end; j++ {
			fmt.Printf("%02X ", data[j])
		}
		fmt.Println()
	}
	
	fmt.Printf("Reset vector area (last 16 bytes):\n")
	start := len(data) - 16
	if start < 0 {
		start = 0
	}
	
	for i := start; i < len(data); i += 16 {
		end := i + 16
		if end > len(data) {
			end = len(data)
		}
		
		fmt.Printf("%04X: ", 0x8000+i)
		for j := i; j < end; j++ {
			fmt.Printf("%02X ", data[j])
		}
		fmt.Println()
	}
	
	resetVectorOffset := 0x7FFC
	if resetVectorOffset < len(data) {
		fmt.Printf("Reset vector at ROM offset 0x%04X (address 0xFFFC): %02X %02X\n", 
			resetVectorOffset, data[resetVectorOffset], data[resetVectorOffset+1])
	}
}