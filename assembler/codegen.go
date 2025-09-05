package assembler

import (
	"fmt"
)

type CodeGenerator struct {
	symbols *SymbolTable
	pc      uint16
	verbose bool
}

type OpcodeInfo struct {
	Opcode uint8
	Size   int
}

var opcodeTable = map[string]map[AddressingMode]OpcodeInfo{
	"ADC": {
		AddrImmediate:       {0x69, 2},
		AddrZeroPage:        {0x65, 2},
		AddrZeroPageX:       {0x75, 2},
		AddrAbsolute:        {0x6D, 3},
		AddrAbsoluteX:       {0x7D, 3},
		AddrAbsoluteY:       {0x79, 3},
		AddrIndexedIndirect: {0x61, 2},
		AddrIndirectIndexed: {0x71, 2},
	},
	"AND": {
		AddrImmediate:       {0x29, 2},
		AddrZeroPage:        {0x25, 2},
		AddrZeroPageX:       {0x35, 2},
		AddrAbsolute:        {0x2D, 3},
		AddrAbsoluteX:       {0x3D, 3},
		AddrAbsoluteY:       {0x39, 3},
		AddrIndexedIndirect: {0x21, 2},
		AddrIndirectIndexed: {0x31, 2},
	},
	"ASL": {
		AddrAccumulator: {0x0A, 1},
		AddrZeroPage:    {0x06, 2},
		AddrZeroPageX:   {0x16, 2},
		AddrAbsolute:    {0x0E, 3},
		AddrAbsoluteX:   {0x1E, 3},
	},
	"BCC": {AddrRelative: {0x90, 2}},
	"BCS": {AddrRelative: {0xB0, 2}},
	"BEQ": {AddrRelative: {0xF0, 2}},
	"BIT": {
		AddrZeroPage: {0x24, 2},
		AddrAbsolute: {0x2C, 3},
	},
	"BMI": {AddrRelative: {0x30, 2}},
	"BNE": {AddrRelative: {0xD0, 2}},
	"BPL": {AddrRelative: {0x10, 2}},
	"BRK": {AddrImplicit: {0x00, 1}},
	"BVC": {AddrRelative: {0x50, 2}},
	"BVS": {AddrRelative: {0x70, 2}},
	"CLC": {AddrImplicit: {0x18, 1}},
	"CLD": {AddrImplicit: {0xD8, 1}},
	"CLI": {AddrImplicit: {0x58, 1}},
	"CLV": {AddrImplicit: {0xB8, 1}},
	"CMP": {
		AddrImmediate:       {0xC9, 2},
		AddrZeroPage:        {0xC5, 2},
		AddrZeroPageX:       {0xD5, 2},
		AddrAbsolute:        {0xCD, 3},
		AddrAbsoluteX:       {0xDD, 3},
		AddrAbsoluteY:       {0xD9, 3},
		AddrIndexedIndirect: {0xC1, 2},
		AddrIndirectIndexed: {0xD1, 2},
	},
	"CPX": {
		AddrImmediate: {0xE0, 2},
		AddrZeroPage:  {0xE4, 2},
		AddrAbsolute:  {0xEC, 3},
	},
	"CPY": {
		AddrImmediate: {0xC0, 2},
		AddrZeroPage:  {0xC4, 2},
		AddrAbsolute:  {0xCC, 3},
	},
	"DEC": {
		AddrZeroPage:  {0xC6, 2},
		AddrZeroPageX: {0xD6, 2},
		AddrAbsolute:  {0xCE, 3},
		AddrAbsoluteX: {0xDE, 3},
	},
	"DEX": {AddrImplicit: {0xCA, 1}},
	"DEY": {AddrImplicit: {0x88, 1}},
	"EOR": {
		AddrImmediate:       {0x49, 2},
		AddrZeroPage:        {0x45, 2},
		AddrZeroPageX:       {0x55, 2},
		AddrAbsolute:        {0x4D, 3},
		AddrAbsoluteX:       {0x5D, 3},
		AddrAbsoluteY:       {0x59, 3},
		AddrIndexedIndirect: {0x41, 2},
		AddrIndirectIndexed: {0x51, 2},
	},
	"INC": {
		AddrZeroPage:  {0xE6, 2},
		AddrZeroPageX: {0xF6, 2},
		AddrAbsolute:  {0xEE, 3},
		AddrAbsoluteX: {0xFE, 3},
	},
	"INX": {AddrImplicit: {0xE8, 1}},
	"INY": {AddrImplicit: {0xC8, 1}},
	"JMP": {
		AddrAbsolute: {0x4C, 3},
		AddrIndirect: {0x6C, 3},
	},
	"JSR": {AddrAbsolute: {0x20, 3}},
	"LDA": {
		AddrImmediate:       {0xA9, 2},
		AddrZeroPage:        {0xA5, 2},
		AddrZeroPageX:       {0xB5, 2},
		AddrAbsolute:        {0xAD, 3},
		AddrAbsoluteX:       {0xBD, 3},
		AddrAbsoluteY:       {0xB9, 3},
		AddrIndexedIndirect: {0xA1, 2},
		AddrIndirectIndexed: {0xB1, 2},
	},
	"LDX": {
		AddrImmediate: {0xA2, 2},
		AddrZeroPage:  {0xA6, 2},
		AddrZeroPageY: {0xB6, 2},
		AddrAbsolute:  {0xAE, 3},
		AddrAbsoluteY: {0xBE, 3},
	},
	"LDY": {
		AddrImmediate: {0xA0, 2},
		AddrZeroPage:  {0xA4, 2},
		AddrZeroPageX: {0xB4, 2},
		AddrAbsolute:  {0xAC, 3},
		AddrAbsoluteX: {0xBC, 3},
	},
	"LSR": {
		AddrAccumulator: {0x4A, 1},
		AddrZeroPage:    {0x46, 2},
		AddrZeroPageX:   {0x56, 2},
		AddrAbsolute:    {0x4E, 3},
		AddrAbsoluteX:   {0x5E, 3},
	},
	"NOP": {AddrImplicit: {0xEA, 1}},
	"ORA": {
		AddrImmediate:       {0x09, 2},
		AddrZeroPage:        {0x05, 2},
		AddrZeroPageX:       {0x15, 2},
		AddrAbsolute:        {0x0D, 3},
		AddrAbsoluteX:       {0x1D, 3},
		AddrAbsoluteY:       {0x19, 3},
		AddrIndexedIndirect: {0x01, 2},
		AddrIndirectIndexed: {0x11, 2},
	},
	"PHA": {AddrImplicit: {0x48, 1}},
	"PHP": {AddrImplicit: {0x08, 1}},
	"PLA": {AddrImplicit: {0x68, 1}},
	"PLP": {AddrImplicit: {0x28, 1}},
	"ROL": {
		AddrAccumulator: {0x2A, 1},
		AddrZeroPage:    {0x26, 2},
		AddrZeroPageX:   {0x36, 2},
		AddrAbsolute:    {0x2E, 3},
		AddrAbsoluteX:   {0x3E, 3},
	},
	"ROR": {
		AddrAccumulator: {0x6A, 1},
		AddrZeroPage:    {0x66, 2},
		AddrZeroPageX:   {0x76, 2},
		AddrAbsolute:    {0x6E, 3},
		AddrAbsoluteX:   {0x7E, 3},
	},
	"RTI": {AddrImplicit: {0x40, 1}},
	"RTS": {AddrImplicit: {0x60, 1}},
	"SBC": {
		AddrImmediate:       {0xE9, 2},
		AddrZeroPage:        {0xE5, 2},
		AddrZeroPageX:       {0xF5, 2},
		AddrAbsolute:        {0xED, 3},
		AddrAbsoluteX:       {0xFD, 3},
		AddrAbsoluteY:       {0xF9, 3},
		AddrIndexedIndirect: {0xE1, 2},
		AddrIndirectIndexed: {0xF1, 2},
	},
	"SEC": {AddrImplicit: {0x38, 1}},
	"SED": {AddrImplicit: {0xF8, 1}},
	"SEI": {AddrImplicit: {0x78, 1}},
	"STA": {
		AddrZeroPage:        {0x85, 2},
		AddrZeroPageX:       {0x95, 2},
		AddrAbsolute:        {0x8D, 3},
		AddrAbsoluteX:       {0x9D, 3},
		AddrAbsoluteY:       {0x99, 3},
		AddrIndexedIndirect: {0x81, 2},
		AddrIndirectIndexed: {0x91, 2},
	},
	"STX": {
		AddrZeroPage:  {0x86, 2},
		AddrZeroPageY: {0x96, 2},
		AddrAbsolute:  {0x8E, 3},
	},
	"STY": {
		AddrZeroPage:  {0x84, 2},
		AddrZeroPageX: {0x94, 2},
		AddrAbsolute:  {0x8C, 3},
	},
	"TAX": {AddrImplicit: {0xAA, 1}},
	"TAY": {AddrImplicit: {0xA8, 1}},
	"TSX": {AddrImplicit: {0xBA, 1}},
	"TXA": {AddrImplicit: {0x8A, 1}},
	"TXS": {AddrImplicit: {0x9A, 1}},
	"TYA": {AddrImplicit: {0x98, 1}},
}

