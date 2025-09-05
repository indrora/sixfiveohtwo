package emulator

func (cpu *CPU) LDA(addr uint16) {
	cpu.A = cpu.ReadByte(addr)
	cpu.UpdateZeroAndNegative(cpu.A)
}

func (cpu *CPU) LDX(addr uint16) {
	cpu.X = cpu.ReadByte(addr)
	cpu.UpdateZeroAndNegative(cpu.X)
}

func (cpu *CPU) LDY(addr uint16) {
	cpu.Y = cpu.ReadByte(addr)
	cpu.UpdateZeroAndNegative(cpu.Y)
}

func (cpu *CPU) STA(addr uint16) {
	cpu.WriteByte(addr, cpu.A)
}

func (cpu *CPU) STX(addr uint16) {
	cpu.WriteByte(addr, cpu.X)
}

func (cpu *CPU) STY(addr uint16) {
	cpu.WriteByte(addr, cpu.Y)
}

func (cpu *CPU) ADC(addr uint16) {
	value := cpu.ReadByte(addr)
	carry := uint8(0)
	if cpu.GetFlag(CARRY_FLAG) {
		carry = 1
	}
	
	result := uint16(cpu.A) + uint16(value) + uint16(carry)
	
	cpu.SetFlag(CARRY_FLAG, result > 0xFF)
	cpu.SetFlag(OVERFLOW_FLAG, ((cpu.A^value)&0x80) == 0 && ((cpu.A^uint8(result))&0x80) != 0)
	
	cpu.A = uint8(result)
	cpu.UpdateZeroAndNegative(cpu.A)
}

func (cpu *CPU) SBC(addr uint16) {
	value := cpu.ReadByte(addr)
	carry := uint8(1)
	if cpu.GetFlag(CARRY_FLAG) {
		carry = 1
	} else {
		carry = 0
	}
	
	result := int16(cpu.A) - int16(value) - int16(1-carry)
	
	cpu.SetFlag(CARRY_FLAG, result >= 0)
	cpu.SetFlag(OVERFLOW_FLAG, ((cpu.A^value)&0x80) != 0 && ((cpu.A^uint8(result))&0x80) != 0)
	
	cpu.A = uint8(result)
	cpu.UpdateZeroAndNegative(cpu.A)
}

func (cpu *CPU) AND(addr uint16) {
	cpu.A = cpu.A & cpu.ReadByte(addr)
	cpu.UpdateZeroAndNegative(cpu.A)
}

func (cpu *CPU) ORA(addr uint16) {
	cpu.A = cpu.A | cpu.ReadByte(addr)
	cpu.UpdateZeroAndNegative(cpu.A)
}

func (cpu *CPU) EOR(addr uint16) {
	cpu.A = cpu.A ^ cpu.ReadByte(addr)
	cpu.UpdateZeroAndNegative(cpu.A)
}

func (cpu *CPU) CMP(addr uint16) {
	value := cpu.ReadByte(addr)
	result := cpu.A - value
	cpu.SetFlag(CARRY_FLAG, cpu.A >= value)
	cpu.UpdateZeroAndNegative(result)
}

func (cpu *CPU) CPX(addr uint16) {
	value := cpu.ReadByte(addr)
	result := cpu.X - value
	cpu.SetFlag(CARRY_FLAG, cpu.X >= value)
	cpu.UpdateZeroAndNegative(result)
}

func (cpu *CPU) CPY(addr uint16) {
	value := cpu.ReadByte(addr)
	result := cpu.Y - value
	cpu.SetFlag(CARRY_FLAG, cpu.Y >= value)
	cpu.UpdateZeroAndNegative(result)
}

func (cpu *CPU) BIT(addr uint16) {
	value := cpu.ReadByte(addr)
	result := cpu.A & value
	cpu.SetFlag(ZERO_FLAG, result == 0)
	cpu.SetFlag(OVERFLOW_FLAG, (value&0x40) != 0)
	cpu.SetFlag(NEGATIVE_FLAG, (value&0x80) != 0)
}

func (cpu *CPU) INC(addr uint16) {
	value := cpu.ReadByte(addr) + 1
	cpu.WriteByte(addr, value)
	cpu.UpdateZeroAndNegative(value)
}

