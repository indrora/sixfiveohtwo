package emulator

type AddressingMode int

const (
	Implicit AddressingMode = iota
	Accumulator
	Immediate
	ZeroPage
	ZeroPageX
	ZeroPageY
	Absolute
	AbsoluteX
	AbsoluteY
	Indirect
	IndexedIndirect
	IndirectIndexed
	Relative
)

type Instruction struct {
	Name        string
	AddressMode AddressingMode
	Cycles      int
	Execute     func(*CPU, uint16)
}

var instructions = [256]Instruction{
	0x00: {"BRK", Implicit, 7, (*CPU).BRK},
	0x01: {"ORA", IndexedIndirect, 6, (*CPU).ORA},
	0x05: {"ORA", ZeroPage, 3, (*CPU).ORA},
	0x06: {"ASL", ZeroPage, 5, (*CPU).ASL},
	0x08: {"PHP", Implicit, 3, (*CPU).PHP},
	0x09: {"ORA", Immediate, 2, (*CPU).ORA},
	0x0A: {"ASL", Accumulator, 2, (*CPU).ASLA},
	0x0D: {"ORA", Absolute, 4, (*CPU).ORA},
	0x0E: {"ASL", Absolute, 6, (*CPU).ASL},
	0x10: {"BPL", Relative, 2, (*CPU).BPL},
	0x11: {"ORA", IndirectIndexed, 5, (*CPU).ORA},
	0x15: {"ORA", ZeroPageX, 4, (*CPU).ORA},
	0x16: {"ASL", ZeroPageX, 6, (*CPU).ASL},
	0x18: {"CLC", Implicit, 2, (*CPU).CLC},
	0x19: {"ORA", AbsoluteY, 4, (*CPU).ORA},
	0x1D: {"ORA", AbsoluteX, 4, (*CPU).ORA},
	0x1E: {"ASL", AbsoluteX, 7, (*CPU).ASL},
	0x20: {"JSR", Absolute, 6, (*CPU).JSR},
	0x21: {"AND", IndexedIndirect, 6, (*CPU).AND},
	0x24: {"BIT", ZeroPage, 3, (*CPU).BIT},
	0x25: {"AND", ZeroPage, 3, (*CPU).AND},
	0x26: {"ROL", ZeroPage, 5, (*CPU).ROL},
	0x28: {"PLP", Implicit, 4, (*CPU).PLP},
	0x29: {"AND", Immediate, 2, (*CPU).AND},
	0x2A: {"ROL", Accumulator, 2, (*CPU).ROLA},
	0x2C: {"BIT", Absolute, 4, (*CPU).BIT},
	0x2D: {"AND", Absolute, 4, (*CPU).AND},
	0x2E: {"ROL", Absolute, 6, (*CPU).ROL},
	0x30: {"BMI", Relative, 2, (*CPU).BMI},
	0x31: {"AND", IndirectIndexed, 5, (*CPU).AND},
	0x35: {"AND", ZeroPageX, 4, (*CPU).AND},
	0x36: {"ROL", ZeroPageX, 6, (*CPU).ROL},
	0x38: {"SEC", Implicit, 2, (*CPU).SEC},
	0x39: {"AND", AbsoluteY, 4, (*CPU).AND},
	0x3D: {"AND", AbsoluteX, 4, (*CPU).AND},
	0x3E: {"ROL", AbsoluteX, 7, (*CPU).ROL},
	0x40: {"RTI", Implicit, 6, (*CPU).RTI},
	0x41: {"EOR", IndexedIndirect, 6, (*CPU).EOR},
	0x45: {"EOR", ZeroPage, 3, (*CPU).EOR},
	0x46: {"LSR", ZeroPage, 5, (*CPU).LSR},
	0x48: {"PHA", Implicit, 3, (*CPU).PHA},
	0x49: {"EOR", Immediate, 2, (*CPU).EOR},
	0x4A: {"LSR", Accumulator, 2, (*CPU).LSRA},
	0x4C: {"JMP", Absolute, 3, (*CPU).JMP},
	0x4D: {"EOR", Absolute, 4, (*CPU).EOR},
	0x4E: {"LSR", Absolute, 6, (*CPU).LSR},
	0x50: {"BVC", Relative, 2, (*CPU).BVC},
	0x51: {"EOR", IndirectIndexed, 5, (*CPU).EOR},
	0x55: {"EOR", ZeroPageX, 4, (*CPU).EOR},
	0x56: {"LSR", ZeroPageX, 6, (*CPU).LSR},
	0x58: {"CLI", Implicit, 2, (*CPU).CLI},
	0x59: {"EOR", AbsoluteY, 4, (*CPU).EOR},
	0x5D: {"EOR", AbsoluteX, 4, (*CPU).EOR},
	0x5E: {"LSR", AbsoluteX, 7, (*CPU).LSR},
	0x60: {"RTS", Implicit, 6, (*CPU).RTS},
	0x61: {"ADC", IndexedIndirect, 6, (*CPU).ADC},
	0x65: {"ADC", ZeroPage, 3, (*CPU).ADC},
	0x66: {"ROR", ZeroPage, 5, (*CPU).ROR},
	0x68: {"PLA", Implicit, 4, (*CPU).PLA},
	0x69: {"ADC", Immediate, 2, (*CPU).ADC},
	0x6A: {"ROR", Accumulator, 2, (*CPU).RORA},
	0x6C: {"JMP", Indirect, 5, (*CPU).JMPI},
	0x6D: {"ADC", Absolute, 4, (*CPU).ADC},
	0x6E: {"ROR", Absolute, 6, (*CPU).ROR},
	0x70: {"BVS", Relative, 2, (*CPU).BVS},
	0x71: {"ADC", IndirectIndexed, 5, (*CPU).ADC},
	0x75: {"ADC", ZeroPageX, 4, (*CPU).ADC},
	0x76: {"ROR", ZeroPageX, 6, (*CPU).ROR},
	0x78: {"SEI", Implicit, 2, (*CPU).SEI},
	0x79: {"ADC", AbsoluteY, 4, (*CPU).ADC},
	0x7D: {"ADC", AbsoluteX, 4, (*CPU).ADC},
	0x7E: {"ROR", AbsoluteX, 7, (*CPU).ROR},
	0x81: {"STA", IndexedIndirect, 6, (*CPU).STA},
	0x84: {"STY", ZeroPage, 3, (*CPU).STY},
	0x85: {"STA", ZeroPage, 3, (*CPU).STA},
	0x86: {"STX", ZeroPage, 3, (*CPU).STX},
	0x88: {"DEY", Implicit, 2, (*CPU).DEY},
	0x8A: {"TXA", Implicit, 2, (*CPU).TXA},
	0x8C: {"STY", Absolute, 4, (*CPU).STY},
	0x8D: {"STA", Absolute, 4, (*CPU).STA},
	0x8E: {"STX", Absolute, 4, (*CPU).STX},
	0x90: {"BCC", Relative, 2, (*CPU).BCC},
	0x91: {"STA", IndirectIndexed, 6, (*CPU).STA},
	0x94: {"STY", ZeroPageX, 4, (*CPU).STY},
	0x95: {"STA", ZeroPageX, 4, (*CPU).STA},
	0x96: {"STX", ZeroPageY, 4, (*CPU).STX},
	0x98: {"TYA", Implicit, 2, (*CPU).TYA},
	0x99: {"STA", AbsoluteY, 5, (*CPU).STA},
	0x9A: {"TXS", Implicit, 2, (*CPU).TXS},
	0x9D: {"STA", AbsoluteX, 5, (*CPU).STA},
	0xA0: {"LDY", Immediate, 2, (*CPU).LDY},
	0xA1: {"LDA", IndexedIndirect, 6, (*CPU).LDA},
	0xA2: {"LDX", Immediate, 2, (*CPU).LDX},
	0xA4: {"LDY", ZeroPage, 3, (*CPU).LDY},
	0xA5: {"LDA", ZeroPage, 3, (*CPU).LDA},
	0xA6: {"LDX", ZeroPage, 3, (*CPU).LDX},
	0xA8: {"TAY", Implicit, 2, (*CPU).TAY},
	0xA9: {"LDA", Immediate, 2, (*CPU).LDA},
	0xAA: {"TAX", Implicit, 2, (*CPU).TAX},
	0xAC: {"LDY", Absolute, 4, (*CPU).LDY},
	0xAD: {"LDA", Absolute, 4, (*CPU).LDA},
	0xAE: {"LDX", Absolute, 4, (*CPU).LDX},
	0xB0: {"BCS", Relative, 2, (*CPU).BCS},
	0xB1: {"LDA", IndirectIndexed, 5, (*CPU).LDA},
	0xB4: {"LDY", ZeroPageX, 4, (*CPU).LDY},
	0xB5: {"LDA", ZeroPageX, 4, (*CPU).LDA},
	0xB6: {"LDX", ZeroPageY, 4, (*CPU).LDX},
	0xB8: {"CLV", Implicit, 2, (*CPU).CLV},
	0xB9: {"LDA", AbsoluteY, 4, (*CPU).LDA},
	0xBA: {"TSX", Implicit, 2, (*CPU).TSX},
	0xBC: {"LDY", AbsoluteX, 4, (*CPU).LDY},
	0xBD: {"LDA", AbsoluteX, 4, (*CPU).LDA},
	0xBE: {"LDX", AbsoluteY, 4, (*CPU).LDX},
	0xC0: {"CPY", Immediate, 2, (*CPU).CPY},
	0xC1: {"CMP", IndexedIndirect, 6, (*CPU).CMP},
	0xC4: {"CPY", ZeroPage, 3, (*CPU).CPY},
	0xC5: {"CMP", ZeroPage, 3, (*CPU).CMP},
	0xC6: {"DEC", ZeroPage, 5, (*CPU).DEC},
	0xC8: {"INY", Implicit, 2, (*CPU).INY},
	0xC9: {"CMP", Immediate, 2, (*CPU).CMP},
	0xCA: {"DEX", Implicit, 2, (*CPU).DEX},
	0xCC: {"CPY", Absolute, 4, (*CPU).CPY},
	0xCD: {"CMP", Absolute, 4, (*CPU).CMP},
	0xCE: {"DEC", Absolute, 6, (*CPU).DEC},
	0xD0: {"BNE", Relative, 2, (*CPU).BNE},
	0xD1: {"CMP", IndirectIndexed, 5, (*CPU).CMP},
	0xD5: {"CMP", ZeroPageX, 4, (*CPU).CMP},
	0xD6: {"DEC", ZeroPageX, 6, (*CPU).DEC},
	0xD8: {"CLD", Implicit, 2, (*CPU).CLD},
	0xD9: {"CMP", AbsoluteY, 4, (*CPU).CMP},
	0xDD: {"CMP", AbsoluteX, 4, (*CPU).CMP},
	0xDE: {"DEC", AbsoluteX, 7, (*CPU).DEC},
	0xE0: {"CPX", Immediate, 2, (*CPU).CPX},
	0xE1: {"SBC", IndexedIndirect, 6, (*CPU).SBC},
	0xE4: {"CPX", ZeroPage, 3, (*CPU).CPX},
	0xE5: {"SBC", ZeroPage, 3, (*CPU).SBC},
	0xE6: {"INC", ZeroPage, 5, (*CPU).INC},
	0xE8: {"INX", Implicit, 2, (*CPU).INX},
	0xE9: {"SBC", Immediate, 2, (*CPU).SBC},
	0xEA: {"NOP", Implicit, 2, (*CPU).NOP},
	0xEC: {"CPX", Absolute, 4, (*CPU).CPX},
	0xED: {"SBC", Absolute, 4, (*CPU).SBC},
	0xEE: {"INC", Absolute, 6, (*CPU).INC},
	0xF0: {"BEQ", Relative, 2, (*CPU).BEQ},
	0xF1: {"SBC", IndirectIndexed, 5, (*CPU).SBC},
	0xF5: {"SBC", ZeroPageX, 4, (*CPU).SBC},
	0xF6: {"INC", ZeroPageX, 6, (*CPU).INC},
	0xF8: {"SED", Implicit, 2, (*CPU).SED},
	0xF9: {"SBC", AbsoluteY, 4, (*CPU).SBC},
	0xFD: {"SBC", AbsoluteX, 4, (*CPU).SBC},
	0xFE: {"INC", AbsoluteX, 7, (*CPU).INC},
}

