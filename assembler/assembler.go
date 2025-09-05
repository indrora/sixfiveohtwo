package assembler

import (
	"fmt"
	"io/ioutil"
)

type Assembler struct {
	symbols *SymbolTable
	output  []byte
	pc      uint16
	verbose bool
}

type AssemblerError struct {
	Line    int
	Message string
}

func (e *AssemblerError) Error() string {
	return fmt.Sprintf("line %d: %s", e.Line, e.Message)
}

func NewAssembler() *Assembler {
	return &Assembler{
		symbols: NewSymbolTable(),
		output:  make([]byte, 65536),
		pc:      0,
		verbose: false,
	}
}

func (a *Assembler) SetVerbose(verbose bool) {
	a.verbose = verbose
}

func (a *Assembler) AssembleFile(filename string) error {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	
	return a.Assemble(string(source))
}

func (a *Assembler) Assemble(source string) error {
	lexer := NewLexer(source)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return err
	}
	
	parser := NewParser(tokens, a.symbols)
	instructions, err := parser.Parse()
	if err != nil {
		return err
	}
	
	codegen := NewCodeGenerator(a.symbols)
	codegen.SetVerbose(a.verbose)
	return codegen.Generate(instructions, a.output)
}

func (a *Assembler) WriteROM(filename string, startAddr, size uint16) error {
	rom := make([]byte, size)
	
	if int(startAddr) < len(a.output) {
		endAddr := int(startAddr) + int(size)
		if endAddr > len(a.output) {
			endAddr = len(a.output)
		}
		
		if endAddr > int(startAddr) {
			copy(rom, a.output[startAddr:endAddr])
		}
	}
	
	return ioutil.WriteFile(filename, rom, 0644)
}