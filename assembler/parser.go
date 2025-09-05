package assembler

import (
	"fmt"
	"strconv"
	"strings"
)

type InstructionType int

const (
	InstMnemonic InstructionType = iota
	InstDirective
	InstLabel
)

type AddressingMode int

const (
	AddrImplicit AddressingMode = iota
	AddrAccumulator
	AddrImmediate
	AddrZeroPage
	AddrZeroPageX
	AddrZeroPageY
	AddrAbsolute
	AddrAbsoluteX
	AddrAbsoluteY
	AddrIndirect
	AddrIndexedIndirect
	AddrIndirectIndexed
	AddrRelative
)

type Instruction struct {
	Type          InstructionType
	Mnemonic      string
	AddressMode   AddressingMode
	Operand       uint16
	OperandLabel  string
	Line          int
	Address       uint16
	DirectiveName string
	DirectiveData []uint8
}

type Parser struct {
	tokens   []Token
	position int
	symbols  *SymbolTable
}

func NewParser(tokens []Token, symbols *SymbolTable) *Parser {
	return &Parser{
		tokens:  tokens,
		symbols: symbols,
	}
}

func (p *Parser) Parse() ([]Instruction, error) {
	var instructions []Instruction
	
	for p.position < len(p.tokens) {
		token := p.currentToken()
		
		if token.Type == TokenEOF {
			break
		}
		
		if token.Type == TokenNewline || token.Type == TokenComment {
			p.advance()
			continue
		}
		
		instruction, err := p.parseInstruction()
		if err != nil {
			return nil, err
		}
		
		if instruction != nil {
			instructions = append(instructions, *instruction)
		}
	}
	
	return instructions, nil
}

func (p *Parser) parseInstruction() (*Instruction, error) {
	token := p.currentToken()
	
	switch token.Type {
	case TokenLabel:
		return p.parseLabel()
		
	case TokenMnemonic:
		return p.parseMnemonic()
		
	case TokenDirective:
		return p.parseDirective()
		
	default:
		return nil, fmt.Errorf("unexpected token '%s' at line %d", token.Value, token.Line)
	}
}

func (p *Parser) parseLabel() (*Instruction, error) {
	token := p.currentToken()
	p.advance()
	
	return &Instruction{
		Type:     InstLabel,
		Mnemonic: token.Value,
		Line:     token.Line,
	}, nil
}

func (p *Parser) parseMnemonic() (*Instruction, error) {
	token := p.currentToken()
	mnemonic := token.Value
	line := token.Line
	p.advance()
	
	addressMode, operand, operandLabel, err := p.parseOperand()
	if err != nil {
		return nil, err
	}
	
	return &Instruction{
		Type:         InstMnemonic,
		Mnemonic:     mnemonic,
		AddressMode:  addressMode,
		Operand:      operand,
		OperandLabel: operandLabel,
		Line:         line,
	}, nil
}

func (p *Parser) parseDirective() (*Instruction, error) {
	token := p.currentToken()
	directive := token.Value
	line := token.Line
	p.advance()
	
	switch strings.ToLower(directive) {
	case ".org":
		token := p.currentToken()
		var operand uint16
		var err error
		
		if token.Type == TokenAbsolute {
			p.advance()
			operand, err = p.parseHexNumber(token.Value[1:])
		} else {
			if token.Type == TokenDollar {
				p.advance()
			}
			operand, err = p.parseNumber()
		}
		
		if err != nil {
			return nil, fmt.Errorf("expected address after .org: %v", err)
		}
		
		return &Instruction{
			Type:          InstDirective,
			DirectiveName: "org",
			Operand:       operand,
			Line:          line,
		}, nil
		
	case ".word":
		token := p.currentToken()
		if token.Type == TokenIdentifier {
			label := token.Value
			p.advance()
			return &Instruction{
				Type:          InstDirective,
				DirectiveName: "word",
				OperandLabel:  label,
				Line:          line,
			}, nil
		}
		
		data, err := p.parseWordData()
		if err != nil {
			return nil, fmt.Errorf("error parsing .word data: %v", err)
		}
		
		return &Instruction{
			Type:          InstDirective,
			DirectiveName: "word",
			DirectiveData: data,
			Line:          line,
		}, nil
		
	case ".byte":
		data, err := p.parseByteData()
		if err != nil {
			return nil, fmt.Errorf("error parsing .byte data: %v", err)
		}
		
		return &Instruction{
			Type:          InstDirective,
			DirectiveName: "byte",
			DirectiveData: data,
			Line:          line,
		}, nil
		
	default:
		return nil, fmt.Errorf("unknown directive '%s' at line %d", directive, line)
	}
}