func NewCodeGenerator(symbols *SymbolTable) *CodeGenerator {
	return &CodeGenerator{
		symbols: symbols,
		pc:      0,
		verbose: false,
	}
}

func (cg *CodeGenerator) SetVerbose(verbose bool) {
	cg.verbose = verbose
}

func (cg *CodeGenerator) Generate(instructions []Instruction, output []byte) error {
	if cg.verbose {
		fmt.Printf("Generating code for %d instructions\n", len(instructions))
	}
	
	if err := cg.firstPass(instructions); err != nil {
		return err
	}
	
	if err := cg.secondPass(instructions, output); err != nil {
		return err
	}
	
	if !cg.symbols.IsResolved() {
		undefined := cg.symbols.GetUndefined()
		return fmt.Errorf("undefined symbols: %v", undefined)
	}
	
	if cg.verbose {
		fmt.Printf("Code generation complete\n")
	}
	return nil
}

func (cg *CodeGenerator) firstPass(instructions []Instruction) error {
	cg.pc = 0
	
	for i := range instructions {
		inst := &instructions[i]
		
		switch inst.Type {
		case InstDirective:
			switch inst.DirectiveName {
			case "org":
				cg.pc = inst.Operand
			case "word":
				inst.Address = cg.pc
				cg.pc += 2
			case "byte":
				inst.Address = cg.pc
				cg.pc += uint16(len(inst.DirectiveData))
			}
			
		case InstLabel:
			cg.symbols.Define(inst.Mnemonic, cg.pc)
			
		case InstMnemonic:
			inst.Address = cg.pc
			
			opcodeMap, exists := opcodeTable[inst.Mnemonic]
			if !exists {
				return fmt.Errorf("unknown mnemonic '%s' at line %d", inst.Mnemonic, inst.Line)
			}
			
			opcodeInfo, exists := opcodeMap[inst.AddressMode]
			if !exists {
				return fmt.Errorf("invalid addressing mode for '%s' at line %d", inst.Mnemonic, inst.Line)
			}
			
			cg.pc += uint16(opcodeInfo.Size)
		}
	}
	
	return nil
}

