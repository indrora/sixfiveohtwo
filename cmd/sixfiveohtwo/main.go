package main

import (
	"fmt"
	"os"
	"github.com/indrora/sixfiveohtwo/emulator"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: sixfiveohtwo <rom_file>")
		os.Exit(1)
	}

	romFile := os.Args[1]
	
	cpu := emulator.NewCPU()
	
	if err := cpu.LoadROM(romFile); err != nil {
		fmt.Printf("Error loading ROM: %v\n", err)
		os.Exit(1)
	}
	
	cpu.Reset()
	
	fmt.Println("6502 Emulator started")
	cpu.Run()
}