func (p *Parser) parseOperand() (AddressingMode, uint16, string, error) {
	token := p.currentToken()
	
	if token.Type == TokenEOF || token.Type == TokenNewline || token.Type == TokenComment {
		return AddrImplicit, 0, "", nil
	}
	
	if token.Type == TokenHash {
		p.advance()
		nextToken := p.currentToken()
		
		if nextToken.Type == TokenAbsolute {
			p.advance()
			operand, err := p.parseHexNumber(nextToken.Value[1:])
			if err != nil {
				return AddrImmediate, 0, "", err
			}
			return AddrImmediate, operand, "", nil
		} else if nextToken.Type == TokenDollar {
			p.advance()
			operand, err := p.parseNumber()
			if err != nil {
				return AddrImmediate, 0, "", err
			}
			return AddrImmediate, operand, "", nil
		} else {
			operand, err := p.parseNumber()
			if err != nil {
				return AddrImmediate, 0, "", err
			}
			return AddrImmediate, operand, "", nil
		}
	}
	
	if token.Type == TokenDollar {
		p.advance()
		operand, err := p.parseNumber()
		if err != nil {
			return AddrAbsolute, 0, "", err
		}
		
		if operand <= 0xFF {
			return AddrZeroPage, operand, "", nil
		}
		return AddrAbsolute, operand, "", nil
	}
	
	if token.Type == TokenAbsolute {
		p.advance()
		operand, err := p.parseHexNumber(token.Value[1:])
		if err != nil {
			return AddrAbsolute, 0, "", err
		}
		
		if operand <= 0xFF {
			return AddrZeroPage, operand, "", nil
		}
		return AddrAbsolute, operand, "", nil
	}
	
	if token.Type == TokenIdentifier {
		label := token.Value
		p.advance()
		return AddrAbsolute, 0, label, nil
	}
	
	if token.Type == TokenNumber {
		operand, err := p.parseNumber()
		if err != nil {
			return AddrAbsolute, 0, "", err
		}
		
		if operand <= 0xFF {
			return AddrZeroPage, operand, "", nil
		}
		return AddrAbsolute, operand, "", nil
	}
	
	return AddrImplicit, 0, "", nil
}

func (p *Parser) parseNumber() (uint16, error) {
	token := p.currentToken()
	if token.Type != TokenNumber {
		return 0, fmt.Errorf("expected number, got '%s'", token.Value)
	}
	
	p.advance()
	
	if strings.Contains(strings.ToUpper(token.Value), "A") ||
		strings.Contains(strings.ToUpper(token.Value), "B") ||
		strings.Contains(strings.ToUpper(token.Value), "C") ||
		strings.Contains(strings.ToUpper(token.Value), "D") ||
		strings.Contains(strings.ToUpper(token.Value), "E") ||
		strings.Contains(strings.ToUpper(token.Value), "F") {
		val, err := strconv.ParseUint(token.Value, 16, 16)
		return uint16(val), err
	}
	
	val, err := strconv.ParseUint(token.Value, 10, 16)
	return uint16(val), err
}

func (p *Parser) parseHexNumber(hexStr string) (uint16, error) {
	val, err := strconv.ParseUint(hexStr, 16, 16)
	return uint16(val), err
}

func (p *Parser) parseWordData() ([]uint8, error) {
	token := p.currentToken()
	var data []uint8
	
	if token.Type == TokenIdentifier {
		p.advance()
		data = append(data, 0, 0)
		return data, nil
	}
	
	var operand uint16
	var err error
	
	if token.Type == TokenAbsolute {
		p.advance()
		operand, err = p.parseHexNumber(token.Value[1:])
	} else {
		if token.Type == TokenDollar {
			p.advance()
		}
		operand, err = p.parseNumber()
	}
	
	if err != nil {
		return nil, err
	}
	
	data = append(data, uint8(operand&0xFF), uint8((operand>>8)&0xFF))
	return data, nil
}

func (p *Parser) parseByteData() ([]uint8, error) {
	var data []uint8
	
	for {
		token := p.currentToken()
		if token.Type == TokenEOF || token.Type == TokenNewline || token.Type == TokenComment {
			break
		}
		
		if token.Type == TokenDollar {
			p.advance()
		}
		
		operand, err := p.parseNumber()
		if err != nil {
			return nil, err
		}
		
		data = append(data, uint8(operand&0xFF))
		
		if p.currentToken().Type == TokenComma {
			p.advance()
		} else {
			break
		}
	}
	
	return data, nil
}

func (p *Parser) currentToken() Token {
	if p.position >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.position]
}

func (p *Parser) advance() {
	if p.position < len(p.tokens) {
		p.position++
	}
}