func (cg *CodeGenerator) secondPass(instructions []Instruction, output []byte) error {
	for _, inst := range instructions {
		switch inst.Type {
		case InstDirective:
			switch inst.DirectiveName {
			case "word":
				if inst.OperandLabel != "" {
					operand, err := cg.symbols.Resolve(inst.OperandLabel)
					if err != nil {
						return fmt.Errorf("line %d: %v", inst.Line, err)
					}
					if cg.verbose {
						fmt.Printf("Writing word label '%s' = 0x%04X at address 0x%04X\n", inst.OperandLabel, operand, inst.Address)
					}
					output[inst.Address] = uint8(operand & 0xFF)
					output[inst.Address+1] = uint8((operand >> 8) & 0xFF)
				} else if len(inst.DirectiveData) >= 2 {
					if cg.verbose {
						fmt.Printf("Writing word data %02X %02X at address 0x%04X\n", inst.DirectiveData[0], inst.DirectiveData[1], inst.Address)
					}
					output[inst.Address] = inst.DirectiveData[0]
					output[inst.Address+1] = inst.DirectiveData[1]
				}
				
			case "byte":
				for i, b := range inst.DirectiveData {
					output[inst.Address+uint16(i)] = b
				}
			}
			
		case InstMnemonic:
			opcodeMap := opcodeTable[inst.Mnemonic]
			opcodeInfo := opcodeMap[inst.AddressMode]
			
			if cg.verbose {
				fmt.Printf("Writing opcode %s (0x%02X) at address 0x%04X\n", inst.Mnemonic, opcodeInfo.Opcode, inst.Address)
			}
			output[inst.Address] = opcodeInfo.Opcode
			
			if opcodeInfo.Size > 1 {
				operand := inst.Operand
				
				if inst.OperandLabel != "" {
					var err error
					operand, err = cg.symbols.Resolve(inst.OperandLabel)
					if err != nil {
						return fmt.Errorf("line %d: %v", inst.Line, err)
					}
				}
				
				if inst.AddressMode == AddrRelative {
					offset := int16(operand) - int16(inst.Address+2)
					if offset < -128 || offset > 127 {
						return fmt.Errorf("branch out of range at line %d", inst.Line)
					}
					output[inst.Address+1] = uint8(offset)
				} else if opcodeInfo.Size == 2 {
					output[inst.Address+1] = uint8(operand & 0xFF)
				} else if opcodeInfo.Size == 3 {
					output[inst.Address+1] = uint8(operand & 0xFF)
					output[inst.Address+2] = uint8((operand >> 8) & 0xFF)
				}
			}
		}
	}
	
	return nil
}