func (cpu *CPU) GetAddress(mode AddressingMode) uint16 {
	switch mode {
	case Immediate:
		addr := cpu.PC
		cpu.PC++
		return addr
	case ZeroPage:
		addr := uint16(cpu.ReadByte(cpu.PC))
		cpu.PC++
		return addr
	case ZeroPageX:
		addr := uint16(cpu.ReadByte(cpu.PC) + cpu.X)
		cpu.PC++
		return addr & 0xFF
	case ZeroPageY:
		addr := uint16(cpu.ReadByte(cpu.PC) + cpu.Y)
		cpu.PC++
		return addr & 0xFF
	case Absolute:
		addr := cpu.ReadWord(cpu.PC)
		cpu.PC += 2
		return addr
	case AbsoluteX:
		addr := cpu.ReadWord(cpu.PC) + uint16(cpu.X)
		cpu.PC += 2
		return addr
	case AbsoluteY:
		addr := cpu.ReadWord(cpu.PC) + uint16(cpu.Y)
		cpu.PC += 2
		return addr
	case Indirect:
		indirect := cpu.ReadWord(cpu.PC)
		cpu.PC += 2
		return cpu.ReadWord(indirect)
	case IndexedIndirect:
		base := cpu.ReadByte(cpu.PC)
		cpu.PC++
		addr := uint16(base + cpu.X)
		return cpu.ReadWord(addr & 0xFF)
	case IndirectIndexed:
		base := cpu.ReadByte(cpu.PC)
		cpu.PC++
		addr := cpu.ReadWord(uint16(base)) + uint16(cpu.Y)
		return addr
	case Relative:
		offset := int8(cpu.ReadByte(cpu.PC))
		cpu.PC++
		return uint16(int32(cpu.PC) + int32(offset))
	default:
		return 0
	}
}