func (cpu *CPU) DEC(addr uint16) {
	value := cpu.ReadByte(addr) - 1
	cpu.WriteByte(addr, value)
	cpu.UpdateZeroAndNegative(value)
}

func (cpu *CPU) ASL(addr uint16) {
	value := cpu.ReadByte(addr)
	cpu.SetFlag(CARRY_FLAG, (value&0x80) != 0)
	value <<= 1
	cpu.WriteByte(addr, value)
	cpu.UpdateZeroAndNegative(value)
}

func (cpu *CPU) LSR(addr uint16) {
	value := cpu.ReadByte(addr)
	cpu.SetFlag(CARRY_FLAG, (value&0x01) != 0)
	value >>= 1
	cpu.WriteByte(addr, value)
	cpu.UpdateZeroAndNegative(value)
}

func (cpu *CPU) ROL(addr uint16) {
	value := cpu.ReadByte(addr)
	carry := cpu.GetFlag(CARRY_FLAG)
	cpu.SetFlag(CARRY_FLAG, (value&0x80) != 0)
	value <<= 1
	if carry {
		value |= 0x01
	}
	cpu.WriteByte(addr, value)
	cpu.UpdateZeroAndNegative(value)
}

func (cpu *CPU) ROR(addr uint16) {
	value := cpu.ReadByte(addr)
	carry := cpu.GetFlag(CARRY_FLAG)
	cpu.SetFlag(CARRY_FLAG, (value&0x01) != 0)
	value >>= 1
	if carry {
		value |= 0x80
	}
	cpu.WriteByte(addr, value)
	cpu.UpdateZeroAndNegative(value)
}

func (cpu *CPU) ASLA(addr uint16) {
	cpu.SetFlag(CARRY_FLAG, (cpu.A&0x80) != 0)
	cpu.A <<= 1
	cpu.UpdateZeroAndNegative(cpu.A)
}

func (cpu *CPU) LSRA(addr uint16) {
	cpu.SetFlag(CARRY_FLAG, (cpu.A&0x01) != 0)
	cpu.A >>= 1
	cpu.UpdateZeroAndNegative(cpu.A)
}

func (cpu *CPU) ROLA(addr uint16) {
	carry := cpu.GetFlag(CARRY_FLAG)
	cpu.SetFlag(CARRY_FLAG, (cpu.A&0x80) != 0)
	cpu.A <<= 1
	if carry {
		cpu.A |= 0x01
	}
	cpu.UpdateZeroAndNegative(cpu.A)
}

func (cpu *CPU) RORA(addr uint16) {
	carry := cpu.GetFlag(CARRY_FLAG)
	cpu.SetFlag(CARRY_FLAG, (cpu.A&0x01) != 0)
	cpu.A >>= 1
	if carry {
		cpu.A |= 0x80
	}
	cpu.UpdateZeroAndNegative(cpu.A)
}

func (cpu *CPU) TAX(addr uint16) {
	cpu.X = cpu.A
	cpu.UpdateZeroAndNegative(cpu.X)
}

func (cpu *CPU) TAY(addr uint16) {
	cpu.Y = cpu.A
	cpu.UpdateZeroAndNegative(cpu.Y)
}

func (cpu *CPU) TXA(addr uint16) {
	cpu.A = cpu.X
	cpu.UpdateZeroAndNegative(cpu.A)
}

func (cpu *CPU) TYA(addr uint16) {
	cpu.A = cpu.Y
	cpu.UpdateZeroAndNegative(cpu.A)
}

func (cpu *CPU) TXS(addr uint16) {
	cpu.SP = cpu.X
}

func (cpu *CPU) TSX(addr uint16) {
	cpu.X = cpu.SP
	cpu.UpdateZeroAndNegative(cpu.X)
}

func (cpu *CPU) INX(addr uint16) {
	cpu.X++
	cpu.UpdateZeroAndNegative(cpu.X)
}

func (cpu *CPU) INY(addr uint16) {
	cpu.Y++
	cpu.UpdateZeroAndNegative(cpu.Y)
}

