package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/indrora/sixfiveohtwo/assembler"
)

func main() {
	var outputFile string
	var startAddr uint
	var romSize uint
	var verbose bool
	
	flag.StringVar(&outputFile, "o", "", "output ROM file")
	flag.UintVar(&startAddr, "start", 0x8000, "ROM start address")
	flag.UintVar(&romSize, "size", 32768, "ROM size in bytes")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()
	
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] input.asm\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	
	inputFile := flag.Arg(0)
	
	if outputFile == "" {
		base := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
		outputFile = base + ".rom"
	}
	
	if verbose {
		fmt.Printf("Assembling %s -> %s\n", inputFile, outputFile)
		fmt.Printf("ROM start: 0x%04X, size: %d bytes\n", startAddr, romSize)
	}
	
	asm := assembler.NewAssembler()
	asm.SetVerbose(verbose)
	
	if err := asm.AssembleFile(inputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Assembly error: %v\n", err)
		os.Exit(1)
	}
	
	if err := asm.WriteROM(outputFile, uint16(startAddr), uint16(romSize)); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing ROM: %v\n", err)
		os.Exit(1)
	}
	
	if verbose {
		fmt.Printf("Successfully assembled %s\n", outputFile)
	}
}
