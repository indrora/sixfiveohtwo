package emulator

import (
	"fmt"
	"io/ioutil"
)

const (
	CARRY_FLAG     = 0x01
	ZERO_FLAG      = 0x02
	INTERRUPT_FLAG = 0x04
	DECIMAL_FLAG   = 0x08
	BREAK_FLAG     = 0x10
	UNUSED_FLAG    = 0x20
	OVERFLOW_FLAG  = 0x40
	NEGATIVE_FLAG  = 0x80
)

const (
	STACK_BASE   = 0x0100
	RESET_VECTOR = 0xFFFC
	IRQ_VECTOR   = 0xFFFE
	NMI_VECTOR   = 0xFFFA
)

type CPU struct {
	A  uint8
	X  uint8
	Y  uint8
	SP uint8
	PC uint16
	P  uint8
	
	memory [65536]uint8
	cycles uint64
	
	running bool
}

func NewCPU() *CPU {
	cpu := &CPU{
		SP: 0xFF,
		P:  UNUSED_FLAG,
	}
	
	return cpu
}

func (cpu *CPU) Reset() {
	resetAddr := cpu.ReadWord(RESET_VECTOR)
	cpu.PC = resetAddr
	cpu.SP = 0xFF
	cpu.P = UNUSED_FLAG
	cpu.cycles = 0
	cpu.running = true
}

func (cpu *CPU) LoadROM(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	
	romStart := 0x8000
	if len(data) > 32768 {
		return fmt.Errorf("ROM too large")
	}
	
	copy(cpu.memory[romStart:], data)
	
	return nil
}

func (cpu *CPU) ReadByte(addr uint16) uint8 {
	if addr >= 0xF000 && addr <= 0xFFFF {
		return cpu.handleIORead(addr)
	}
	return cpu.memory[addr]
}

func (cpu *CPU) WriteByte(addr uint16, value uint8) {
	if addr >= 0xF000 && addr <= 0xFFFF {
		cpu.handleIOWrite(addr, value)
		return
	}
	cpu.memory[addr] = value
}

func (cpu *CPU) ReadWord(addr uint16) uint16 {
	lo := uint16(cpu.ReadByte(addr))
	hi := uint16(cpu.ReadByte(addr + 1))
	return (hi << 8) | lo
}

func (cpu *CPU) WriteWord(addr uint16, value uint16) {
	cpu.WriteByte(addr, uint8(value&0xFF))
	cpu.WriteByte(addr+1, uint8((value>>8)&0xFF))
}

func (cpu *CPU) Push(value uint8) {
	cpu.WriteByte(STACK_BASE+uint16(cpu.SP), value)
	cpu.SP--
}

func (cpu *CPU) Pop() uint8 {
	cpu.SP++
	return cpu.ReadByte(STACK_BASE + uint16(cpu.SP))
}

func (cpu *CPU) PushWord(value uint16) {
	cpu.Push(uint8((value >> 8) & 0xFF))
	cpu.Push(uint8(value & 0xFF))
}

func (cpu *CPU) PopWord() uint16 {
	lo := uint16(cpu.Pop())
	hi := uint16(cpu.Pop())
	return (hi << 8) | lo
}

func (cpu *CPU) SetFlag(flag uint8, set bool) {
	if set {
		cpu.P |= flag
	} else {
		cpu.P &^= flag
	}
}

func (cpu *CPU) GetFlag(flag uint8) bool {
	return (cpu.P & flag) != 0
}

func (cpu *CPU) UpdateZeroAndNegative(value uint8) {
	cpu.SetFlag(ZERO_FLAG, value == 0)
	cpu.SetFlag(NEGATIVE_FLAG, (value&0x80) != 0)
}

func (cpu *CPU) handleIORead(addr uint16) uint8 {
	switch addr {
	case 0xF004:
		return cpu.readKeyboard()
	default:
		return cpu.memory[addr]
	}
}

func (cpu *CPU) handleIOWrite(addr uint16, value uint8) {
	switch addr {
	case 0xF001:
		cpu.writeDisplay(value)
	default:
		cpu.memory[addr] = value
	}
}

func (cpu *CPU) readKeyboard() uint8 {
	return 0x00
}

func (cpu *CPU) writeDisplay(value uint8) {
	if value >= 0x20 && value <= 0x7E {
		fmt.Printf("%c", value)
	} else if value == 0x0A || value == 0x0D {
		fmt.Println()
	}
}

func (cpu *CPU) Run() {
	for cpu.running {
		cpu.Step()
	}
}

func (cpu *CPU) Step() {
	opcode := cpu.ReadByte(cpu.PC)
	cpu.PC++
	
	instruction := instructions[opcode]
	if instruction.Execute == nil {
		fmt.Printf("Unknown opcode: 0x%02X at PC: 0x%04X\n", opcode, cpu.PC-1)
		cpu.running = false
		return
	}
	
	var addr uint16
	switch instruction.AddressMode {
	case Implicit, Accumulator:
		addr = 0
	default:
		addr = cpu.GetAddress(instruction.AddressMode)
	}
	
	instruction.Execute(cpu, addr)
	cpu.cycles += uint64(instruction.Cycles)
	
	if cpu.PC == 0 {
		cpu.running = false
	}
}