func (cpu *CPU) DEX(addr uint16) {
	cpu.X--
	cpu.UpdateZeroAndNegative(cpu.X)
}

func (cpu *CPU) DEY(addr uint16) {
	cpu.Y--
	cpu.UpdateZeroAndNegative(cpu.Y)
}

func (cpu *CPU) PHA(addr uint16) {
	cpu.Push(cpu.A)
}

func (cpu *CPU) PLA(addr uint16) {
	cpu.A = cpu.Pop()
	cpu.UpdateZeroAndNegative(cpu.A)
}

func (cpu *CPU) PHP(addr uint16) {
	cpu.Push(cpu.P | BREAK_FLAG)
}

func (cpu *CPU) PLP(addr uint16) {
	cpu.P = cpu.Pop()
	cpu.P |= UNUSED_FLAG
	cpu.P &^= BREAK_FLAG
}

func (cpu *CPU) JMP(addr uint16) {
	cpu.PC = addr
}

func (cpu *CPU) JMPI(addr uint16) {
	cpu.PC = cpu.ReadWord(addr)
}

func (cpu *CPU) JSR(addr uint16) {
	cpu.PC--
	cpu.PushWord(cpu.PC)
	cpu.PC = addr
}

func (cpu *CPU) RTS(addr uint16) {
	cpu.PC = cpu.PopWord() + 1
}

func (cpu *CPU) RTI(addr uint16) {
	cpu.P = cpu.Pop()
	cpu.P |= UNUSED_FLAG
	cpu.P &^= BREAK_FLAG
	cpu.PC = cpu.PopWord()
}

func (cpu *CPU) BEQ(addr uint16) {
	if cpu.GetFlag(ZERO_FLAG) {
		cpu.PC = addr
	}
}

func (cpu *CPU) BNE(addr uint16) {
	if !cpu.GetFlag(ZERO_FLAG) {
		cpu.PC = addr
	}
}

func (cpu *CPU) BCS(addr uint16) {
	if cpu.GetFlag(CARRY_FLAG) {
		cpu.PC = addr
	}
}

func (cpu *CPU) BCC(addr uint16) {
	if !cpu.GetFlag(CARRY_FLAG) {
		cpu.PC = addr
	}
}

func (cpu *CPU) BMI(addr uint16) {
	if cpu.GetFlag(NEGATIVE_FLAG) {
		cpu.PC = addr
	}
}

func (cpu *CPU) BPL(addr uint16) {
	if !cpu.GetFlag(NEGATIVE_FLAG) {
		cpu.PC = addr
	}
}

func (cpu *CPU) BVS(addr uint16) {
	if cpu.GetFlag(OVERFLOW_FLAG) {
		cpu.PC = addr
	}
}

func (cpu *CPU) BVC(addr uint16) {
	if !cpu.GetFlag(OVERFLOW_FLAG) {
		cpu.PC = addr
	}
}

func (cpu *CPU) CLC(addr uint16) {
	cpu.SetFlag(CARRY_FLAG, false)
}

func (cpu *CPU) SEC(addr uint16) {
	cpu.SetFlag(CARRY_FLAG, true)
}

func (cpu *CPU) CLI(addr uint16) {
	cpu.SetFlag(INTERRUPT_FLAG, false)
}

func (cpu *CPU) SEI(addr uint16) {
	cpu.SetFlag(INTERRUPT_FLAG, true)
}

func (cpu *CPU) CLV(addr uint16) {
	cpu.SetFlag(OVERFLOW_FLAG, false)
}

func (cpu *CPU) CLD(addr uint16) {
	cpu.SetFlag(DECIMAL_FLAG, false)
}

func (cpu *CPU) SED(addr uint16) {
	cpu.SetFlag(DECIMAL_FLAG, true)
}

func (cpu *CPU) NOP(addr uint16) {
}

func (cpu *CPU) BRK(addr uint16) {
	cpu.PC++
	cpu.PushWord(cpu.PC)
	cpu.Push(cpu.P | BREAK_FLAG)
	cpu.SetFlag(INTERRUPT_FLAG, true)
	cpu.PC = cpu.ReadWord(IRQ_VECTOR)
}