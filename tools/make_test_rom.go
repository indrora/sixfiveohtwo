package main

import (
	"os"
)

func main() {
	rom := make([]byte, 32768)
	
	prog := []byte{
		0xA9, 0x48,       // LDA #$48 ('H')
		0x8D, 0x01, 0xF0, // STA $F001
		0xA9, 0x65,       // LDA #$65 ('e')
		0x8D, 0x01, 0xF0, // STA $F001  
		0xA9, 0x6C,       // LDA #$6C ('l')
		0x8D, 0x01, 0xF0, // STA $F001
		0x8D, 0x01, 0xF0, // STA $F001 (second 'l')
		0xA9, 0x6F,       // LDA #$6F ('o')
		0x8D, 0x01, 0xF0, // STA $F001
		0xA9, 0x0A,       // LDA #$0A (newline)
		0x8D, 0x01, 0xF0, // STA $F001
		0x00,             // BRK
	}
	
	copy(rom[0:], prog)
	
	rom[0x7FFC] = 0x00
	rom[0x7FFD] = 0x80
	
	f, err := os.Create("test.rom")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	
	f.Write(rom)
}